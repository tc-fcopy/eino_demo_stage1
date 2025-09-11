package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/browseruse"
	"github.com/joho/godotenv"
)

// demoBrowserUse 演示 BrowserUse 工具：启动浏览器并访问指定 URL。
// 依赖：本地环境支持与可能的浏览器依赖。默认访问 bilibili。
func demoBrowserUse() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	ctx := context.Background()

	but, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		log.Fatal(err)
	}
	url := "https://www.bilibili.com"
	result, err := but.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	time.Sleep(10 * time.Second)
	but.Cleanup()
}
