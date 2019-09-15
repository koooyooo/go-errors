# errors

## Features

### Nested Error Message
Message
```bash
fail in calling f1() || fail in calling f2() || API call failed || {"Msg":"Auth API call","URL":"https://test.auth.com","Method":"Get","StCode":500}
```

### Custom & Raw StackTraces
- Custom Error Stacks could be registered manually
  - By wrapping errors as many times as you like
  - This manual registration could reduce StackTrace noise
- Raw stack trace will be created automatically
  - When you Create or Wrap errors at the first time

Custom StackTrace
```bash
0. @github.com/xxxxxx/xx-errors/errors.f0 (errors_test.go: 27) -> fail in calling f1() [Authentication] 
1. @github.com/xxxxxx/xx-errors/errors.f1 (errors_test.go: 32) -> fail in calling f2() [Authentication,w/API-Call] 
2. @github.com/xxxxxx/xx-errors/errors.f2 (errors_test.go: 42) -> API call failed [Authentication,w/API-Call] 
Cause: {"Msg":"Auth API call","URL":"https://test.auth.com","Method":"Get","StCode":500} 
```

Raw StackTrace
```bash
0. @runtime.goexit (/usr/local/Cellar/go/1.13/libexec/src/runtime/asm_amd64.s: 1357) 
1. @testing.tRunner (/usr/local/Cellar/go/1.13/libexec/src/testing/testing.go: 909) 
2. @github.com/xxxxxx/xx-errors/errors.TestCheckOutput (/Users/xxx/go/src/github.com/xxxxxx/xx-errors/errors/errors_test.go: 128) 
3. @github.com/xxxxxx/xx-errors/errors.f0 (/Users/xxx/go/src/github.com/xxxxxx/xx-errors/errors/errors_test.go: 26) 
4. @github.com/xxxxxx/xx-errors/errors.f1 (/Users/xxx/go/src/github.com/xxxxxx/xx-errors/errors/errors_test.go: 31) 
5. @github.com/xxxxxx/xx-errors/errors.f2 (/Users/xxx/go/src/github.com/xxxxxx/xx-errors/errors/errors_test.go: 42) 
```

### Labelling
- You can define some tags which describe some situations.
- You can use the tags afterwards.

## Usage
### Create New Error

New Error with message
```go
    e := New(NoLabel, "new error")
    return e;
```

New Error with message by format and args
```go
    e := Errorf(NoLabel, "new %s", "error")
    return e;
```

### Wrap other errors
```go
    err := errors.New("original")
    e := Wrap(NoLabel, "doing xxx errors", err)
    return e
```

### Label Error

You can define and apply custom Labels to logs, to specify what kind of operation it is.
```go
var LabelDomainAuthorization = Label("Authorization")
var LabelDomainUserOperation = Label("UserOperation")

var LabelWithLock = Label("w/Lock")
var LabelWithGoroutine = Label("w/Goroutine")
var LabelWithFileAccess = Label("w/FileAccess")
var LabelWithAPICall = Label("w/API-Call")
```
```go
    err := f2()
    return Wrap(L(LabelDomainAuthentication, LabelWithAPICall), "fail in calling f2()", err)
```

And the log will be labeled like this.
```bash
1. @github.com/xxxxxx/xx-errors/errors.f1 (errors_test.go: 32) -> fail in calling f2() [Authentication,w/API-Call] 
```