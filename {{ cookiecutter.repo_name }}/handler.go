package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Financial-Times/service-status-go/buildinfo"
	tidutils "github.com/Financial-Times/transactionid-utils-go"
)

type requestHandler struct {
}

func (handler *requestHandler) sampleMessage(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	// todo: implement handler logic

	err := error(nil)
	switch err {
	case nil:
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(`{"output":"Sample output"}`))
	default:
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func newHttpRequest(ctx context.Context, method string, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	tid, err := tidutils.GetTransactionIDFromContext(ctx)
	if err == nil {
		req.Header.Set(tidutils.TransactionIDHeader, tid)
	}

	req.Header.Set("User-Agent", "{{ cookiecutter.ft_platform }}-{{ cookiecutter.repo_name }}/"+strings.Replace(buildinfo.GetBuildInfo().Version, " ", "-", -1))
	return req, nil
}

type message struct {
	Message string `json:"message"`
}

func writeJSONMessage(w http.ResponseWriter, status int, msg string) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	return enc.Encode(&message{Message: msg})
}
