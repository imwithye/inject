Inject
===
[![Build Status](https://travis-ci.org/imwithye/inject.svg)](https://travis-ci.org/imwithye/inject)
[![GoDoc](https://godoc.org/github.com/imwithye/inject?status.svg)](https://godoc.org/github.com/imwithye/inject)
Dependency injection library for golang.

## Usage

```go
package main

import "github.com/imwithye/inject"

func main() {
	type user struct {
		Name     string
		Password string `inject:"password"`
		Usertype string `inject:"usertype"`
	}
	injector := New()

	name := "Ciel"
	injector.Map(name)
	password := "123456"
	injector.MapTag(password, "password")
	usertype := "normal"
	injector.MapTag(usertype, "usertype")

	u := user{}
	injector.Inject(&u)
	// u.Name == "Ciel"
	// u.Password == "Ciel"
	// u.Usertype == "normal"

	fn := func(name string) {
		// name == "Ciel"
	}
	injector.Invoke(fn)
}
```
