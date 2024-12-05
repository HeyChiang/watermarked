package watermark

import (
	"bytes"
	"image"
	"image/color"
	stdDraw "image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"

	"encoding/base64"
	"fmt"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Position 水印位置
type Position string

const (
	Center      Position = "center"
	TopLeft     Position = "topLeft"
	TopRight    Position = "topRight"
	BottomLeft  Position = "bottomLeft"
	BottomRight Position = "bottomRight"
	Tiled       Position = "tiled" // 平铺
)

// WatermarkOptions 水印选项
type WatermarkOptions struct {
	// 文本水印选项
	Text       string     `json:"text"`       // 水印文本
	TextSize   float64    `json:"textSize"`   // 字体大小
	TextColor  color.RGBA `json:"textColor"`  // 文字颜色
	FontFamily string     `json:"fontFamily"` // 字体名称

	// 图片水印选项
	Scale float64 `json:"scale"` // 图片水印缩放比例

	// 通用选项
	Opacity  float64  `json:"opacity"`  // 透明度 (0-1)
	Angle    float64  `json:"angle"`    // 旋转角度
	Spacing  float64  `json:"spacing"`  // 水印间距
	Position Position `json:"position"` // 水印位置
	Margin   int      `json:"margin"`   // 边距
}

type Watermarker struct {
	font *truetype.Font
}

func NewWatermarker() (*Watermarker, error) {
	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	return &Watermarker{font: f}, nil
}

// calculatePosition 计算水印位置
func calculatePosition(baseWidth, baseHeight, watermarkWidth, watermarkHeight int, pos Position, margin int) (x, y int) {
	switch pos {
	case TopLeft:
		return margin, margin
	case TopRight:
		return baseWidth - watermarkWidth - margin, margin
	case BottomLeft:
		return margin, baseHeight - watermarkHeight - margin
	case BottomRight:
		return baseWidth - watermarkWidth - margin, baseHeight - watermarkHeight - margin
	case Center:
		return (baseWidth - watermarkWidth) / 2, (baseHeight - watermarkHeight) / 2
	default:
		return (baseWidth - watermarkWidth) / 2, (baseHeight - watermarkHeight) / 2
	}
}

// rotatePoint 旋转点坐标
func rotatePoint(x, y, cx, cy int, angle float64) (int, int) {
	rad := angle * math.Pi / 180
	cos := math.Cos(rad)
	sin := math.Sin(rad)

	// 将点移动到原点
	nx := float64(x - cx)
	ny := float64(y - cy)

	// 旋转
	rx := nx*cos - ny*sin
	ry := nx*sin + ny*cos

	// 移回原位置
	return int(rx) + cx, int(ry) + cy
}

// AddTextWatermark 添加文本水印
func (w *Watermarker) AddTextWatermark(img image.Image, opt WatermarkOptions) (image.Image, error) {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	stdDraw.Draw(rgba, bounds, img, image.Point{}, stdDraw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(w.font)
	c.SetFontSize(opt.TextSize)
	c.SetClip(bounds)
	c.SetDst(rgba)

	// 设置文字颜色和透明度
	textColor := opt.TextColor
	textColor.A = uint8(float64(255) * opt.Opacity)
	c.SetSrc(image.NewUniform(textColor))

	// 计算文本大小
	opts := truetype.Options{
		Size: opt.TextSize,
	}
	face := truetype.NewFace(w.font, &opts)
	textWidth := font.MeasureString(face, opt.Text).Round()
	textHeight := int(opt.TextSize)

	if opt.Position == Tiled {
		// 平铺水印
		for y := opt.Margin; y < bounds.Dy(); y += int(opt.Spacing) + textHeight {
			for x := opt.Margin; x < bounds.Dx(); x += int(opt.Spacing) + textWidth {
				// 计算旋转后的位置
				rx, ry := rotatePoint(x, y, x+textWidth/2, y+textHeight/2, opt.Angle)
				pt := freetype.Pt(rx, ry+textHeight)
				_, err := c.DrawString(opt.Text, pt)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		// 单个水印
		x, y := calculatePosition(bounds.Dx(), bounds.Dy(), textWidth, textHeight, opt.Position, opt.Margin)
		// 计算旋转后的位置
		rx, ry := rotatePoint(x, y, x+textWidth/2, y+textHeight/2, opt.Angle)
		pt := freetype.Pt(rx, ry+textHeight)
		_, err := c.DrawString(opt.Text, pt)
		if err != nil {
			return nil, err
		}
	}

	return rgba, nil
}

// AddImageWatermark 添加图片水印
func (w *Watermarker) AddImageWatermark(baseImg image.Image, watermarkPath string, opt WatermarkOptions) (image.Image, error) {
	// 读取水印图片
	wmFile, err := os.Open(watermarkPath)
	if err != nil {
		return nil, err
	}
	defer wmFile.Close()

	watermark, _, err := image.Decode(wmFile)
	if err != nil {
		return nil, err
	}

	// 创建新的RGBA图像
	bounds := baseImg.Bounds()
	rgba := image.NewRGBA(bounds)
	stdDraw.Draw(rgba, bounds, baseImg, image.Point{}, stdDraw.Src)

	// 计算水印大小
	wmBounds := watermark.Bounds()
	wmWidth := int(float64(wmBounds.Dx()) * opt.Scale)
	wmHeight := int(float64(wmBounds.Dy()) * opt.Scale)

	// 创建缩放后的水印图片
	scaledWatermark := image.NewRGBA(image.Rect(0, 0, wmWidth, wmHeight))
	draw.BiLinear.Scale(scaledWatermark, scaledWatermark.Bounds(), watermark, wmBounds, draw.Over, nil)

	if opt.Position == Tiled {
		// 平铺水印
		for y := opt.Margin; y < bounds.Dy(); y += int(opt.Spacing) + wmHeight {
			for x := opt.Margin; x < bounds.Dx(); x += int(opt.Spacing) + wmWidth {
				// 计算旋转后的位置
				rx, ry := rotatePoint(x, y, x+wmWidth/2, y+wmHeight/2, opt.Angle)
				drawImageWithOpacity(rgba, scaledWatermark, rx, ry, opt.Opacity)
			}
		}
	} else {
		// 单个水印
		x, y := calculatePosition(bounds.Dx(), bounds.Dy(), wmWidth, wmHeight, opt.Position, opt.Margin)
		// 计算旋转后的位置
		rx, ry := rotatePoint(x, y, x+wmWidth/2, y+wmHeight/2, opt.Angle)
		drawImageWithOpacity(rgba, scaledWatermark, rx, ry, opt.Opacity)
	}

	return rgba, nil
}

// drawImageWithOpacity 绘制带透明度的图片
func drawImageWithOpacity(dst *image.RGBA, src image.Image, x, y int, opacity float64) {
	bounds := src.Bounds()
	for py := 0; py < bounds.Dy(); py++ {
		for px := 0; px < bounds.Dx(); px++ {
			srcColor := src.At(px, py)
			r, g, b, a := srcColor.RGBA()
			a = uint32(float64(a) * opacity)
			color := color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			}
			dst.Set(x+px, y+py, color)
		}
	}
}

// LoadImage 加载图片
func LoadImage(r io.Reader) (image.Image, error) {
	// 读取所有数据
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// 检测图片格式
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// 如果不是RGBA格式，转换为RGBA
	if _, ok := img.(*image.RGBA); !ok {
		bounds := img.Bounds()
		rgba := image.NewRGBA(bounds)
		stdDraw.Draw(rgba, bounds, img, bounds.Min, stdDraw.Src)
		return rgba, nil
	}

	return img, nil
}

// SaveImage 保存图片
func SaveImage(img image.Image, outputPath string, quality int) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	ext := filepath.Ext(outputPath)
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
	case ".png":
		return png.Encode(out, img)
	default:
		return jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
	}
}

// GetImagePreview returns a base64 encoded image preview
func GetImagePreview(filePath string) (string, error) {
	// Read the image file
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	// Encode to base64
	base64String := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imageData)
	return base64String, nil
}
