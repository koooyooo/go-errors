package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func f0() error {
	err := f1()
	return Wrap(L(LabelDomainAuthorization), "fail in calling f1()", err)
}

func f1() error {
	err := f2()
	kindB := L(LabelDomainAuthorization, LabelWithAPICallExternal)
	return Wrap(kindB, "fail in calling f2()", err)
}

func f2() error {
	return Wrap(L(LabelDomainAuthorization, LabelWithAPICallExternal), "API call failed", ApiCallError {
		Msg: "Auth API call",
		URL: "https://test.auth.com",
		Method: "Get",
		StCode: 500,
	})
}

func TestDo(t *testing.T) {
	err := f0()
	error, ok := err.(*Error)
	fmt.Println(err)
	fmt.Println()
	if !ok {
		fmt.Println(err.Error())
	} else {
		error.PrintStackTrace(os.Stdout)
		fmt.Println()
		error.PrintRawStackTrace(os.Stdout)
		fmt.Println()
		if api, ok := error.Cause().(ApiCallError); ok {
			fmt.Printf("Operate API Call Error: %v", api)
		}
	}
}

type ApiCallError struct {
	Msg string
	URL string
	Method string
	StCode int
}

func (e ApiCallError) Error() string {
	j, _ := json.Marshal(e)
	return string(j)
}
