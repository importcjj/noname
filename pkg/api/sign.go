package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dop251/goja"
)

type SignResult struct {
	Nars string `json:"nars"`
	Sesi string `json:"sesi"`
}

type Signer struct {
	jsCode   string
	vm       *goja.Runtime
	signFunc goja.Callable
}

func NewSigner(jsFile string) (*Signer, error) {
	b, err := ioutil.ReadFile(jsFile)
	if err != nil {
		return nil, err
	}

	vm := goja.New()
	_, err = vm.RunString(string(b))
	if err != nil {
		return nil, err
	}

	signFunc, ok := goja.AssertFunction(vm.Get("sign"))
	if !ok {
		panic("Not a function")
	}

	return &Signer{
		jsCode:   string(b),
		vm:       vm,
		signFunc: signFunc,
	}, nil
}

func (s *Signer) Sign(body interface{}) (*SignResult, error) {
	ret, err := s.signFunc(goja.Undefined(), s.vm.ToValue(body))
	if err != nil {
		return nil, err
	}

	retS := ret.String()
	if err != nil {
		return nil, err
	}

	var result SignResult

	err = json.Unmarshal([]byte(retS), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
