package endpoint

import (
	"addService/service"
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

// go-kit endpoint 层
// 负责请求和响应的编解码、类型转换等工作

type SumRequest struct {
	A int
	B int
}

type SumResponse struct {
	Result int
	Err    string
}

type ConcatRequest struct {
	S1 string
	S2 string
}

type ConcatResponse struct {
	Result string
	Err    string
}

func makeSumEndpoint(svc *service.AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(SumRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		result, err := (*svc).Sum(req.A, req.B)
		return SumResponse{Result: result, Err: err.Error()}, nil
	}
}

func makeConcatEndpoint(svc *service.AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(ConcatRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		result, err := (*svc).Concat(req.S1, req.S2)
		return ConcatResponse{Result: result, Err: err.Error()}, nil
	}
}
