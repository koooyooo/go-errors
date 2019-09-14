package errors

import (
	"bytes"
	"fmt"
	"strings"
)

type Error struct {
	Msg            string
	Labels         *Labels
	Err            error
	Stack          *Stack
	RawStackTraces []*Stack // Error階層の源流からコピーした全体Stack情報
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// StackTraceString は登録スタック基づいた範囲でトレース情報を出力
func (e *Error) StackTraceString() string {
	getFileNameFromPath := func(s string) string {
		fragments := strings.Split(s, "/")
		return fragments[len(fragments)-1]
	}
	buff := bytes.Buffer{}
	errors := ListErrorChain(e)
	for i, err := range errors {
		buff.WriteString(fmt.Sprintf("%d. @%s -> (%s: %d) %s [%v] \n", i, err.Stack.FuncName, getFileNameFromPath(err.Stack.File), err.Stack.Line, err.Msg, err.Labels))
		if i == (len(errors) - 1) && err.Err != nil {
			buff.WriteString(err.Err.Error() + "\n")
		}
	}
	return buff.String()
}

// RawStackTraceString は発生地点から全てのスタックを走査したトレース情報を出力
func (e *Error) RawStackTraceString() string {
	buff := bytes.Buffer{}
	for i, s := range e.RawStackTraces {
		buff.WriteString(fmt.Sprintf("%d. %s (%s: %d) \n", i, s.FuncName, s.File, s.Line))
	}
	return buff.String()
}

func (e *Error) Cause() error  {
	chain := ListErrorChain(e)
	if len := len(chain); len > 0 {
		return chain[0].Err
	}
	return nil
}

func ListErrorChain(origin error) []*Error {
	err, ok := origin.(*Error)
	if !ok {
		return []*Error{}
	}
	list := []*Error{err}
	child, ok := err.Err.(*Error)
	if !ok {
		return list
	}
	return append(list, ListErrorChain(child)...)
}

func New(kind *Labels, msg string) *Error {
	stack, _ := NewStack(2)
	return &Error{
		Msg:            msg,
		Labels:         kind,
		Err:            nil,
		Stack:          stack,
		RawStackTraces: NewStackTrace(2),
	}
}

func Errorf(kind *Labels, format string, args ...interface{}) *Error {
	msg := fmt.Sprintf(format, args...)
	stack, _ := NewStack(2)
	return &Error{
		Msg:            msg,
		Labels:         kind,
		Err:            nil,
		Stack:          stack,
		RawStackTraces: NewStackTrace(2),
	}
}

func Wrap(kind *Labels, msg string, err error) *Error {
	stack, _ := NewStack(2)
	e, ok := err.(*Error)
	if !ok {
		return &Error{
			Msg:            msg,
			Labels:         kind,
			Err:            err,
			Stack:          stack,
			RawStackTraces: NewStackTrace(2),
		}
	}
	return &Error{
		Msg:            msg,
		Labels:         kind,
		Err:            e,
		Stack:          stack,
		RawStackTraces: e.RawStackTraces,
	}
}

