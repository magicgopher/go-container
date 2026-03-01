package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type Image struct {
	Repository string
	Tag        string
	ImageID    string
	Created    string
	Size       string
}

func InitClient() *client.Client {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.47"))
	if err != nil {
		panic(err)
	}
	return client
}

// ImageList 获取所有镜像列表
func ImageList() ([]Image, error) {
	client := InitClient()
	images, err := client.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		panic(err)
	}
	var imageList []Image
	for _, image := range images {
		imageList = append(imageList, toImage(image))
		// log.Println(image.Created)
	}
	return imageList, nil
}

// toImage 将image.Summary转换为Image
func toImage(image image.Summary) Image {
	return Image{
		Repository: image.RepoTags[0],
		Tag:        version(image.RepoTags[0]),
		ImageID:    imageID(image.ID),
		Created:    created(image.Created),
		Size:       fmt.Sprintf("%d MB", image.Size/1024/1024),
	}
}

// version 处理镜像版本
func version(repoTags string) string {
	res := strings.Split(repoTags, ":")
	return res[len(res)-1]
}

// imageID 处理镜像ID
func imageID(id string) string {
	withoutPrefix := strings.TrimPrefix(id, "sha256:")
	imageID := withoutPrefix[:12]
	return imageID
}

func created(timestamp int64) string {
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
		mins := diff / 60
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
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
