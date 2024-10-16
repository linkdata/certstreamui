package stream

import (
	"io"
	"net/http"

	"github.com/google/certificate-transparency-go/loglist3"
)

var (
	// cache response from the "all log list" URL
	logList *loglist3.LogList
)

func GetLogList() (*loglist3.LogList, error) {
	if logList != nil {
		return logList, nil
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Get(loglist3.AllLogListURL)
	if err != nil {
		return nil, err
	}
	json, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ll, err := loglist3.NewFromJSON(json)
	if err != nil {
		return nil, err
	}
	logList = ll
	return ll, nil
}
