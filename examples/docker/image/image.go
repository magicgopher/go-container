package image

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// ImgInfo 镜像信息
type ImgInfo struct {
	Repository string
	Tag        string
	ImageID    string
	Created    string
	Size       string
}

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

// List 获取所有镜像列表
func List() ([]ImgInfo, error) {
	cli := InitClient()
	defer Close(cli) // 关闭客户端
	images, err := cli.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		panic(err)
	}
	var imageList []ImgInfo
	for _, img := range images {
		imageList = append(imageList, toImage(img))
	}
	return imageList, nil
}

// toImage 将 image.Summary 转换为内部 Image 结构
func toImage(img image.Summary) ImgInfo {
	repo, tag := parseRepoTag(img.RepoTags)
	return ImgInfo{
		Repository: repo,
		Tag:        tag,
		ImageID:    parseImageID(img.ID),
		Created:    parseCreated(img.Created),
		Size:       parseSize(img.Size),
	}
}

// parseImageID 处理镜像ID
func parseImageID(id string) string {
	withoutPrefix := strings.TrimPrefix(id, "sha256:")
	if len(withoutPrefix) >= 12 {
		return withoutPrefix[:12]
	}
	return withoutPrefix
}

// parseRepoTag 安全地解析 Repository 和 Tag
func parseRepoTag(repoTags []string) (string, string) {
	// 处理悬空镜像 (Dangling images)
	if len(repoTags) == 0 {
		return "<none>", "<none>"
	}
	full := repoTags[0]
	// 使用 LastIndex 防止 URL 中自带端口导致切割错误 (例如 registry:5000/image:tag)
	lastColon := strings.LastIndex(full, ":")
	if lastColon == -1 {
		return full, "latest"
	}
	return full[:lastColon], full[lastColon+1:]
}

// parseCreated 处理镜像列表输出的创建的时间
func parseCreated(timestamp int64) string {
	now := time.Now().Unix()
	diff := now - timestamp
	if diff < 0 {
		return "in the future"
	}
	if diff < 60 {
		return "just now"
	}
	// 分钟
	if diff < 3600 {
		minutes := diff / 60
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	// 小时
	if diff < 86400 {
		hours := diff / 3600
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	// 天
	if diff < 86400*30 {
		days := diff / 86400
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
	// 月（粗略用 30 天）
	if diff < 86400*365 {
		months := diff / (86400 * 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}
	// 年（粗略用 365 天）
	years := diff / (86400 * 365)
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

// parseSize 处理镜像列表输出的镜像大小格式
func parseSize(size int64) string {
	// 处理负数的情况（虽然镜像大小通常没有负数）
	sign := ""
	if size < 0 {
		sign = "-"
		size = -size
	}
	if size < 1000 {
		return fmt.Sprintf("%s%dB", sign, size)
	}
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	unitIndex := 0
	floatSize := float64(size)
	// Docker 默认使用的是 Base 1000 (SI标准单位)
	for floatSize >= 1000 && unitIndex < len(units)-1 {
		floatSize /= 1000
		unitIndex++
	}
	// 针对边界进位处理 (比如 999.9 会进位成 1000，此时应该转为下一个单位)
	if math.Round(floatSize*100)/100 >= 1000 && unitIndex < len(units)-1 {
		floatSize /= 1000
		unitIndex++
	}
	var formattedStr string
	// 根据有效数字长度，动态决定小数位数 (保留3位有效数字)
	if floatSize >= 100 {
		formattedStr = fmt.Sprintf("%.0f", floatSize) // 示例: 749
	} else if floatSize >= 10 {
		formattedStr = fmt.Sprintf("%.1f", floatSize) // 示例: 30.3
	} else {
		formattedStr = fmt.Sprintf("%.2f", floatSize) // 示例: 1.07
	}
	// 移除末尾多余的 0 和小数点，比如 1.00 变成 1，1.20 变成 1.2
	if strings.Contains(formattedStr, ".") {
		formattedStr = strings.TrimRight(formattedStr, "0")
		formattedStr = strings.TrimRight(formattedStr, ".")
	}
	return fmt.Sprintf("%s%s%s", sign, formattedStr, units[unitIndex])
}

// Pull 镜像下载
func Pull(name string) error {
	cli := InitClient()
	defer Close(cli) // 关闭客户端
	pull, err := cli.ImagePull(context.Background(), name, image.PullOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := pull.Close(); err != nil {
			log.Printf("关闭镜像 [%s] 的拉取流时发生错误: %v\n", name, closeErr)
		}
	}()
	// 将拉取过程的日志输出到标准输出（控制台），这样你就能看到下载进度了
	_, err = io.Copy(os.Stdout, pull)
	if err != nil {
		return err
	}
	return nil
}

// Remove 镜像删除
func Remove(imageTag string) (bool, error) {
	// 获取client实例
	cli := InitClient()
	defer Close(cli) // 关闭客户端
	// 设置删除选项
	removeOpts := image.RemoveOptions{
		Force:         true, // 强制删除，即使被容器使用 [7]
		PruneChildren: true, // 删除未使用的父镜像
	}
	// 执行删除
	deleted, err := cli.ImageRemove(context.Background(), imageTag, removeOpts)
	if err != nil {
		// 判断是否是“镜像不存在”的情况
		var target notFoundErr
		if errors.As(err, &target) {
			return false, nil // 镜像本来就不存在，视为成功但没删除
		}
		// 其他错误正常上抛
		return false, fmt.Errorf("镜像删除失败 %s: %w", imageTag, err)
	}
	// deleted 通常会有内容，代表真的删除了某些东西
	if len(deleted) > 0 {
		return true, nil
	}
	// 理论上很少走到这里
	return false, nil
}

type notFoundErr interface {
	NotFound()
}
