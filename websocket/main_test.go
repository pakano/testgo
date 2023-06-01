package main

import (
	"encoding/json"
	"test/easyjson"
	"testing"
)

func fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

func BenchmarkStdJsonMarshal(b *testing.B) {
	b.ResetTimer()
	p := easyjson.Person{Name: "zhansgan", Age: 11}
	for i := 0; i < b.N; i++ {
		json.Marshal(&p)
	}
}

func BenchmarkEasyJsonMarshal(b *testing.B) {
	b.ResetTimer()
	p := easyjson.Person{Name: "zhansgan", Age: 11}
	for i := 0; i < b.N; i++ {
		p.MarshalJSON()
	}
}

var str = `{
	"name":"zhangsan",
	"age":10
}`

func BenchmarkStdUnmarshal(b *testing.B) {
	b.ResetTimer()
	p := easyjson.Person{}
	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(str), &p)
	}
}

func BenchmarkEasyUnmarshal(b *testing.B) {
	b.ResetTimer()
	p := easyjson.Person{}
	for i := 0; i < b.N; i++ {
		p.UnmarshalJSON([]byte(str))
	}
}
