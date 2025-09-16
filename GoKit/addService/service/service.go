package service

import "errors"

// go-kit serivce 层
// 专注业务逻辑代码

var (
	ErrBothParametersZero = errors.New("both parameters are zero")
)

func NewAddService() AddService {
	return &addService{}
}

type AddService interface {
	Sum(a, b int) (int, error)
	Concat(a, b string) (string, error)
}

type addService struct{}

func (as *addService) Sum(a, b int) (int, error) {
	if a == 0 && b == 0 {
		return 0, ErrBothParametersZero
	}
	return a + b, nil
}

func (as *addService) Concat(a, b string) (string, error) {
	if a == "" && b == "" {
		return "", ErrBothParametersZero
	}
	return a + b, nil
}
