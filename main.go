package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/graphqlcalls"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/rest"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/traceparser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/workmanager"
	"gopkg.in/yaml.v2"
)

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

func (c *Conf) getConf() *Conf {
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

type CallData struct {
	Key    string `json:"key"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

var restTests rest.RestData

func httpServe() {
	r := mux.NewRouter()
	r.HandleFunc("/tests/{test}", TestHandler).Methods("POST")
	r.HandleFunc("/rest/tests", ListRestTests).Methods("GET")
	r.HandleFunc("/rest/tests", PostRestTests).Methods("POST")
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
	json.NewEncoder(w).Encode(calls)
}

func ListRestTests(w http.ResponseWriter, r *http.Request) {
	//restCallPath := "rest_data.yaml"
	log.Println("Handling ", r.URL.Path)
	var calls []CallData
	for k, spec := range restTests.FuncMap {
		calls = append(calls, CallData{Key: k, Path: spec.Path, Method: spec.Method})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calls)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	test := vars["test"]
	fmt.Fprintf(w, "Monitor the log for status of the test, or the grafana dashboard %s \n", test)
	decoder := json.NewDecoder(r.Body)
	var conf Conf
	decoder.Decode(&conf)
	log.Println(conf)
	go func() {

		//restCallPath := "rest_data.yaml"
		apiGateway := conf.Server.URL

		gqlURL := apiGateway + "/apis/graphql/v1/query"
		// define rest calls
		restTests.ProcessRestCalls(apiGateway)
		// Prepare and run graphql tests
		// GraphQL query worker

		w := workmanager.Worker{
			FuncMap: restTests.FuncMap,
		}
		for _, test := range conf.Tests {
			// Default sample rate
			log.Println("Test Name: ", test.Name)
			var samplingRate float32 = 0.1
			if test.SampleRate > 0 {
				samplingRate = test.SampleRate
			}
			w.ZipkinEndPoint = conf.Server.Zipkin
			w.SampleRate = samplingRate
			w.OpsIterations = test.OpsCount
			w.GqlURL = gqlURL
			// initialize graphql
			w.GraphqlFuncMap = graphqlcalls.GraphqlFuncMap
			if test.OpsCount == 0 && test.Timeout == 0 {
				log.Printf("Connot run tests, One of ops count or timeout for tests have to be provided\n")
			}
			for _, funcKey := range test.Rest {
				w.WorkerType = 0
				testRunner(&w, funcKey, test.Concurrency, test.Timeout, conf.Server.Tsdb)
			}
			for _, funcKey := range test.Graphql {
				w.WorkerType = 1
				testRunner(&w, funcKey, test.Concurrency, test.Timeout, conf.Server.Tsdb)
			}
		}
	}()
	fmt.Println(conf.Tests)
}

func main() {
	// read conf
	httpServe()
	time.Sleep(10 * time.Second)
}

func testRunner(w *workmanager.Worker, funcKey string, concurrency int, timeout int, tsdbConnStr string) {
	log.Println(funcKey)
	w.WorkerStart(funcKey, concurrency, timeout)
	time.Sleep(5 * time.Second)
	content, err := w.GatherTestTraces()
	if err != nil {
		log.Printf("Error getting trace content %v", err)
	} else {
		// retrieve data from zipkin backend
		tsData := traceparser.RetrieveData(funcKey, content)
		for _, data := range tsData {
			log.Printf("%d, %f, %d\n", data.Timestamp, data.Duration, data.Error)
		}
		// insert data onto timescale db
		traceparser.InsertData(tsdbConnStr, tsData)

	}
}
