package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func stroage_upload(userID int64, c *gin.Context) (string, error) {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	URL_COS := "https://tiktok-videostorage-1308838593.cos.ap-shanghai.myqcloud.com"
	u, _ := url.Parse(URL_COS)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		//Timeout: 30 * time.Second, //超时时间
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			// SecretID: os.Getenv("SECRETID_TC"), //TC-TencentCloud
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			// SecretKey: os.Getenv("SECRETKEY_TC"),
			SecretID:  "AKIDOfjY9WfZEepkPn0yh9MK2Uiv12yunY5N",
			SecretKey: "81f4FEyTknbBWbR1LQFpAG9mBQSGqJXV",
		},
	})

	// 对象键（Key）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
	// name := "test/objectPut.go"

	//获取路径
	data, err := c.FormFile("data")
	if err != nil {
		return "", err
	}
	file, err := data.Open()
	if err != nil {
		return "", err
	}
	finalName := fmt.Sprintf("%d_%s", userID, filepath.Base(data.Filename))
	filePath := filepath.Join("./tiktokVideo/", finalName)
	//上传对象
	_, err = client.Object.Put(context.Background(), filePath, file, nil)
	if err != nil {
		return "", err
	}
	// // 1.通过字符串上传对象
	// f := strings.NewReader("test")
	// _, err := client.Object.Put(context.Background(), name, f, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// // 2.通过本地文件上传对象
	// _, err = client.Object.PutFromFile(context.Background(), name, "../test", nil)
	// if err != nil {
	// 	panic(err)
	// }
	// // 3.通过文件流上传对象
	// fd, err := os.Open("./test")
	// if err != nil {
	// 	panic(err)
	// }
	// defer fd.Close()
	// _, err = client.Object.Put(context.Background(), name, fd, nil)
	// if err != nil {
	// 	panic(err)
	// }

	return client.Object.GetObjectURL(filePath).String(), nil
}
