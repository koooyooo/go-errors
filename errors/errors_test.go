package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func f0() error {
	err := f1()
	return Wrap(L(LabelDomainAuthorization), "fail call f1", err)
}

func f1() error {
	err := f2()
	kindB := L(LabelDomainAuthorization, LabelWithAPICallExternal)
	return Wrap(kindB, "fail call f2", err)
}

func f2() error {
	return Wrap(L(LabelDomainAuthorization, LabelWithAPICallExternal), "something bad happened", ApiCallError {
		Msg: "Hello",
		URL: "https://test.jp",
		Method: "Get",
		StCode: 500,
	})
}

func TestDo(t *testing.T) {
	err := f0()
	error, ok := err.(*Error)
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
