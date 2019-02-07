# mapprint [![Travis](https://img.shields.io/travis/Eun/mapprint.svg)](https://travis-ci.org/Eun/mapprint) [![Codecov](https://img.shields.io/codecov/c/github/Eun/mapprint.svg)](https://codecov.io/gh/Eun/mapprint) [![GoDoc](https://godoc.org/github.com/Eun/mapprint?status.svg)](https://godoc.org/github.com/Eun/mapprint) [![go-report](https://goreportcard.com/badge/github.com/Eun/mapprint)](https://goreportcard.com/report/github.com/Eun/mapprint)
Printf for maps and structs

```bash
go get github.com/Eun/mapprint
```

```go
// prints `[14:01:08.005674] Database initialized'
mapprint.Printf("[%0H:%02m:%02s.%06ms] Database initialized", map[string]string{
    "H": 14,
    "m": 1,
    "s": 8,
    "ms": 5674,
})
```

## structs
```go
type User struct {
    Name    string
    Balance int
}

// prints `Hello Joe, your balance is 104€!'
mapprint.Printf("Hello %Name, your balance is %Balance€!", User{
    Name:    "Joe",
    Balance: 104,
})
```

## formating
```go
// prints `Hello        Joe'
mapprint.Printf("Hello %10user!", map[string]string{
    "user": "Joe",
})

// prints `Hello Joe       '
mapprint.Printf("Hello %-10user!", map[string]string{
    "user": "Joe",
})

// prints `Hello    Joe    '
mapprint.Printf("Hello |-10user!", map[string]string{
    "user": "Joe",
})

// prints `Hello Joe, your balance is 0000104.00€!'
mapprint.Printf("Hello %Name, your balance is %010.2Balance€!", map[string]interface{}{
    "Name": "Joe",
    "Balance": 104,
})

// maybe something more exotic?
// prints `Hello Joe, your balance is ABAB104.00€!'
mapprint.Printf("Hello %Name, your balance is %AB10.2Balance€!", map[string]interface{}{
    "Name": "Joe",
    "Balance": 104,
})
```