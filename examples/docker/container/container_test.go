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
