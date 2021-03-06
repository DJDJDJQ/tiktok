package service

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Upload(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("content-type")
	contentLen := req.ContentLength

	fmt.Printf("upload content-type:%s,content-length:%d", contentType, contentLen)
	if !strings.Contains(contentType, "multipart/form-data") {
		w.Write([]byte("content-type must be multipart/form-data"))
		return
	}
	if contentLen >= 4*1024*1024 { // 10 MB
		w.Write([]byte("file to large,limit 4MB"))
		return
	}

	err := req.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte("ParseMultipartForm error:" + err.Error()))
		return
	}

	if len(req.MultipartForm.File) == 0 {
		w.Write([]byte("not have any file"))
		return
	}

	for name, files := range req.MultipartForm.File {
		fmt.Printf("req.MultipartForm.File,name=%s", name)

		if len(files) != 1 {
			w.Write([]byte("too many files"))
			return
		}
		if name == "" {
			w.Write([]byte("is not FileData"))
			return
		}

		for _, f := range files {
			handle, err := f.Open()
			if err != nil {
				w.Write([]byte(fmt.Sprintf("unknown error,fileName=%s,fileSize=%d,err:%s", f.Filename, f.Size, err.Error())))
				return
			}

			path := "./" + f.Filename
			dst, _ := os.Create(path)
			io.Copy(dst, handle)
			dst.Close()
			fmt.Printf("successful uploaded,fileName=%s,fileSize=%.2f MB,savePath=%s \n", f.Filename, float64(contentLen)/1024/1024, path)

			w.Write([]byte("successful,url=" + url.QueryEscape(f.Filename)))
		}
	}
}

//map for  Http Content-Type  Http 文件类型对应的content-Type
var HttpContentType = map[string]string{
	".avi":  "video/avi",
	".mp3":  "audio/mp3",
	".mp4":  "video/mp4",
	".wmv":  "video/x-ms-wmv",
	".asf":  "video/x-ms-asf",
	".rm":   "application/vnd.rn-realmedia",
	".rmvb": "application/vnd.rn-realmedia-vbr",
	".mov":  "video/quicktime",
	".m4v":  "video/mp4",
	".flv":  "video/x-flv",
	".jpg":  "image/jpeg",
	".png":  "image/png",
}

func FileStream(c *gin.Context) {
	fileName := c.Query("url")
	filePath := "D:\\学习\\Go\\字节后端\\抖音项目\\titok\\public\\" + fileName

	//打开文件
	fileTmp, errByOpenFile := os.Open(filePath)
	if errByOpenFile != nil {
		c.JSON(http.StatusOK, http.Response{StatusCode: 7001, Status: "获取文件失败"})
	}
	defer fileTmp.Close()
	c.File(filePath)
}

func GetContentType(fileName string) (extension, contentType string) {
	arr := strings.Split(fileName, ".")

	// see: https://tool.oschina.net/commons/
	if len(arr) >= 2 {
		extension = arr[len(arr)-1]
		switch extension {
		case "jpeg", "jpe", "jpg":
			contentType = "image/jpeg"
		case "png":
			contentType = "image/png"
		case "gif":
			contentType = "image/gif"
		case "mp4":
			contentType = "video/mpeg4"
		case "mp3":
			contentType = "audio/mp3"
		}
	}
	// .*（ 二进制流，不知道下载文件类型）
	contentType = "application/octet-stream"
	return
}

func Download(w http.ResponseWriter, req *http.Request) {
	if req.RequestURI == "/favicon.ico" {
		return
	}

	fmt.Printf("download url=%s \n", req.RequestURI)

	filename := req.RequestURI[1:]
	enEscapeUrl, err := url.QueryUnescape(filename)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	f, err := os.Open("./service/data/video/" + enEscapeUrl)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	info, err := f.Stat()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	_, contentType := GetContentType(filename)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))

	f.Seek(0, 0)
	io.Copy(w, f)
}
