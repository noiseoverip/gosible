package main

import (
	"flag"
	"fmt"
	"reflect"
)
import "github.com/traefik/yaegi/interp"
import "github.com/traefik/yaegi/stdlib"

func main() {
	flag.Parse()
	path := flag.Arg(0)
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)
	i.Use(i.Symbols("ansiblego/internal"))

	foo := &Foo{Result: "from another boss"}
	var Symbols = interp.Exports{
		"ansiblego/internal": map[string]reflect.Value{
			"foo": reflect.ValueOf(Bar),
			"Foo2" : reflect.ValueOf(foo),
		},
	}

	i.Use(Symbols)

	_, err := i.EvalPath(path)
	if err != nil {
		panic(err)
	}
	v, err := i.Eval("main.gosible")
	if err != nil {
		panic(err)
	}
	module := v.Interface().(func(string) string)
	result := module("params")
	fmt.Printf("Result: %s", result)
	if err != nil {
		panic(err)
	}
}

type Foo struct {
	Result string
}

func(f Foo) Bar() string {
	return f.Result
}

func Bar() string {
	return "from boss"
}
