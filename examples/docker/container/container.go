package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"log"
)

// InitClient 初始化 Docker 客户端
func InitClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.47"))
	if err != nil {
		panic(err)
	}
	return cli
}

// Close 关闭 Docker 客户端
func Close(cli *client.Client) {
	if err := cli.Close(); err != nil {
		log.Fatalf("Docker客户端关闭失败: %s\n", err)
	}
	log.Printf("Docker客户端关闭成功\n")
}

// RunList 查询正在运行的容器列表
func RunList() ([]container.Summary, error) {
	cli := InitClient()
	defer Close(cli) // 关闭客户端
	// 如果是 docker ps -a 参数，那么需要设置ListOptions结构体中的All字段为true
	options := container.ListOptions{}
	list, err := cli.ContainerList(context.Background(), options)
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

// Create 根据镜像创建一个容器
func Create(image, containerName string) (string, error) {
	ctx := context.Background()
	cli := InitClient()
	defer Close(cli) // 关闭客户端
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

// Start 根据容器ID运行容器
func Start(containerID string) (bool, error) {
	ctx := context.Background()
	cli := InitClient()
	defer Close(cli) // 关闭客户端
	err := cli.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		return false, fmt.Errorf("启动容器失败: %w", err)
	}
	fmt.Printf("容器启动成功！ID: %s\n", containerID[:12])
	return true, nil
}

// GetContainerID 根据容器名称获取容器ID
func GetContainerID(containerName string) (string, error) {
	cli := InitClient()
	defer Close(cli) // 关闭客户端
	ctx := context.Background()
	containerInfo, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		// 如果容器不存在，err 通常会是 client.IsErrNotFound(err) == true
		return "", fmt.Errorf("获取容器信息失败 (名称: %s): %w", containerName, err)
	}
	// 成功获取，返回完整容器 ID
	return containerInfo.ID, nil
}

// Remove 删除容器
func Remove(identifier string, force bool) (bool, error) {
	cli := InitClient()
	defer Close(cli) // 关闭客户端
	ctx := context.Background()
	// 准备删除选项
	options := container.RemoveOptions{
		RemoveVolumes: true,  // 删除挂载的匿名卷（通常建议开启）
		RemoveLinks:   false, // 是否删除链接（一般用不到）
		Force:         force, // 是否强制删除运行中的容器
	}
	// 删除容器
	err := cli.ContainerRemove(ctx, identifier, options)
	if err == nil {
		fmt.Printf("容器已成功删除: %s\n", identifier)
		return true, nil
	}
	return false, err
}
