package errors

import "runtime"

type Stack struct {
	Pc       uintptr
	FuncName string
	File     string
	Line     int
}

func NewStack(skip int) (*Stack, bool) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil, false
	}
	fn := runtime.FuncForPC(pc)
	return &Stack{pc, fn.Name(), file, line}, true
}

func NewStackTrace(skip int) (stackTrace []*Stack) {
	for i := skip; ; i++ {
		stack, ok := NewStack(i)
		if !ok {
			return
		}
		stackTrace = append([]*Stack{stack}, stackTrace...)
	}
}
