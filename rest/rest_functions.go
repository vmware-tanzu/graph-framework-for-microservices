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
	//url := "http://localhost:45192/hr/test2"
	url := fmt.Sprintf("%s/hr/test2", r.ApiGateway)
	method := "GET"

	//client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf("Failed to build request %v", err)
	}
	return req

}

func (r *RestFuncs) PutEmployee() *http.Request {
	//empName O= "e5"
	empName := RandomString(10)
	url := fmt.Sprintf("%s/root/default/employee/%s", r.ApiGateway, empName)
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

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
