package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/callbacks"
)

func genCallBack() callbacks.Handler {
	handler := callbacks.NewHandlerBuilder().OnStartFn(
		func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
			fmt.Printf("当前节点的%v输入：%v\n", info.Component, input)
			return ctx
		}).OnEndFn(
		func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
			fmt.Printf("当前节点的%v输出：%v\n", info.Component, output)
			return ctx
		}).Build()
	return handler
}
