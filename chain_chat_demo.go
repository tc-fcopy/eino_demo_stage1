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

// demoChainChat 演示使用 Ark ChatModel 与 compose.Chain 进行简单对话，
// 通过 Lambda 对输入进行预处理后再交给模型回答。
func demoChainChat() {
	ctx := context.Background()
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	if err != nil {
		panic(err)
	}

	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output []*schema.Message, err error) {
		processed := input + "回答加上desuwa"
		output = []*schema.Message{{
			Role:    schema.User,
			Content: processed,
		}}
		return output, nil
	})

	chain := compose.NewChain[string, *schema.Message]()
	chain.AppendLambda(lambda).AppendChatModel(model)
	r, err := chain.Compile(ctx)
	if err != nil {
		panic(err)
	}
	answer, err := r.Invoke(ctx, "你好,你可以告诉我你的名字吗？")
	if err != nil {
		panic(err)
	}
	fmt.Println(answer.Content)
}
