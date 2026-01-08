package main

import (
	"fmt"
	"log"

	"github.com/bytectlgo/mem0-go/client"
	"github.com/bytectlgo/mem0-go/types"
)

func main() {
	// 创建客户端
	mem0, err := client.NewMemoryClient(client.ClientOptions{
		APIKey: "your-api-key",
	})
	if err != nil {
		log.Fatal(err)
	}

	// 添加内存
	memories, err := mem0.Add("Hello, World!", types.MemoryOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added memory: %+v\n", memories[0])

	// 搜索内存
	results, err := mem0.Search("Hello", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Search results: %+v\n", results)

	// 获取内存
	memories_, err := mem0.Search("Hello", nil)
	if err != nil {
		log.Fatal(err)
	}
	if len(memories) == 0 {
		log.Fatal("No memory found")
	}
	memory := memories_[0]
	fmt.Printf("Got memory: %+v\n", memory)

	// 更新内存
	updated, err := mem0.Update(memory.ID, "Hello, Updated World!")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated memory: %+v\n", updated[0])

	// 获取内存历史
	history, err := mem0.History(memory.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Memory history: %+v\n", history)

	// 删除内存
	err = mem0.Delete(memory.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Memory deleted successfully")
}
