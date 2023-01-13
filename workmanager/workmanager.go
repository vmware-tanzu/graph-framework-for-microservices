package workmanager

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/rest"
)

type Worker struct {
	WorkerType     int
	GqlURL         string
	zipkinClient   *zipkinhttp.Client
	tracer         *zipkin.Tracer
	gclient        graphql.Client
	ZipkinEndPoint string
	httpClient     *http.Client
	FuncMap        map[string]rest.SpecData
	GraphqlFuncMap map[string]func(context.Context, graphql.Client)
	stop           chan bool
	WorkData       WorkData
	SampleRate     float32
	moduloRate     int
	OpsIterations  int
	m              sync.Mutex
}

type WorkData struct {
	Duration     []int64
	Average      int64
	High         int64
	Low          int64
	ErrCount     int
	OpsCount     int
	TestStart    time.Time
	TestDuration int64
	TestName     string
}

// workManager - starts and stops workers, manages concurrency and time
func (w *Worker) WorkerStart(job string, concurrency int, runFor int) {
	// wait for start and stop singal for the job
	// using buffered channels, keeps signal sender moving
	w.m.Lock()
	defer w.m.Unlock()
	w.WorkData = WorkData{}
	w.WorkData.TestName = job
	w.WorkData.TestStart = time.Now()
	w.stop = make(chan bool, 1)
	var err error
	if w.ZipkinEndPoint != "" {
		w.tracer, err = w.NewTracer()
		if err != nil {
			log.Fatalf("unable to create tracer: %+v\n", err)
		}
	}
	if w.tracer != nil {
		w.zipkinClient, err = zipkinhttp.NewClient(w.tracer, zipkinhttp.ClientTrace(true))
		if err != nil {
			log.Fatalf("unable to create client: %+v\n", err)
		}
		w.gclient = graphql.NewClient(w.GqlURL, w.zipkinClient)
	} else {
		w.gclient = graphql.NewClient(w.GqlURL, http.DefaultClient)
	}

	// set moduloRats
	w.moduloRate = int(1 / w.SampleRate)
	log.Printf("Sampling rate %f, modulo no - %d\n", w.SampleRate, w.moduloRate)
	go w.StopOnTimeOut(runFor)
	w.startWorkers(concurrency, job)
	// channels have been used so that explicit stop can be added . (It has been removed for now )
	w.WorkData.TestDuration = time.Since(w.WorkData.TestStart).Milliseconds()
}

func (w *Worker) StopOnTimeOut(runFor int) {
	if w.OpsIterations == 0 && runFor > 0 {
		time.Sleep(time.Second * time.Duration(runFor))
		log.Println("Stopping worker after runFor : ", runFor)
		w.stop <- true
		log.Println("Work stopped, closing worker automatically...")
	}
}

func (w *Worker) WorkerStop() {
	log.Println("Stoppping worker on demand")
	w.stop <- true
}

func (w *Worker) startWorkers(concurrency int, job string) {
	// concurrent work that can be done = no. of bool set in the channel
	log.Printf("Starting workers for job %s \n", job)
	work := make(chan bool, concurrency)
	for i := 0; i < concurrency; i++ {
		work <- true
	}

	var workerFunc func(string, chan bool)
	// set workerFunc based on the workerType
	switch w.WorkerType {
	case 0:
		workerFunc = w.doWork
	case 1:
		workerFunc = w.doGraphqlQuery
	}
	// http worker
	w.WorkData.OpsCount = 0
	for loop := true; loop; {
		//count the number of ops
		select {
		case <-work:
			w.WorkData.OpsCount++
			// consume work
			start := time.Now()
			workerFunc(job, work)
			elapsed := time.Since(start)
			if (w.WorkData.OpsCount % w.moduloRate) == 0 {
				w.WorkData.Duration = append(w.WorkData.Duration, elapsed.Milliseconds())
			}
		case <-w.stop:
			loop = false
		}
		// stop on ops count
		if w.OpsIterations > 0 && w.OpsIterations == w.WorkData.OpsCount {
			w.stop <- true
		}
	}
}

// async rest client worker
func (w *Worker) doWork(job string, work chan bool) {
	// get work
	specData := w.FuncMap[job]
	req := GetRestReq(specData)
	var err error
	var res *http.Response
	if w.zipkinClient == nil {
		log.Println(req)

		res, err = w.httpClient.Do(req)
		if err != nil {
			log.Printf("err: unable to do http request: %+v\n", err)
		} else {
			res.Body.Close()
		}
	} else {
		res, err = w.zipkinClient.DoWithAppSpan(req, job)
		if err != nil {
			log.Printf("Err: zipkinclient : unable to do http request: %+v\n", err)
		} else {
			res.Body.Close()
		}
	}
	// work done
	work <- true
	//log.Println("status : ", res.StatusCode)
}

// async work graphql worker
func (w *Worker) doGraphqlQuery(job string, work chan bool) {
	gqlFunc := w.GraphqlFuncMap[job]
	var ctx context.Context
	ctx = context.Background()
	if w.tracer == nil {
		gqlFunc(ctx, w.gclient)
	} else {
		span, _ := w.tracer.StartSpanFromContext(ctx, job)
		ctx = zipkin.NewContext(ctx, span)
		gqlFunc(ctx, w.gclient)
		span.Finish()
	}
	// work done
	work <- true
}

func (d *WorkData) CalculateAverage() {
	d.Low = math.MaxInt64
	d.High = 0
	var sum int64 = 0
	for _, v := range d.Duration {
		if v < d.Low {
			d.Low = v
		}
		if v > d.High {
			d.High = v
		}
		sum += v
	}
	d.Average = sum / int64(len(d.Duration))
}

func (w *Worker) GatherTestTraces(test string) ([]byte, error) {
	url := fmt.Sprintf("%s/zipkin/api/v2/traces?serviceName=http_client&spanName=%s&endTs=%d&limit=2000&lookback=%d", w.ZipkinEndPoint, test, w.WorkData.TestStart.UnixMilli()+w.WorkData.TestDuration, w.WorkData.TestDuration)
	log.Println("tarce URL ", url)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return body, nil

}

func GetRestReq(specData rest.SpecData) *http.Request {
	randString := RandomString(10)
	//newPath := strings.ReplaceAll(spec.Path, "{{random}}", randString)
	url := strings.ReplaceAll(specData.Url, "{{random}}", randString)
	payload := strings.NewReader(specData.Data)
	req, err := http.NewRequest(specData.Method, url, payload)
	if err != nil {
		log.Printf("Error : Failed to build request %v", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return req
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func (w *Worker) NewTracer() (*zipkin.Tracer, error) {
	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(w.ZipkinEndPoint + "/api/v2/spans")

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: "http_client", Port: 8080}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(float64(w.SampleRate))
	if err != nil {
		return nil, err
	}

	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	return t, err
}
