package main

import (
	"context"
	"fmt"
	"os"

	ark "github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

// State 用于在 Graph 中保存上下文状态（示例记录不同角色的前置事件）。
type State struct {
	History map[string]any
}

func genFunc(ctx context.Context) *State {
	return &State{History: make(map[string]any)}
}

// demoStatefulGraph 演示带本地状态的 Graph：
// - 根据输入角色路由到不同节点（特朗普/拜登）
// - 在前置处理阶段读写状态，影响最终提问内容
// - 由 Ark ChatModel 生成回答
func DemoStatefulGraph() {
	if err := runStatefulGraph(); err != nil {
		panic(err)
	}
}

func runStatefulGraph() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	ctx := context.Background()

	g := compose.NewGraph[map[string]string, *schema.Message](
		compose.WithGenLocalState(genFunc),
	)

	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (map[string]string, error) {
		_ = compose.ProcessState[*State](ctx, func(ctx context.Context, state *State) error {
			state.History["共和党人_action"] = "话筒不小心抵到你的鼻子"
			state.History["民主党人_action"] = "说话途中话筒突然卡了一下，没声音，几秒后又恢复了"
			return nil
		})

		if input["role"] == "共和党人" {
			return map[string]string{"role": "特朗普", "content": input["content"]}, nil
		}
		if input["role"] == "民主党人" {
			return map[string]string{"role": "拜登", "content": input["content"]}, nil
		}
		return map[string]string{"role": "user", "content": input["content"]}, nil
	})

	trumpLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) ([]*schema.Message, error) {
		return []*schema.Message{
			{Role: schema.System, Content: "你是美国总统特朗普，每次都会站在共和党的角度回答我的问题"},
			{Role: schema.User, Content: input["content"]},
		}, nil
	})

	trumpPreHandler := func(ctx context.Context, input map[string]string, state *State) (map[string]string, error) {
		input["content"] = input["content"] + state.History["共和党人_action"].(string)
		return input, nil
	}

	bidenLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) ([]*schema.Message, error) {
		_ = compose.ProcessState[*State](ctx, func(ctx context.Context, state *State) error {
			input["content"] = input["content"] + state.History["民主党人_action"].(string)
			return nil
		})
		return []*schema.Message{
			{Role: schema.System, Content: "你是美国总统拜登，每次都会站在民主党的角度回答我的问题"},
			{Role: schema.User, Content: input["content"]},
		}, nil
	})

	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{APIKey: os.Getenv("ARK_API_KEY"), Model: os.Getenv("MODEL")})
	if err != nil {
		return err
	}

	if err := g.AddLambdaNode("lambda", lambda); err != nil {
		return err
	}
	if err := g.AddLambdaNode("trumpLambda", trumpLambda, compose.WithStatePreHandler(trumpPreHandler)); err != nil {
		return err
	}
	if err := g.AddLambdaNode("bidenLambda", bidenLambda); err != nil {
		return err
	}
	if err := g.AddChatModelNode("model", model); err != nil {
		return err
	}

	if err := g.AddBranch("lambda", compose.NewGraphBranch(func(ctx context.Context, in map[string]string) (string, error) {
		if in["role"] == "特朗普" {
			return "trumpLambda", nil
		}
		if in["role"] == "拜登" {
			return "bidenLambda", nil
		}
		return "trumpLambda", nil
	}, map[string]bool{"trumpLambda": true, "bidenLambda": true})); err != nil {
		return err
	}

	if err := g.AddEdge(compose.START, "lambda"); err != nil {
		return err
	}
	if err := g.AddEdge("trumpLambda", "model"); err != nil {
		return err
	}
	if err := g.AddEdge("bidenLambda", "model"); err != nil {
		return err
	}
	if err := g.AddEdge("model", compose.END); err != nil {
		return err
	}

	r, err := g.Compile(ctx)
	if err != nil {
		return err
	}
	input := map[string]string{"role": "民主党人", "content": "我应该支持民主党还是共和党？"}
	answer, err := r.Invoke(ctx, input)
	if err != nil {
		return err
	}
	fmt.Println(answer.Content)
	return nil
}
