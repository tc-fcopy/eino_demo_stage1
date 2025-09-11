package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
)

// demoGraphBranch 演示使用 compose.Graph 构建分支图，
// 根据输入路由到不同节点并返回不同结果。
func demoGraphBranch() {
	if err := runGraphBranch(); err != nil {
		panic(err)
	}
}

func runGraphBranch() error {
	ctx := context.Background()

	g := compose.NewGraph[string, string]()

	lambda0 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		if input == "1" {
			return "猫猫", nil
		} else if input == "2" {
			return "耄耋", nil
		} else if input == "3" {
			return "device", nil
		}
		return "", nil
	})

	lambda1 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "喵", nil
	})
	lambda2 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "嗨", nil
	})
	lambda3 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "没用人类了吗", nil
	})

	if err := g.AddLambdaNode("lambda0", lambda0); err != nil {
		return err
	}
	if err := g.AddLambdaNode("lambda1", lambda1); err != nil {
		return err
	}
	if err := g.AddLambdaNode("lambda2", lambda2); err != nil {
		return err
	}
	if err := g.AddLambdaNode("lambda3", lambda3); err != nil {
		return err
	}

	if err := g.AddBranch("lambda0", compose.NewGraphBranch(func(ctx context.Context, in string) (endNode string, err error) {
		if in == "猫猫" {
			return "lambda1", nil
		}
		if in == "耄耋" {
			return "lambda2", nil
		}
		if in == "device" {
			return "lambda3", nil
		}
		return compose.END, nil
	}, map[string]bool{"lambda1": true, "lambda2": true, "lambda3": true, compose.END: true})); err != nil {
		return err
	}

	if err := g.AddEdge(compose.START, "lambda0"); err != nil {
		return err
	}
	if err := g.AddEdge("lambda1", compose.END); err != nil {
		return err
	}
	if err := g.AddEdge("lambda2", compose.END); err != nil {
		return err
	}
	if err := g.AddEdge("lambda3", compose.END); err != nil {
		return err
	}

	r, err := g.Compile(ctx)
	if err != nil {
		return err
	}

	answer, err := r.Invoke(ctx, "1")
	if err != nil {
		return err
	}
	fmt.Println(answer)
	return nil
}
