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

var tracer *zipkin.Tracer

var zipkinClient *zipkinhttp.Client

var gclient graphql.Client

// true for http, false for graphql
var workerType bool

type server struct {
	URL    string `yaml:"url"`
	Zipkin string `yaml:"zipkin"`
}

type conf struct {
	Server      server   `yaml:"server"`
	Concurrency int      `yaml:"concurrency"`
	Timeout     int      `yaml:"timeout"`
	Rest        []string `yaml:"rest"`
	Graphql     []string `yaml:"graphql"`
}

func (c *conf) getConf() *conf {
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
	log.Println("con ", c.Concurrency, " timeout ", c.Timeout)
	log.Println(c.Rest, c.Graphql)

	apiGateway = c.Server.URL
	endpointURL = c.Server.Zipkin + "/api/v2/spans"
	gqlURL = apiGateway + "/apis/graphql/v1/query"

	var err error
	tracer, err = newTracer()
	if err != nil {
		log.Fatalf("error out %v", err)
	}
	log.Println("tracer added")
	// add functions

	//workManager(GET_HR, c.Concurrency, c.Timeout)
	//time.Sleep(10 * time.Second)
	zipkinClient, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("error out %v", err)
	}
	gclient = graphql.NewClient(gqlURL, zipkinClient)

	restFuncs := rest.RestFuncs{
		ApiGateway: apiGateway,
	}
	// initialize rest functions
	restFuncs.Init()
	// REST worker
	// create regualr rest client to skip zipkin
	w := workmanager.Worker{
		Tracer:     tracer,
		WorkerType: 0,
		FuncMap:    restFuncs.FuncMap,
		SampleRate: 0.15,
	}

	for _, k := range c.Rest {
		log.Println(k)
		w.WorkManagerInit(k, c.Concurrency)
		w.StartWithAutoStop(c.Timeout)
		w.WorkData.CalculateAverage()
		log.Println(w.WorkData.Average, w.WorkData.Low, w.WorkData.High)
		log.Printf("Work data :- \n \tops count: %d \t err count: %d", w.WorkData.OpsCount, w.WorkData.ErrCount)
	}

	// Prepare and run graphql tests
	graphqlFuncs := graphqlcalls.GraphqlFuncs{
		Gclient: gclient,
	}
	graphqlFuncs.Init()
	// GraphQL query worker
	w2 := workmanager.Worker{
		Tracer:         tracer,
		GraphqlFuncMap: graphqlFuncs.GraphqlFuncMap,
		WorkerType:     1,
		SampleRate:     0.2,
	}
	for _, k := range c.Graphql {
		w2.WorkManagerInit(k, c.Concurrency)
		w2.StartWithAutoStop(c.Timeout)
		w2.WorkData.CalculateAverage()
		log.Println(w2.WorkData.Average, w2.WorkData.Low, w2.WorkData.High)
		log.Printf("Work data :- \n \tops count: %d \t err count: %d", w2.WorkData.OpsCount, w2.WorkData.ErrCount)
	}

	time.Sleep(10 * time.Second)
}
func newTracer() (*zipkin.Tracer, error) {
	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(endpointURL)

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: "http_client", Port: 8080}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(0.1)
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
