package rest

import (
	"go/ast"
	"net/http"
	"strconv"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

var httpStatusCodes = map[string]int{
	"StatusContinue":                      100, // RFC 7231, 6.2.1
	"StatusSwitchingProtocols":            101, // RFC 7231, 6.2.2
	"StatusProcessing":                    102, // RFC 2518, 10.1
	"StatusEarlyHints":                    103, // RFC 8297
	"StatusOK":                            200, // RFC 7231, 6.3.1
	"StatusCreated":                       201, // RFC 7231, 6.3.2
	"StatusAccepted":                      202, // RFC 7231, 6.3.3
	"StatusNonAuthoritativeInfo":          203, // RFC 7231, 6.3.4
	"StatusNoContent":                     204, // RFC 7231, 6.3.5
	"StatusResetContent":                  205, // RFC 7231, 6.3.6
	"StatusPartialContent":                206, // RFC 7233, 4.1
	"StatusMultiStatus":                   207, // RFC 4918, 11.1
	"StatusAlreadyReported":               208, // RFC 5842, 7.1
	"StatusIMUsed":                        226, // RFC 3229, 10.4.1
	"StatusMultipleChoices":               300, // RFC 7231, 6.4.1
	"StatusMovedPermanently":              301, // RFC 7231, 6.4.2
	"StatusFound":                         302, // RFC 7231, 6.4.3
	"StatusSeeOther":                      303, // RFC 7231, 6.4.4
	"StatusNotModified":                   304, // RFC 7232, 4.1
	"StatusUseProxy":                      305, // RFC 7231, 6.4.5
	"_":                                   306, // RFC 7231, 6.4.6 (Unused)
	"StatusTemporaryRedirect":             307, // RFC 7231, 6.4.7
	"StatusPermanentRedirect":             308, // RFC 7538, 3
	"StatusBadRequest":                    400, // RFC 7231, 6.5.1
	"StatusUnauthorized":                  401, // RFC 7235, 3.1
	"StatusPaymentRequired":               402, // RFC 7231, 6.5.2
	"StatusForbidden":                     403, // RFC 7231, 6.5.3
	"StatusNotFound":                      404, // RFC 7231, 6.5.4
	"StatusMethodNotAllowed":              405, // RFC 7231, 6.5.5
	"StatusNotAcceptable":                 406, // RFC 7231, 6.5.6
	"StatusProxyAuthRequired":             407, // RFC 7235, 3.2
	"StatusRequestTimeout":                408, // RFC 7231, 6.5.7
	"StatusConflict":                      409, // RFC 7231, 6.5.8
	"StatusGone":                          410, // RFC 7231, 6.5.9
	"StatusLengthRequired":                411, // RFC 7231, 6.5.10
	"StatusPreconditionFailed":            412, // RFC 7232, 4.2
	"StatusRequestEntityTooLarge":         413, // RFC 7231, 6.5.11
	"StatusRequestURITooLong":             414, // RFC 7231, 6.5.12
	"StatusUnsupportedMediaType":          415, // RFC 7231, 6.5.13
	"StatusRequestedRangeNotSatisfiable":  416, // RFC 7233, 4.4
	"StatusExpectationFailed":             417, // RFC 7231, 6.5.14
	"StatusTeapot":                        418, // RFC 7168, 2.3.3
	"StatusMisdirectedRequest":            421, // RFC 7540, 9.1.2
	"StatusUnprocessableEntity":           422, // RFC 4918, 11.2
	"StatusLocked":                        423, // RFC 4918, 11.3
	"StatusFailedDependency":              424, // RFC 4918, 11.4
	"StatusTooEarly":                      425, // RFC 8470, 5.2.
	"StatusUpgradeRequired":               426, // RFC 7231, 6.5.15
	"StatusPreconditionRequired":          428, // RFC 6585, 3
	"StatusTooManyRequests":               429, // RFC 6585, 4
	"StatusRequestHeaderFieldsTooLarge":   431, // RFC 6585, 5
	"StatusUnavailableForLegalReasons":    451, // RFC 7725, 3
	"StatusInternalServerError":           500, // RFC 7231, 6.6.1
	"StatusNotImplemented":                501, // RFC 7231, 6.6.2
	"StatusBadGateway":                    502, // RFC 7231, 6.6.3
	"StatusServiceUnavailable":            503, // RFC 7231, 6.6.4
	"StatusGatewayTimeout":                504, // RFC 7231, 6.6.5
	"StatusHTTPVersionNotSupported":       505, // RFC 7231, 6.6.6
	"StatusVariantAlsoNegotiates":         506, // RFC 2295, 8.1
	"StatusInsufficientStorage":           507, // RFC 4918, 11.5
	"StatusLoopDetected":                  508, // RFC 5842, 7.2
	"StatusNotExtended":                   510, // RFC 2774, 7
	"StatusNetworkAuthenticationRequired": 511, // RFC 6585, 6
}

// GetHttpCodesResponses will extract all variables which type is HTTPCodesResponse
func GetHttpCodesResponses(p parser.Package, responsesMap map[string]nexus.HTTPCodesResponse) map[string]nexus.HTTPCodesResponse {
	for _, spec := range parser.GetNexusSpecs(p, "nexus.HTTPCodesResponse") {
		responsesMap[spec.Name] = extractHttpCodesResponse(spec.Value)
	}
	return responsesMap
}

func extractHttpCodesResponse(val *ast.CompositeLit) nexus.HTTPCodesResponse {
	response := nexus.HTTPCodesResponse{}
	for _, elt := range val.Elts {
		kv := elt.(*ast.KeyValueExpr)
		responseKey := extractHttpCodesKey(kv.Key)
		responseValue := extractHttpCodesValue(kv.Value)
		response[responseKey] = responseValue
	}

	return response
}

func extractHttpCodesKey(key ast.Expr) nexus.ResponseCode {
	switch k := key.(type) {
	case *ast.SelectorExpr:
		if k.Sel.String() == "DefaultHTTPErrorCode" {
			return nexus.DefaultHTTPErrorCode
		}

		return nexus.ResponseCode(httpStatusCodes[k.Sel.String()])
	}

	return 0
}

func extractHttpCodesValue(value ast.Expr) nexus.HTTPResponse {
	res := nexus.HTTPResponse{}
	switch val := value.(type) {
	case *ast.CompositeLit:
		if desc, ok := val.Elts[0].(*ast.KeyValueExpr); ok {
			switch descVal := desc.Value.(type) {
			case *ast.BasicLit:
				descStrVal, err := strconv.Unquote(descVal.Value)
				if err != nil {
					panic(err)
				}
				res.Description = descStrVal
			case *ast.CallExpr:
				if descVal.Fun.(*ast.SelectorExpr).Sel.String() == "StatusText" {
					arg := descVal.Args[0].(*ast.SelectorExpr)
					res.Description = http.StatusText(httpStatusCodes[arg.Sel.String()])
				}
			}
		}
	case *ast.SelectorExpr:
		if val.Sel.String() == "DefaultHTTPError" {
			res = nexus.DefaultHTTPError
		}
	}

	return res
}
