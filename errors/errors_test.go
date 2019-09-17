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
	return Wrap(err, L(LabelDomainAuthorization), "fail in calling f1()")
}

func f1() error {
	err := f2()
	return Wrap(err, L(LabelDomainAuthorization, LabelWithAPICallExternal), "fail in calling f2()")
}

func f2() error {
	err := ApiCallError {
		Msg: "Auth API call",
		URL: "https://test.auth.com",
		Method: "Get",
		StCode: 500,
	}
	return Wrap(err, L(LabelDomainAuthorization, LabelWithAPICallExternal), "API call failed")
}

func TestNew(t *testing.T) {
	e := New(L(), "new error")
	assert.Equal(t, "new error", e.Error())
	assert.Equal(t,"new error", e.Msg)
	assert.Nil(t, e.Cause())
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
	assert.Nil(t, e.Cause())
	assert.True(t, strings.HasSuffix(e.Stack.FuncName, "TestErrorf"))

	// runtime.goexit & testing.tRunner & TestNew => 3
	len := len(e.RawStackTraces)
	assert.Equal(t, len, 3)
	assert.True(t, strings.HasSuffix(e.RawStackTraces[len-1].FuncName, "TestErrorf"))
}

func TestWrap(t *testing.T) {
	cause := errors.New("original")
	e := Wrap(cause, NoLabel, "doing xxx errors")
	assert.Equal(t, "doing xxx errors || original", e.Error())
	assert.Equal(t, "doing xxx errors", e.Msg)
	assert.Equal(t, cause, e.Cause())
	assert.True(t, strings.HasSuffix(e.Stack.FuncName, "TestWrap"))

	// runtime.goexit & testing.tRunner & TestNew => 3
	len := len(e.RawStackTraces)
	assert.Equal(t, len, 3)
	assert.True(t, strings.HasSuffix(e.RawStackTraces[len-1].FuncName, "TestWrap"))
}

func TestWrapf(t *testing.T) {
	cause := errors.New("original")
	e := Wrapf(cause, NoLabel, "doing xxx %s", "errors")
	assert.Equal(t, "doing xxx errors || original", e.Error())
	assert.Equal(t, "doing xxx errors", e.Msg)
	assert.Equal(t, cause, e.Cause())
	assert.True(t, strings.HasSuffix(e.Stack.FuncName, "TestWrapf"))

	// runtime.goexit & testing.tRunner & TestNew => 3
	len := len(e.RawStackTraces)
	assert.Equal(t, len, 3)
	assert.True(t, strings.HasSuffix(e.RawStackTraces[len-1].FuncName, "TestWrapf"))
}

func TestErrorMethod(t *testing.T) {
	err := f0()
	error, _ := err.(*Error)
	assert.Equal(t, `fail in calling f1() || fail in calling f2() || API call failed || {"Msg":"Auth API call","URL":"https://test.auth.com","Method":"Get","StCode":500}`, error.Error())
}

func TestStackTraceString(t *testing.T) {
	err := f0()
	error, _ := err.(*Error)
	s := error.StackTraceString()
	// contains FuncNames
	assert.Contains(t, s, "errors.f0")
	assert.Contains(t, s, "errors.f1")
	assert.Contains(t, s, "errors.f2")
	// contains Msgs
	assert.Contains(t, s, "fail in calling f1()")
	assert.Contains(t, s, "fail in calling f2()")
	assert.Contains(t, s, "API call failed")
}

func TestRawStackTraceString(t *testing.T) {
	err := f0()
	error, _ := err.(*Error)
	s := error.RawStackTraceString()
	// contains FuncNames
	assert.Contains(t, s, "errors.f0")
	assert.Contains(t, s, "errors.f1")
	assert.Contains(t, s, "errors.f2")
}

func TestGetCauseAndOperateByTypes(t *testing.T) {
	err := f0()
	error, _ := err.(*Error)

	// get cause and use its info to operate
	apiErr, _ := error.Cause().(ApiCallError)
	assert.Equal(t, "Auth API call", apiErr.Msg)
	assert.Equal(t, "Get", apiErr.Method)
	assert.Equal(t, "https://test.auth.com", apiErr.URL)
	assert.Equal(t, 500, apiErr.StCode)
}

func TestCheckOutput(t *testing.T) {
	err := f0()
	error, _ := err.(*Error)

	// Output Msg
	fmt.Println(err)
	fmt.Println()

	// Output StackTrace
	error.PrintStackTrace(os.Stdout)
	fmt.Println()

	// Output RawStackTrace
	error.PrintRawStackTrace(os.Stdout)
	fmt.Println()

	// Check Cause
	if apiErr, ok := error.Cause().(ApiCallError); ok {
		fmt.Println(apiErr)
	}
}
