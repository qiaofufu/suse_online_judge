package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"path"
)

// savePicture 保存form中上传的图片
// name: form-key
// prefix: 保存图片的前缀
// path: 保存文件的路径， 相对于./file/
func savePicture(ctx *gin.Context, name string, prefix string, paths string) (string, error) {
	picture, err := ctx.FormFile(name)
	if err != nil {
		return "", errors.New( "获取图片失败： " + name)
	}
	suffix := path.Ext(picture.Filename)
	if suffix != ".png" && suffix != ".jpg" {
		return "", errors.New("图片格式错误")
	}
	picture.Filename = prefix + suffix
	if err := ctx.SaveUploadedFile(picture, "./file/" + paths + "/" + picture.Filename); err != nil {
		return "", errors.New("图片保存失败")
	}
	return "/file/"+ paths + "/" + picture.Filename, nil
}
