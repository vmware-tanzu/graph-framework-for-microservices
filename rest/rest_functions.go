package rest

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// define all the functions to be used

func (r *RestFuncs) GetHR() *http.Request {
	url := fmt.Sprintf("%s/hr/test2", r.ApiGateway)
	return r.GetCall(url)
}

func (r *RestFuncs) PutEmployee() *http.Request {
	//empName O= "e5"
	empName := RandomString(10)
	url := fmt.Sprintf("%s/root/default/employee/%s", r.ApiGateway, empName)
	return r.PutCall(url, `{}`)
}

func (r *RestFuncs) PutCall(url string, json_string string) *http.Request {
	//empName O= "e5"
	//empName := RandomString(10)
	//url := fmt.Sprintf("%s/root/default/employee/%s", r.ApiGateway, empName)
	method := "PUT"

	payload := strings.NewReader(json_string)

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Fatalf("Failed to build request %v", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return req
}

func (r *RestFuncs) GetCall(url string) *http.Request {
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

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
