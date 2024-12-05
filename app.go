package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"watermarked/internal/file"
	"watermarked/internal/watermark"
)

// App struct
type App struct {
	ctx         context.Context
	watermarker *watermark.Watermarker
}

// NewApp creates a new App application struct
func NewApp() *App {
	wm, err := watermark.NewWatermarker()
	if err != nil {
		panic(err)
	}
	return &App{
		watermarker: wm,
	}
}

// startup is called when the app starts. The context is saved
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// 清理临时文件
	file.CleanupTempFiles()
}

// UploadImage 上传图片
func (a *App) UploadImage(imageData []byte, filename string) (*file.FileInfo, error) {
	// 保存文件到临时目录
	tempFile, err := file.SaveFile(imageData, filename)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 验证图片文件
	if err := file.ValidateImage(tempFile); err != nil {
		os.Remove(tempFile) // 清理无效文件
		return nil, err
	}

	// 获取文件信息
	return file.GetFileInfo(tempFile)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// AddTextWatermark 添加文本水印
func (a *App) AddTextWatermark(imagePath string, options watermark.WatermarkOptions) error {
	// 验证图片文件
	if err := file.ValidateImage(imagePath); err != nil {
		return err
	}

	// 打开原始图片
	img, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("打开图片失败: %v", err)
	}
	defer img.Close()

	// 加载图片
	srcImg, err := watermark.LoadImage(img)
	if err != nil {
		return fmt.Errorf("加载图片失败: %v", err)
	}

	// 添加水印
	result, err := a.watermarker.AddTextWatermark(srcImg, options)
	if err != nil {
		return fmt.Errorf("添加水印失败: %v", err)
	}

	// 生成输出文件路径
	dir := filepath.Dir(imagePath)
	filename := filepath.Base(imagePath)
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	outputPath := filepath.Join(dir, name+"_watermarked"+ext)

	// 保存图片
	err = watermark.SaveImage(result, outputPath, 90)
	if err != nil {
		return fmt.Errorf("保存图片失败: %v", err)
	}

	return nil
}

// AddImageWatermark 添加图片水印
func (a *App) AddImageWatermark(imagePath string, watermarkPath string, options watermark.WatermarkOptions) error {
	// 验证源图片
	if err := file.ValidateImage(imagePath); err != nil {
		return err
	}

	// 验证水印图片
	if err := file.ValidateImage(watermarkPath); err != nil {
		return fmt.Errorf("水印图片无效: %v", err)
	}

	// 打开原始图片
	img, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("打开图片失败: %v", err)
	}
	defer img.Close()

	// 加载图片
	srcImg, err := watermark.LoadImage(img)
	if err != nil {
		return fmt.Errorf("加载图片失败: %v", err)
	}

	// 添加水印
	result, err := a.watermarker.AddImageWatermark(srcImg, watermarkPath, options)
	if err != nil {
		return fmt.Errorf("添加水印失败: %v", err)
	}

	// 生成输出文件路径
	dir := filepath.Dir(imagePath)
	filename := filepath.Base(imagePath)
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	outputPath := filepath.Join(dir, name+"_watermarked"+ext)

	// 保存图片
	err = watermark.SaveImage(result, outputPath, 90)
	if err != nil {
		return fmt.Errorf("保存图片失败: %v", err)
	}

	return nil
}

// GetImagePreview returns a base64 encoded preview of the image
func (a *App) GetImagePreview(filePath string) (string, error) {
	return watermark.GetImagePreview(filePath)
}

// GetSupportedFormats 获取支持的图片格式
func (a *App) GetSupportedFormats() []string {
	formats := make([]string, 0, len(file.SupportedFormats))
	for format := range file.SupportedFormats {
		formats = append(formats, format)
	}
	return formats
}
