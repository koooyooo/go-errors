# errors

## Features
### Custom & Raw StackTraces
- Custom Error Stacks could be registered manually
  - By wrapping errors as many times as you like
  - This manual registration could reduce StackTrace noise
- Raw stack trace will be created automatically
  - When you Create or Wrap errors at the first time

### Labelling
- You can define some tags which describe some situations.
- You can use the tags afterwards.

## Usage
### Create New Error

New Error with message
```go
    e := New(L(), "new error")
    return e;
```

New Error with message by format and args
```go


```
