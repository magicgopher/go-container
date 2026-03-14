package container

import "testing"

// TestContainerRunList 获取正在运行的容器列表单元测试
func TestContainerRunList(t *testing.T) {
	containers, err := ContainerRunList()
	// 如果发生错误，直接终止当前测试
	if err != nil {
		t.Fatalf("执行 ContainerRunList 时发生错误: %v", err)
	}
	// 记录成功获取的数量
	t.Logf("成功获取到 %d 个正在运行的容器", len(containers))
	// 如果有容器在运行，可以通过测试日志打印出它们的名字以供核对
	for _, c := range containers {
		t.Logf("测试检查 - 容器 ID: %s, 容器名称: %v", c.ID[:10], c.Names)
	}
}

// CreateContainer 根据镜像创建容器单元测试
func TestCreateContainer(t *testing.T) {
	image := "hello-world:latest"
	container, err := CreateContainer(image, "hello")
	if err != nil {
		t.Logf("创建容器失败: %v\n", err)
		return
	}
	t.Logf("根据镜像名: %s 成功创建容器ID为 %s 的容器!", image, container)
}

// TestStartContainer 根据容器ID运行容器单元测试
func TestStartContainer(t *testing.T) {
	image := "hello-world:latest"
	container, err := CreateContainer(image, "hello")
	if err != nil {
		t.Logf("创建容器失败: %v\n", err)
		return
	}
	t.Logf("根据镜像名: %s 成功创建容器ID为 %s 的容器!", image, container)
	startContainer, err := StartContainer(container)
	if err != nil {
		t.Logf("容器运行失败: %v\n", err)
		return
	}
	t.Logf("容器ID为 %s 的容器，是否成功运行: %v\n", container, startContainer)
}

// TestGetContainerID 容器名称获取容器ID
func TestGetContainerID(t *testing.T) {
	containerID, err := GetContainerID("hello")
	if err != nil {
		t.Logf("获取容器ID失败: %v\n", err)
		return
	}
	t.Logf("容器ID: %s\n", containerID)
}

// TestRemoveContainer 移除某个不在运行的容器
func TestRemoveContainer(t *testing.T) {
	container, err := RemoveContainer("hello", false)
	if err != nil {
		t.Logf("容器移除失败: %v\n", err)
	}
	t.Logf("容器移除成功! ID: %v\n", container)
}
