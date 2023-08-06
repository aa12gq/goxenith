package file

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// Put 将数据存入文件
func Put(data []byte, to string) error {
	err := os.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Exists 判断文件是否存在
func Exists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func SaveUploadImage(c *gin.Context, file *multipart.FileHeader) (avatarPath string, err error) {

	var image string
	publicPath := "public"
	dirName := fmt.Sprintf("/uploads/images/")
	os.MkdirAll(publicPath+dirName, 0755)

	// 获取文件名和后缀
	filenameWithoutExt := fileNameWithoutExtension(file.Filename)
	fileExt := filepath.Ext(file.Filename)

	// 构建完整的存储路径
	completePath := publicPath + dirName + filenameWithoutExt + fileExt

	if err := c.SaveUploadedFile(file, completePath); err != nil {
		return image, err
	}

	pwd, _ := os.Getwd()
	imgPwd := fmt.Sprintf("%v/%v%v%v", pwd, publicPath, dirName, file.Filename)
	return imgPwd, err
}

// fileNameWithoutExtension 从给定的文件名中获取不带扩展名的文件名
func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func fileNameFromUploadFile(name string, file *multipart.FileHeader) string {
	return name + filepath.Ext(file.Filename)
}
