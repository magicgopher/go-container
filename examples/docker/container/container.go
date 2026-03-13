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
	// 如果是 docker ps -a 参数，那么需要设置ListOptions结构体中的All字段为true
	options := container.ListOptions{}
	list, err := client.ContainerList(context.Background(), options)
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

// CreateContainer 根据镜像创建一个容器
func CreateContainer(image, containerName string) (string, error) {
	ctx := context.Background()
	cli := InitClient()
	defer cli.Close()
	// 容器核心配置（类似 Dockerfile 中的 CMD/ENV/EXPOSE 等）
	config := &container.Config{
		Image: image, // 必须指定的镜像名称（可以带 tag，如 "nginx:alpine"）
		Tty:   false, // 是否分配伪终端（交互式容器可设 true）
		// 可选扩展字段示例：
		// Cmd:   []string{"nginx", "-g", "daemon off;"},  // 覆盖镜像默认命令
		// Env:   []string{"NGINX_PORT=80"},
		// ExposedPorts: map[string]struct{}{"80/tcp": {}},
	}

	// 主机配置（端口映射、资源限制、重启策略等）
	hostConfig := &container.HostConfig{
		// 示例：端口映射（宿主机 8080 -> 容器 80）
		// PortBindings: map[string][]types.PortBinding{
		// 	"80/tcp": {{HostIP: "0.0.0.0", HostPort: "8080"}},
		// },

		// 示例：重启策略（容器退出后自动重启）
		// RestartPolicy: container.RestartPolicy{
		// 	Name:              "unless-stopped",
		// 	MaximumRetryCount: 0,
		// },

		// 示例：内存/CPU 限制
		// Resources: container.Resources{
		// 	Memory:     512 * 1024 * 1024, // 512MB
		// 	MemorySwap: 512 * 1024 * 1024,
		// 	NanoCPUs:   1_000_000_000,     // 1 CPU core
		// },
	}
	// 创建容器
	resp, err := cli.ContainerCreate(
		ctx,
		config,
		hostConfig,
		nil,           // networkingConfig（网络配置，可传 nil 使用默认）
		nil,           // platform（架构，可传 nil 让 daemon 决定）
		containerName, // 容器名称（空字符串让 daemon 自动生成）
	)
	if err != nil {
		return "", fmt.Errorf("创建容器失败: %w", err)
	}
	fmt.Printf("容器创建成功！ID: %s\n", resp.ID[:12])
	return resp.ID, nil
}

// StartContainer 根据容器ID运行容器
func StartContainer(containerID string) (bool, error) {
	ctx := context.Background()
	cli := InitClient()
	defer cli.Close()

	err := cli.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		return false, fmt.Errorf("启动容器失败: %w", err)
	}

	fmt.Printf("容器启动成功！ID: %s\n", containerID[:12])
	return true, nil
}
