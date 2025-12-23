package handler

import (
	"os"

	"github.com/disintegration/imaging"
)

// 判断是否支持压缩
func SupportCompress(filename string) bool {
	_, err := imaging.FormatFromFilename(filename)
	return err == nil
}

// CompressImage 接收图片路径并进行压缩
// 返回错误信息
func CompressImage(f FileHandler, originImagePath string, thumbImagePath string, quality int) error {
	f.base.log.Debug().Msgf("start compressing image, originImagePath=%s, quality=%d", originImagePath, quality)

	originImage, err := imaging.Open(originImagePath, imaging.AutoOrientation(true))
	if err != nil {
		f.base.log.Error().Msgf("failed to open image: %v", err)
		return err
	}

	thumbImageFile, err := os.Create(thumbImagePath)
	if err != nil {
		f.base.log.Error().Msgf("failed to create thumb image: %v", err)
		return err
	}
	defer thumbImageFile.Close()

	compressedImage := imaging.Resize(originImage, 600, 0, imaging.Lanczos)
	err = imaging.Encode(thumbImageFile, compressedImage, imaging.JPEG)
	if err != nil {
		f.base.log.Error().Msgf("failed to save image: %v", err)
		os.Remove(thumbImagePath)
		return err
	}

	f.base.log.Debug().Msgf("compressing image finished, originImagePath=%s, quality=%d", originImagePath, quality)
	return nil
}
