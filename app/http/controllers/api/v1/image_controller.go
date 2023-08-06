package v1

import (
	"fmt"
	aliyun "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"goxenith/app/requests"
	"goxenith/pkg/config"
	"goxenith/pkg/file"
	"goxenith/pkg/logger"
	"goxenith/pkg/oss"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
	"os"
)

type ImageController struct {
	BaseAPIController
}

func (ctrl *ImageController) Upload(c *gin.Context) {
	request := requests.UserUpdateAvatarRequest{}
	if ok := requests.Validate(c, &request, requests.ImageUpload); !ok {
		return
	}

	url, err := file.SaveUploadImage(c, request.File)
	if err != nil {
		logger.Error(fmt.Sprintf("图片保存到本地出错, err: %v", err))
		response.Abort500(c, "上传头像失败，请稍后尝试~")
		return
	}
	objectKey := fmt.Sprintf("images/%v", request.File.Filename)
	err = oss.Bucket.PutObjectFromFile(objectKey, url, aliyun.ContentType("image/jpg"))
	if err != nil {
		logger.Error(fmt.Sprintf("阿里云OSS图片上传出错, err: %v", err))
		response.Abort500(c, "上传头像失败，请稍后尝试~")
		return
	}
	os.Remove(url)
	response.JSON(c, &pb.Image{
		Id:             "",
		Url:            fmt.Sprintf("https://%v/%v", config.Get("oss.access_path"), objectKey),
		OriginFilename: request.File.Filename,
	})
}
