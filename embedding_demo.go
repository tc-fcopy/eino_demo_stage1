package main

import (
	"context"
	"os"

	embedding "github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/joho/godotenv"
)

// demoEmbedding 使用 Ark 向量模型将多段文本转换为向量，并打印维度与结果。
// 依赖：.env 中配置 ARK_API_KEY 与 EMBEDDER。
func demoEmbedding() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	ctx := context.Background()

	embedder, err := embedding.NewEmbedder(ctx, &embedding.EmbeddingConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("EMBEDDER"),
	})
	if err != nil {
		panic(err)
	}
	input := []string{
		"你好，给我两万块钱",
		"你好，借我两万块钱",
		"八百标兵奔北坡",
	}
	embeddings, err := embedder.EmbedStrings(ctx, input)
	if err != nil {
		panic(err)
	}
	for i, v := range embeddings {
		println("文本", i+1, "的向量维度是", len(v))
		println("文本", i+1, "是 ----- ", v)
	}
}
