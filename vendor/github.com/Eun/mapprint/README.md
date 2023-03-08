# mapprint
[![Actions Status](https://github.com/Eun/mapprint/workflows/CI/badge.svg)](https://github.com/Eun/mapprint/actions)
[![Coverage Status](https://coveralls.io/repos/github/Eun/mapprint/badge.svg?branch=master)](https://coveralls.io/github/Eun/mapprint?branch=master)
[![PkgGoDev](https://img.shields.io/badge/pkg.go.dev-reference-blue)](https://pkg.go.dev/github.com/Eun/mapprint)
[![GoDoc](https://godoc.org/github.com/Eun/mapprint?status.svg)](https://godoc.org/github.com/Eun/mapprint)
[![go-report](https://goreportcard.com/badge/github.com/Eun/mapprint)](https://goreportcard.com/report/github.com/Eun/mapprint)
[![go1.15](https://img.shields.io/badge/go-1.15-blue)](#)
---
Printf for maps and structs

```bash
go get github.com/Eun/mapprint
```

```go
// prints `[14:01:08.005674] Database initialized'
mapprint.Printf("[%0H:%02m:%02s.%06ms] Database initialized", map[string]interface{}{
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