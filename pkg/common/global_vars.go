package common

var Mode string
var SSLEnabled string

func IsModeAdmin() bool {
	if Mode == "" {
		return false
	}
	return Mode == "admin"
}

func IsHttpsEnabled() bool {
	if SSLEnabled == "" {
		return false
	}
	return SSLEnabled == "true"
}

var CustomEndpoints = map[string][]string{"allspark-ui": {"/login", "/*.js/", "/home", "/allspark-static/*"}}
var CustomEndpointSvc = "allspark-ui"
