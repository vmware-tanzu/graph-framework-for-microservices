package workmanager

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/rest"
)

type Worker struct {
	WorkerType     int
	zipkinClient   *zipkinhttp.Client
	Tracer         *zipkin.Tracer
	httpClient     *http.Client
	FuncMap        map[string]rest.SpecData
	GraphqlFuncMap map[string]func()
	stop           chan bool
	WorkData       WorkData
	SampleRate     float32
	moduloRate     int
	OpsIterations  int
	m              sync.Mutex
}

type WorkData struct {
	Duration []int64
	Average  int64
	High     int64
	Low      int64
	ErrCount int
	OpsCount int
}

// workManager - starts and stops workers, manages concurrency and time
func (w *Worker) WorkerStart(job string, concurrency int, runFor int) {
	// wait for start and stop singal for the job
	// using buffered channels, keeps signal sender moving
	w.m.Lock()
	defer w.m.Unlock()
	w.stop = make(chan bool, 1)
	w.WorkData = WorkData{}
	var err error
	if w.Tracer != nil {
		w.zipkinClient, err = zipkinhttp.NewClient(w.Tracer, zipkinhttp.ClientTrace(true))
		if err != nil {
			log.Fatalf("unable to create client: %+v\n", err)
		}
	}

	// set moduloRats
	w.moduloRate = int(1 / w.SampleRate)
	log.Printf("Sampling rate %f, modulo no - %d\n", w.SampleRate, w.moduloRate)
	go w.startWorkers(concurrency, job)
	if w.OpsIterations == 0 && runFor > 0 {
		time.Sleep(time.Second * time.Duration(runFor))
		log.Println("Stopping worker after runFor : ", runFor)
		w.stop <- true
		log.Println("Work stopped, closing worker automatically...")
	}
	<-w.stop
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
	for {
		//count the number of ops
		w.WorkData.OpsCount++
		// consume work
		<-work
		start := time.Now()
		workerFunc(job, work)
		elapsed := time.Since(start)
		if (w.WorkData.OpsCount % w.moduloRate) == 0 {
			w.WorkData.Duration = append(w.WorkData.Duration, elapsed.Milliseconds())
		}
		if w.OpsIterations > 0 && w.OpsIterations == w.WorkData.OpsCount {
			w.stop <- true
			break
		}
	}
}

// async rest client worker
func (w *Worker) doWork(job string, work chan bool) {
	// get work
	specData := w.FuncMap[job]
	req := GetRestReq(specData)
	var res *http.Response
	var err error
	if w.zipkinClient == nil {
		log.Println(req)

		res, err = w.httpClient.Do(req)
		if err != nil {
			log.Fatalf("unable to do http request: %+v\n", err)
		}
	} else {
		res, err = w.zipkinClient.DoWithAppSpan(req, job)
		if err != nil {
			log.Fatalf("zipkinclient : unable to do http request: %+v\n", err)
		}
	}
	defer res.Body.Close()
	// work done
	work <- true
	//log.Println("status : ", res.StatusCode)
	if res.StatusCode >= 400 {
		w.WorkData.ErrCount++
	}
}

// async work graphql worker
func (w *Worker) doGraphqlQuery(job string, work chan bool) {
	gqlFunc := w.GraphqlFuncMap[job]
	ctx := context.Background()
	if w.Tracer == nil {
		gqlFunc()
	} else {
		span, _ := w.Tracer.StartSpanFromContext(ctx, job)
		gqlFunc()
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

func GetRestReq(specData rest.SpecData) *http.Request {
	randString := RandomString(10)
	//newPath := strings.ReplaceAll(spec.Path, "{{random}}", randString)
	url := strings.ReplaceAll(specData.Path, "{{random}}", randString)
	payload := strings.NewReader(specData.Data)
	req, err := http.NewRequest(specData.Method, url, payload)
	if err != nil {
		log.Fatalf("Failed to build request %v", err)
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
