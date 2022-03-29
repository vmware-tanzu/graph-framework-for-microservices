package dmcache

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func renderJSON(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	cachedData, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Error marshaling DM cache data %v", err)
		w.Write([]byte(errMsg))
	}

	_, err = w.Write(cachedData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Error writing DM cache to response writer %v", err)
		w.Write([]byte(errMsg))
	}
}
