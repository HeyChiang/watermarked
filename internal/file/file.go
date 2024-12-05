package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 支持的图片格式
var SupportedFormats = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".bmp":  true,
}

// FileInfo 文件信息结构
type FileInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
}

// ValidateImage 验证图片文件
func ValidateImage(filePath string) error {
	// 检查文件是否存在
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("文件不存在: %v", err)
	}

	// 检查是否是文件
	if info.IsDir() {
		return fmt.Errorf("路径指向的是目录而不是文件")
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))
	if !SupportedFormats[ext] {
		return fmt.Errorf("不支持的文件格式: %s", ext)
	}

	return nil
}

// GetFileInfo 获取文件信息
func GetFileInfo(filePath string) (*FileInfo, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:      info.Name(),
		Path:      filePath,
		Size:      info.Size(),
		Extension: strings.ToLower(filepath.Ext(info.Name())),
	}, nil
}

// CreateTempDir 创建临时目录
func CreateTempDir() (string, error) {
	tempDir := filepath.Join(os.TempDir(), "watermarked")
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		return "", err
	}
	return tempDir, nil
}

// SaveFile 保存文件到临时目录
func SaveFile(data []byte, filename string) (string, error) {
	tempDir, err := CreateTempDir()
	if err != nil {
		return "", err
	}

	// 生成唯一的文件名
	tempFile := filepath.Join(tempDir, filename)
	err = os.WriteFile(tempFile, data, 0644)
	if err != nil {
		return "", err
	}

	return tempFile, nil
}

// CleanupTempFiles 清理临时文件
func CleanupTempFiles() error {
	tempDir := filepath.Join(os.TempDir(), "watermarked")
	return os.RemoveAll(tempDir)
}
