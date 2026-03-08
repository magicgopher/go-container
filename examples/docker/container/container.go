package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func InitClient() *client.Client {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.47"))
	if err != nil {
		panic(err)
	}
	return client
}

// ContainerRunList 查询正在运行的容器列表
func ContainerRunList() ([]container.Summary, error) {
	client := InitClient()
	defer client.Close() // 关闭客户端
	list, err := client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("查询容器列表失败: %w", err)
	}
	// 遍历并打印容器的基础信息
	for _, c := range list {
		// 截取前 10 位 ID 保持控制台输出整洁
		fmt.Printf("容器 ID: %s | 镜像: %s | 状态: %s\n", c.ID[:10], c.Image, c.Status)
	}
	return list, nil
}
