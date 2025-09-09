package main

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	if err != nil {
		panic(err)
	}
	//input := []*schema.Message{
	//	schema.SystemMessage("你是美国总统特朗普"),
	//	schema.UserMessage("你好,你叫什么名字"),
	//}

	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个{role}"),
		&schema.Message{
			Role:    schema.User,
			Content: "请帮帮我，我想解决{task}",
		},
	)

	params := map[string]any{
		"role": "美国总统川普",
		"task": "帮我计算一下美国对日本关税增加了多少",
	}

	messages, err := template.Format(ctx, params)

	response, err := model.Generate(ctx, messages)
	if err != nil {
		panic(err)
	}

	println(response.Content)

	//reader, err := model.Stream(ctx, input)
	//if err != nil {
	//	panic(err)
	//}
	//defer reader.Close()
	//for {
	//	chunk, err := reader.Recv()
	//	if err != nil {
	//		break
	//	}
	//	print(chunk.Content)
	//}

}
