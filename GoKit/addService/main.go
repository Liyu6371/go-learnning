package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// EmptyStringError indicates that a string parameter is empty
	EmptyStringError = errors.New("empty string paramater")
)

type AddService interface {
	Sum(context.Context, int, int) (int, error)
	Concat(context.Context, string, string) (string, error)
}

type addServiceInst struct{}

func (a *addServiceInst) Sum(ctx context.Context, a1 int, a2 int) (int, error) {
	return a1 + a2, nil
}

func wrapperSumEndpoint(svc AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(SumRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		result, err := svc.Sum(ctx, req.A, req.B)
		if err != nil {
			return SumResponse{Err: err.Error()}, nil
		}
		return SumResponse{Result: result}, nil
	}
}

func (a *addServiceInst) Concat(ctx context.Context, s1 string, s2 string) (string, error) {
	if s1 == "" && s2 == "" {
		return "", EmptyStringError
	}
	return s1 + s2, nil
}

func wrapperConcatEndpoint(svc AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ConcatRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		result, err := svc.Concat(ctx, req.S1, req.S2)
		if err != nil {
			return ConcatResponse{Err: err.Error()}, nil
		}
		return ConcatResponse{Result: result}, nil
	}
}

type SumRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

type SumResponse struct {
	Result int    `json:"result"`
	Err    string `json:"err,omitempty"`
}

type ConcatRequest struct {
	S1 string `json:"s1"`
	S2 string `json:"s2"`
}

type ConcatResponse struct {
	Result string `json:"result"`
	Err    string `json:"err,omitempty"`
}

func decodeSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request SumRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeConcatRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ConcatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {
	svc := &addServiceInst{}

	sumHandler := httptransport.NewServer(
		wrapperSumEndpoint(svc),
		decodeSumRequest,
		encodeResponse,
	)

	concatHandler := httptransport.NewServer(
		wrapperConcatEndpoint(svc),
		decodeConcatRequest,
		encodeResponse,
	)
	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("failed to start server: %v", err)
	}
}
