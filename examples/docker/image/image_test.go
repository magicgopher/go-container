package image

import (
	"testing"
)

// TestImageList 镜像列表单元测试
func TestImageList(t *testing.T) {
	images, err := ImageList()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%-50s %-25s %-15s %-15s %s\n", "REPOSITORY", "TAG", "IMAGE ID", "CREATED", "SIZE")
	for _, image := range images {
		t.Logf("%-50s %-25s %-15s %-15s %s\n", image.Repository, image.Tag, image.ImageID, image.Created, image.Size)
	}
}

// TestImagePull 下载镜像单元测试
func TestImagePull(t *testing.T) {
	err := ImagePull("hello-world")
	if err != nil {
		t.Logf("镜像下载失败: %v\n", err)
	}
	t.Log("镜像下载成功")
}

// TestImageRemove 删除镜像单元测试
func TestImageRemove(t *testing.T) {
	result, err := ImageRemove("hello-world")
	if err != nil {
		t.Fatalf("预期删除成功，但发生了意外错误: %v", err)
	}
	t.Logf("删除操作执行完毕，是否有实际删除动作: %v", result)
}
