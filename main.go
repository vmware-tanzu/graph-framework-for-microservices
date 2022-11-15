package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/graphqlcalls"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/rest"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/workmanager"
	"gopkg.in/yaml.v2"
)

var endpointURL string

// apiGateway         = "http://localhost:45192"
var apiGateway string
var gqlURL string

var zipkinClient *zipkinhttp.Client

var gclient graphql.Client

// true for http, false for graphql
var workerType bool

type server struct {
	URL    string `yaml:"url"`
	Zipkin string `yaml:"zipkin"`
}

type Test struct {
	Concurrency int      `yaml:"concurrency"`
	Timeout     int      `yaml:"timeout"`
	OpsCount    int      `yaml:"ops_count"`
	SampleRate  float32  `yaml:"sample_rate"`
	Rest        []string `yaml:"rest"`
	Graphql     []string `yaml:"graphql"`
	Name        string   `yaml:"name"`
}

type conf struct {
	Server server `yaml:"server"`
	Tests  []Test `yaml:"tests"`
}

func (c *conf) getConf() *conf {
	//yamlFile, err := ioutil.ReadFile("conf.yaml")
	yamlFile, err := ioutil.ReadFile("/root/config/conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func main() {
	// read conf
	var c conf
	c.getConf()

	var r rest.RestData
	restCallPath := "/root/rest-calls/rest_data.yaml"
	//restCallPath := "rest_data.yaml"
	apiGateway = c.Server.URL

	// define rest calls
	r.GetRestData(restCallPath)
	r.ProcessRestCalls(apiGateway)

	endpointURL = c.Server.Zipkin + "/api/v2/spans"
	gqlURL = apiGateway + "/apis/graphql/v1/query"

	var err error
	var tracer *zipkin.Tracer
	tracer, err = newTracer(0.1)
	if err != nil {
		log.Fatalf("error out %v", err)
	}
	// add functions

	//client := http.Client{}
	for _, v := range r.Spec {
		log.Println(v.Name, v.Path, v.Data, v.Method)
		log.Printf("key %s, req %v", v.Name, r.FuncMap[v.Name])
	}

	//workManager(GET_HR, c.Concurrency, c.Timeout)
	//time.Sleep(10 * time.Second)
	zipkinClient, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("error out %v", err)
	}
	gclient = graphql.NewClient(gqlURL, zipkinClient)
	// Prepare and run graphql tests
	graphqlFuncs := graphqlcalls.GraphqlFuncs{
		Gclient: gclient,
	}
	graphqlFuncs.Init()
	// GraphQL query worker

	w := workmanager.Worker{
		FuncMap:        r.FuncMap,
		GraphqlFuncMap: graphqlFuncs.GraphqlFuncMap,
	}

	for _, test := range c.Tests {
		// Default sample rate
		log.Println("Test Name: ", test.Name)
		var samplingRate float32 = 0.1
		if test.SampleRate > 0 {
			samplingRate = test.SampleRate
		}
		tracer, err := newTracer(test.SampleRate)
		if err != nil {
			log.Fatalf("error out %v", err)
		}
		w.Tracer = tracer
		w.SampleRate = samplingRate
		w.OpsIterations = test.OpsCount
		if test.OpsCount == 0 && test.Timeout == 0 {
			log.Printf("Connot run tests, One of ops count or timeout for tests have to be provided\n")
		}
		for _, funcKey := range test.Rest {
			w.WorkerType = 0
			testRunner(&w, funcKey, test.Concurrency, test.Timeout)
		}
		for _, funcKey := range test.Graphql {
			w.WorkerType = 1
			testRunner(&w, funcKey, test.Concurrency, test.Timeout)
		}
	}
	time.Sleep(10 * time.Second)
}

func testRunner(w *workmanager.Worker, funcKey string, concurrency int, timeout int) {
	log.Println(funcKey)
	w.WorkerStart(funcKey, concurrency, timeout)
	w.WorkData.CalculateAverage()
	log.Println(w.WorkData.Average, w.WorkData.Low, w.WorkData.High)
	log.Printf("Work data :- \n \tops count: %d \t err count: %d", w.WorkData.OpsCount, w.WorkData.ErrCount)

}

func newTracer(sampleRate float32) (*zipkin.Tracer, error) {
	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(endpointURL)

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: "http_client", Port: 8080}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(float64(sampleRate))
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
