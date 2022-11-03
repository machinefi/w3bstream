package main

import "context"

// Foo 结构体
type Foo struct {
	i int
}

// Bar 接口
type Bar interface {
	Do(ctx context.Context) error
}

// main方法
func main() {
	_ = 1
}
