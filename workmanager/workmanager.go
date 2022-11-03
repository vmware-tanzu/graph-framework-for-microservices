package workmanager

import (
	"log"
	"net/http"
	"time"

	"github.com/hasura/go-graphql-client"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

type Worker struct {
	WorkerType     int
	ZipkinClient   *zipkinhttp.Client
	Gclient        graphql.Client
	FuncMap        map[string]func() *http.Request
	GraphqlFuncMap map[string]func(graphql.Client)
	start          chan bool
	stop           chan bool
	started        bool
}

// workManager - starts and stops workers, manages concurrency and time
func (w *Worker) WorkManager(job string, concurrency int) {
	// wait for start and stop singal for the job
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

}

// StartWithAutoStop : runFor - run for n seconds
func (w *Worker) StartWithAutoStop(runFor int) {
	w.start <- true
	w.started = true
	time.Sleep(time.Second * time.Duration(runFor))
	w.stop <- true
	log.Println("Work stopped, closing worker...")
}

func (w *Worker) Start() {
	if w.started {
		log.Println("Worker already started")
	}
	log.Println("Starting worker")
	w.start <- true
}

func (w *Worker) Stop() {
	log.Println("Work stopped, closing worker...")
	w.stop <- true
}
func (w *Worker) startWorkers(concurrency int, job string) {
	// concurrent work that can be done = no. of bool set in the channel
	work := make(chan bool, concurrency)
	for i := 0; i < concurrency; i++ {
		work <- true
	}
	switch w.WorkerType {
	case 0:
		// http worker
		for {
			// consume work
			<-work
			w.doWork(job, work)
		}
	case 1:
		// graphql get worker
		for {
			<-work
			w.doGraphqlQuery(job, work)
		}
	}

}

// async work
func (w *Worker) doWork(job string, work chan bool) {
	// get work
	req := w.FuncMap[job]()
	req.Header.Add("accept", "application/json")
	var res *http.Response
	res, err := w.ZipkinClient.DoWithAppSpan(req, job)
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}
	defer res.Body.Close()
	// work done
	work <- true

}

// async work
func (w *Worker) doGraphqlQuery(job string, work chan bool) {
	// get work
	gqlFunc := w.GraphqlFuncMap[job]
	gqlFunc(w.Gclient)
	// work done
	work <- true
}
