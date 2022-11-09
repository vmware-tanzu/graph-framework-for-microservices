package workmanager

import (
	"context"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

type Worker struct {
	WorkerType     int
	zipkinClient   *zipkinhttp.Client
	Tracer         *zipkin.Tracer
	httpClient     *http.Client
	FuncMap        map[string]func() *http.Request
	GraphqlFuncMap map[string]func()
	start          chan bool
	stop           chan bool
	started        bool
	WorkData       WorkData
	SampleRate     float32
	moduloRate     int
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
func (w *Worker) WorkManagerInit(job string, concurrency int) {
	// wait for start and stop singal for the job
	w.WorkData = WorkData{}
	var err error
	if w.Tracer != nil {
		w.zipkinClient, err = zipkinhttp.NewClient(w.Tracer, zipkinhttp.ClientTrace(true))
		if err != nil {
			log.Fatalf("unable to create client: %+v\n", err)
		}
	}
	w.start = make(chan bool)
	w.stop = make(chan bool)
	go func() {
		for i := 0; i < 2; i++ {
			select {
			// start job on signal
			case <-w.start:
				go w.startWorkers(concurrency, job)
			// stop job on signal
			case <-w.stop:
				log.Println("exiting worker ")
			}
		}
	}()

	// set moduloRate
	w.moduloRate = int(1 / w.SampleRate)
	log.Printf("Sampling rate %f, modulo no - %d\n", w.SampleRate, w.moduloRate)

}

// StartWithAutoStop : runFor - run for n seconds
func (w *Worker) StartWithAutoStop(runFor int) {
	if w.started {
		log.Println("Worker already started")
		return
	}
	w.start <- true
	w.started = true
	time.Sleep(time.Second * time.Duration(runFor))
	w.stop <- true
	w.started = false
	log.Println("Work stopped, closing worker...")
}

func (w *Worker) Start() {
	if w.started {
		log.Println("Worker already started")
		return
	}
	log.Println("Starting worker")
	w.start <- true
}

func (w *Worker) Stop() {
	log.Println("Work stopped, closing worker...")
	w.stop <- true
	w.started = false
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
	}
}

// async rest client worker
func (w *Worker) doWork(job string, work chan bool) {
	// get work
	req := w.FuncMap[job]()
	req.Header.Add("accept", "application/json")
	var res *http.Response
	var err error
	if w.zipkinClient == nil {
		res, err = w.httpClient.Do(req)
		if err != nil {
			log.Fatalf("unable to do http request: %+v\n", err)
		}
	} else {
		res, err = w.zipkinClient.DoWithAppSpan(req, job)
		if err != nil {
			log.Fatalf("unable to do http request: %+v\n", err)
		}
	}
	defer res.Body.Close()
	// work done
	work <- true
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
