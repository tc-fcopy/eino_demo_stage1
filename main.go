package main

import "fmt"

// 入口文件：可根据需要切换要运行的演示函数。

func main() {
	// demoEmbedding()            // Ark 向量化示例
	// demoIndexerVikingDB()      // VikingDB 向量索引示例
	// demoBrowserUse()           // BrowserUse 浏览器操作示例
	// demoChainChat()            // Chain + ChatModel 简单对话示例
	// demoGraphBranch()          // Graph 分支路由示例
	fmt.Println("hello world")
	DemoStatefulGraph() // 带状态的 Graph + ChatModel 示例（默认）
}
