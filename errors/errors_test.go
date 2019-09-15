package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

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

func f0() error {
	err := f1()
	return Wrap(L(LabelDomainAuthorization), "fail in calling f1()", err)
}

func f1() error {
	err := f2()
	return Wrap(L(LabelDomainAuthorization, LabelWithAPICallExternal), "fail in calling f2()", err)
}

func f2() error {
	err := ApiCallError {
		Msg: "Auth API call",
		URL: "https://test.auth.com",
		Method: "Get",
		StCode: 500,
	}
	return Wrap(L(LabelDomainAuthorization, LabelWithAPICallExternal), "API call failed", err)
}

func TestNew(t *testing.T) {
	e := New(L(), "new error")
	assert.Equal(t, "new error", e.Error())
	assert.Equal(t,"new error", e.Msg)
	assert.True(t, strings.HasSuffix(e.Stack.FuncName, "TestNew"))

	// runtime.goexit & testing.tRunner & TestNew => 3
	len := len(e.RawStackTraces)
	assert.Equal(t, len, 3)
	assert.True(t, strings.HasSuffix(e.RawStackTraces[len -1].FuncName, "TestNew"))
}

func TestErrorf(t *testing.T) {
	e := Errorf(L(), "new %s", "error")
	assert.Equal(t, "new error", e.Error())
	assert.Equal(t, "new error", e.Msg)
	assert.True(t, strings.HasSuffix(e.Stack.FuncName, "TestErrorf"))

	// runtime.goexit & testing.tRunner & TestNew => 3
	len := len(e.RawStackTraces)
	assert.Equal(t, len, 3)
	assert.True(t, strings.HasSuffix(e.RawStackTraces[len-1].FuncName, "TestErrorf"))
}

func TestWrap(t *testing.T) {
	cause := errors.New("original")
	e := Wrap(NoLabel, "doing xxx errors", cause)
	assert.Equal(t, "doing xxx errors || original", e.Error())
	assert.Equal(t, "doing xxx errors", e.Msg)
	assert.True(t, strings.HasSuffix(e.Stack.FuncName, "TestWrap"))
}

func TestErrorMethod(t *testing.T) {
	err := f0()
	error, _ := err.(*Error)
	assert.Equal(t, `fail in calling f1() || fail in calling f2() || API call failed || {"Msg":"Auth API call","URL":"https://test.auth.com","Method":"Get","StCode":500}`, error.Error())
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
