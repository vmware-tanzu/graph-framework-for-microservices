package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/graphqlcalls"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/rest"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/traceparser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/workmanager"
)

var testMap map[string]*workmanager.Worker

type Server struct {
	URL    string `yaml:"url" json:"url"`
	Zipkin string `yaml:"zipkin" json:"zipkin"`
	Tsdb   string `yaml:"tsdb" json:"tsdb"`
}

type Test struct {
	Concurrency int      `yaml:"concurrency" json:"concurrency"`
	Timeout     int      `yaml:"timeout" json:"timeout"`
	OpsCount    int      `yaml:"ops_count" json:"ops_count"`
	SampleRate  float32  `yaml:"sample_rate" json:"sample_rate"`
	Rest        []string `yaml:"rest" json:"rest"`
	Graphql     []string `yaml:"graphql" json:"graphql"`
	Name        string   `yaml:"name" json:"name"`
}

type Conf struct {
	Server Server `yaml:"server" json:"server"`
	Tests  []Test `yaml:"tests" json:"tests"`
}

type CallData struct {
	Key    string `json:"key"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type QueryData struct {
	Key   string `json:"key"`
	Query string `json:"query"`
}

var restTests rest.RestData
var gqlTests graphqlcalls.GQLData

func httpServe() {
	r := mux.NewRouter()
	r.HandleFunc("/tests/{test}", TestHandler).Methods("POST")
	r.HandleFunc("/tests/{test}", TestStopHandler).Methods("DELETE")
	r.HandleFunc("/rest/tests", ListRestTests).Methods("GET")
	r.HandleFunc("/rest/tests", PostRestTests).Methods("POST")
	r.HandleFunc("/gql/tests", ListGQLTests).Methods("GET")
	r.HandleFunc("/gql/tests", PostGQLTests).Methods("POST")
	log.Println("Starting http server to serve tests")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func PostRestTests(w http.ResponseWriter, r *http.Request) {
	//restCallPath := "rest_data.yaml"
	log.Println("Handling ", r.URL.Path)
	// define rest calls
	restTests.ReadRestData(r.Body)
	restTests.ProcessRestCalls("")
	var calls []CallData
	for k, spec := range restTests.FuncMap {
		calls = append(calls, CallData{Key: k, Path: spec.Path, Method: spec.Method})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(calls)
}

func ListRestTests(w http.ResponseWriter, r *http.Request) {
	//restCallPath := "rest_data.yaml"
	log.Println("Handling ", r.URL.Path)
	var calls []CallData
	for k, spec := range restTests.FuncMap {
		calls = append(calls, CallData{Key: k, Path: spec.Path, Method: spec.Method})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(calls)
}

func PostGQLTests(w http.ResponseWriter, r *http.Request) {
	//restCallPath := "rest_data.yaml"
	log.Println("Handling ", r.URL.Path)
	// define rest calls
	gqlTests.ReadQueryData(r.Body)
	gqlTests.ProcessGQLCalls()
	var calls []QueryData
	for k, spec := range gqlTests.GQLFuncMap {
		calls = append(calls, QueryData{Key: k, Query: spec.Query})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(calls)
}

func ListGQLTests(w http.ResponseWriter, r *http.Request) {
	//restCallPath := "rest_data.yaml"
	log.Println("Handling ", r.URL.Path)
	var calls []QueryData
	for k, spec := range gqlTests.GQLFuncMap {
		calls = append(calls, QueryData{Key: k, Query: spec.Query})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(calls)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	testKey := vars["test"]
	decoder := json.NewDecoder(r.Body)
	var conf Conf
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Fprintf(w, "Error reading the test conf\n")
		return
	}
	log.Println(conf)
	if testMap[testKey] != nil {
		fmt.Fprintf(w, "The test %s already exists, use a different tests name\n", testKey)
		return
	}
	fmt.Fprintf(w, "Monitor the log for status of the test, or the grafana dashboard %s \n", testKey)
	go func() {

		//restCallPath := "rest_data.yaml"
		apiGateway := conf.Server.URL

		gqlURL := apiGateway + "/apis/graphql/v1/query"
		// define rest calls
		restTests.ProcessRestCalls(apiGateway)
		gqlTests.ProcessGQLCalls()
		// Prepare and run graphql tests
		// GraphQL query worker

		worker := workmanager.Worker{
			FuncMap:    restTests.FuncMap,
			GQLFuncMap: gqlTests.GQLFuncMap,
		}
		for _, test := range conf.Tests {
			// Default sample rate
			log.Println("Test Name: ", test.Name)
			var samplingRate float32 = 0.1
			if test.SampleRate > 0 {
				samplingRate = test.SampleRate
			}
			worker.ZipkinEndPoint = conf.Server.Zipkin
			worker.SampleRate = samplingRate
			worker.OpsIterations = test.OpsCount
			worker.GqlURL = gqlURL
			// initialize graphql
			if test.OpsCount == 0 && test.Timeout == 0 {
				log.Printf("Connot run tests, One of ops count or timeout for tests have to be provided\n")
			}
			for _, funcKey := range test.Rest {
				worker.WorkerType = 0
				testRunner(&worker, funcKey, test.Concurrency, test.Timeout, conf.Server.Tsdb, testKey)
			}
			for _, funcKey := range test.Graphql {
				worker.WorkerType = 1
				testRunner(&worker, funcKey, test.Concurrency, test.Timeout, conf.Server.Tsdb, testKey)
			}

		}
	}()
	fmt.Println(conf.Tests)
}

func TestStopHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	testKey := vars["test"]
	fmt.Fprintf(w, "Test %s is being stopped \n", testKey)
	if testMap[testKey] != nil {
		worker := testMap[testKey]
		worker.WorkerStop()
	} else {
		fmt.Fprintf(w, "Test %s doesn't exist \n", testKey)
	}
}

func main() {
	// read conf
	testMap = make(map[string]*workmanager.Worker)
	httpServe()
}

func testRunner(w *workmanager.Worker, funcKey string, concurrency int, timeout int, tsdbConnStr string, test string) {
	log.Println(funcKey)
	testMap[test] = w
	defer delete(testMap, test)
	/*
		_, ok := w.FuncMap[funcKey]
		if !ok {
			log.Printf("test doesn't exist %s", funcKey)
		}
	*/
	w.WorkerStart(funcKey, concurrency, timeout)
	time.Sleep(5 * time.Second)
	content, err := w.GatherTestTraces(funcKey)
	if err != nil {
		log.Printf("Error getting trace content %v", err)
	} else {
		// retrieve data from zipkin backend
		tsData := traceparser.RetrieveData(funcKey, content)
		for _, data := range tsData {
			log.Printf("%v, %f, %d\n", data.Timestamp, data.Duration, data.Error)
		}
		// insert data onto timescale db
		traceparser.InsertData(tsdbConnStr, tsData)
	}
}
