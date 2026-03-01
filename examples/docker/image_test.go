package docker

import (
	"fmt"
	"testing"
	"time"
)

func TestImageList(t *testing.T) {
	images, err := ImageList()
	if err != nil {
		t.Fatal(err)
	}
	for _, image := range images {
		t.Log(image)
	}
}

func TestVersion(t *testing.T) {
	// version := version("hello-world:latest")
	// t.Log(version)
	t.Log(len("b42de2e2ef12"))
}

func TestTimeToString(t *testing.T) {
	var timestamp int64 = 1765284049
	tm := time.Unix(timestamp, 0)
	fmt.Println(tm)
}
