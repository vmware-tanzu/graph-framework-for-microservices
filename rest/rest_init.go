package rest

import "net/http"

type RestFuncs struct {
	ApiGateway string
	FuncMap    map[string]func() *http.Request
}

// function keys
const (
	// rest function keys
	PUT_EMPLOYEE = "put_employee"
	GET_HR       = "get_hr"
)

// add the map of function keys to function calls
func (r *RestFuncs) Init() {
	r.FuncMap = make(map[string]func() *http.Request)
	funcMap := r.FuncMap
	funcMap[PUT_EMPLOYEE] = r.PutEmployee
	funcMap[GET_HR] = r.GetHR
}
