package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/gqlclient"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/workmanager"
	"gopkg.in/yaml.v2"
)

var endpointURL string

// apiGateway         = "http://localhost:45192"
var apiGateway string
var gqlURL string

// function keys
const (
	// rest function keys
	PUT_EMPLOYEE = "put_employee"
	GET_HR       = "get_hr"
	// graphql function keys
	GET_MANAGERS      = "get_managers"
	GET_EMPLOYEE_ROLE = "get_employee_role"
)

type BuildReq func() *http.Request

var funcMap map[string]func() *http.Request

var graphqlFuncMap map[string]func(graphql.Client)

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

	funcMap = make(map[string]func() *http.Request)
	graphqlFuncMap = make(map[string]func(graphql.Client))
	var err error
	tracer, err = newTracer()
	if err != nil {
		log.Fatalf("error out %v", err)
	}
	log.Println("tracer added")
	// tracer can now be used to create spans.
	// create global zipkin traced http client
	zipkinClient, err = zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	// add functions
	funcMap[PUT_EMPLOYEE] = putEmployee
	funcMap[GET_HR] = getHR

	graphqlFuncMap[GET_MANAGERS] = gqlGetManagers
	graphqlFuncMap[GET_EMPLOYEE_ROLE] = gqlGetEmployeeRole
	//workManager(GET_HR, c.Concurrency, c.Timeout)
	//time.Sleep(10 * time.Second)
	gclient = graphql.NewClient(gqlURL, zipkinClient)

	// REST worker
	// create regualr rest client to skip zipkin
	w := workmanager.Worker{
		ZipkinClient: zipkinClient,
		WorkerType:   0,
		FuncMap:      funcMap,
		SampleRate:   0.15,
	}

	for _, k := range c.Rest {
		log.Println(k)
		w.WorkManager(k, c.Concurrency)
		w.StartWithAutoStop(c.Timeout)
		w.WorkDuration.CalculateAverage()
		log.Println(w.WorkDuration.Average, w.WorkDuration.Low, w.WorkDuration.High)
	}

	// GraphQL query worker
	w2 := workmanager.Worker{
		GraphqlFuncMap: graphqlFuncMap,
		Gclient:        gclient,
		WorkerType:     1,
		SampleRate:     0.2,
	}
	for _, k := range c.Graphql {
		w2.WorkManager(k, c.Concurrency)
		w2.StartWithAutoStop(c.Timeout)
		w2.WorkDuration.CalculateAverage()
		log.Println(w2.WorkDuration.Average, w2.WorkDuration.Low, w2.WorkDuration.High)
	}

	time.Sleep(10 * time.Second)
}

func getHR() *http.Request {
	//url := "http://localhost:45192/hr/test2"
	url := fmt.Sprintf("%s/hr/test2", apiGateway)
	method := "GET"

	//client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf("Failed to build request %v", err)
	}
	return req

}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
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

func putEmployee() *http.Request {
	//empName := "e5"
	empName := RandomString(10)
	url := fmt.Sprintf("%s/root/default/employee/%s", apiGateway, empName)
	method := "PUT"

	payload := strings.NewReader(`{}`)

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Fatalf("Failed to build request %v", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return req
}

func gqlGetManagers(gclient graphql.Client) {
	ctx := context.Background()
	span, ctx := tracer.StartSpanFromContext(ctx, GET_MANAGERS)
	_, err := gqlclient.Managers(ctx, gclient)
	if err != nil {
		log.Printf("Failed to build request %v", err)
	}
	span.Finish()
}

func gqlGetEmployeeRole(gclient graphql.Client) {
	ctx := context.Background()
	span, ctx := tracer.StartSpanFromContext(ctx, GET_EMPLOYEE_ROLE)
	_, err := gqlclient.Employees(ctx, gclient)
	if err != nil {
		log.Printf("Failed to build request %v", err)
	}
	span.Finish()
}
