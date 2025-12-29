// Package file 文件操作辅助函数
package file

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"GoHub-Service/pkg/app"
	"GoHub-Service/pkg/auth"
	"GoHub-Service/pkg/config"
	"GoHub-Service/pkg/helpers"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
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

func SaveUploadAvatar(c *gin.Context, file *multipart.FileHeader) (string, error) {

	if file == nil {
		return "", errors.New("avatar file is empty")
	}

	if err := validateUpload(file); err != nil {
		return "", err
	}

	biz := "avatars"
	datePath := app.TimenowInTimezone().Format("2006/01/02")
	dateParts := strings.Split(datePath, "/")

	basePath := strings.TrimRight(config.GetString("storage.base_path", "public/uploads"), "/")
	if basePath == "" {
		basePath = "public/uploads"
	}
	publicPrefix := config.GetString("storage.public_prefix", "/uploads")

	userID := auth.CurrentUID(c)
	dirPath := filepath.Join(basePath, biz, dateParts[0], dateParts[1], dateParts[2], userID)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileHash, err := fileHash8(file)
	if err != nil {
		return "", err
	}

	finalName := buildUploadFileName(biz, ext, fileHash)
	tmpName := "raw-" + finalName
	tmpPath := filepath.Join(dirPath, tmpName)

	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		return "", err
	}

	img, err := imaging.Open(tmpPath, imaging.AutoOrientation(true))
	if err != nil {
		return "", err
	}
	resizedAvatar := imaging.Thumbnail(img, 256, 256, imaging.Lanczos)
	finalPath := filepath.Join(dirPath, finalName)
	if err := imaging.Save(resizedAvatar, finalPath); err != nil {
		return "", err
	}

	if err := os.Remove(tmpPath); err != nil {
		return "", err
	}

	publicPath := "/" + path.Join(strings.TrimPrefix(publicPrefix, "/"), biz, dateParts[0], dateParts[1], dateParts[2], userID, finalName)
	return publicPath, nil
}

func validateUpload(file *multipart.FileHeader) error {
	maxSizeMB := config.GetInt64("storage.max_size_mb", 5)
	if maxSizeMB > 0 && file.Size > maxSizeMB*1024*1024 {
		return fmt.Errorf("file exceeds %dMB limit", maxSizeMB)
	}

	rawExt := config.GetString("storage.allowed_ext", "jpg,jpeg,png,gif,webp")
	allowed := buildAllowedExtSet(rawExt)
	if len(allowed) == 0 {
		return nil
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(file.Filename)), ".")
	if _, ok := allowed[ext]; !ok {
		return fmt.Errorf("file type %s is not allowed", ext)
	}

	return nil
}

func buildAllowedExtSet(raw string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, item := range strings.Split(raw, ",") {
		trimmed := strings.TrimSpace(strings.ToLower(item))
		if trimmed != "" {
			result[trimmed] = struct{}{}
		}
	}
	return result
}

func fileHash8(file *multipart.FileHeader) (string, error) {
	opened, err := file.Open()
	if err != nil {
		return "", err
	}
	defer opened.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, opened); err != nil {
		return "", err
	}

	sum := hex.EncodeToString(hasher.Sum(nil))
	if len(sum) < 8 {
		return "", errors.New("failed to build file hash")
	}

	return sum[:8], nil
}

func buildUploadFileName(biz, ext, hash8 string) string {
	cleanedExt := ext
	if cleanedExt != "" && !strings.HasPrefix(cleanedExt, ".") {
		cleanedExt = "." + cleanedExt
	}

	suffix := ""
	if hash8 != "" {
		suffix = "-" + hash8
	}

	return fmt.Sprintf("%s-%d-%s%s%s", biz, app.TimenowInTimezone().UnixNano(), helpers.RandomString(8), suffix, cleanedExt)
}
