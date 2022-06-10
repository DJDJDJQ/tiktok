package service

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {
	buf := bytes.NewBuffer(nil)
	videoPath = "D:\\学习\\Go\\字节后端\\抖音项目\\douyin\\public\\" + videoPath
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	saveImg := filepath.Join("./public/", snapshotPath)
	err = imaging.Save(img, saveImg+".png")
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	names := strings.Split(snapshotPath, "\\")
	snapshotName = names[len(names)-1] + ".png"
	return
}

func main() {
	_, err := GetSnapshot("./test.mp4", "test", 1)
	if err != nil {
		return
	}
}
