package workmanager

import (
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
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
	WorkDuration   WorkDuration
	SampleRate     int
}

type WorkDuration struct {
	Duration []int64
	Average  int64
	High     int64
	Low      int64
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
	if w.started {
		log.Println("Worker already started")
		return
	}
	w.WorkDuration = WorkDuration{}
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
	switch w.WorkerType {
	case 0:
		// http worker
		count := 0
		for {
			// consume work
			count++
			<-work
			start := time.Now()
			w.doWork(job, work)
			elapsed := time.Since(start)
			if (count % w.SampleRate) == 0 {
				w.WorkDuration.Duration = append(w.WorkDuration.Duration, elapsed.Milliseconds())
			}
		}
	case 1:
		// graphql get worker
		count := 0
		for {
			// consume work
			count++
			<-work
			start := time.Now()
			w.doGraphqlQuery(job, work)
			elapsed := time.Since(start)
			if (count % w.SampleRate) == 0 {
				w.WorkDuration.Duration = append(w.WorkDuration.Duration, elapsed.Milliseconds())
			}
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

func (d *WorkDuration) CalculateAverage() {
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
