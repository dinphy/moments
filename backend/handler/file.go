package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	fs_util "github.com/kingwrcy/moments/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type FileHandler struct {
	base BaseHandler
}

func NewFileHandler(injector do.Injector) *FileHandler {
	return &FileHandler{do.MustInvoke[BaseHandler](injector)}
}

// Upload godoc
//
//	@Tags		File
//	@Summary	上传图片
//	@Accept		json
//	@Produce	json
//	@Param		x-api-token	header	string	true	"登录TOKEN"
//	@Success	200
//	@Router		/api/file/upload [post]
func (f FileHandler) Upload(c echo.Context) error {
	var (
		result []string
	)

	form, err := c.MultipartForm()
	if err != nil {
		f.base.log.Error().Msgf("读取上传图片异常:%s", err)
		return FailRespWithMsg(c, Fail, "上传图片异常")
	}

	if err := os.MkdirAll(f.base.cfg.UploadDir, 0755); err != nil {
		f.base.log.Error().Msgf("创建父级目录异常:%s", err)
		return FailRespWithMsg(c, Fail, "创建父级目录异常")
	}

	files := form.File["files"]
	for _, file := range files {
		// 原始图片
		src, err := file.Open()
		if err != nil {
			f.base.log.Error().Msgf("打开上传图片异常:%s", err)
			return FailRespWithMsg(c, Fail, "上传图片异常")
		}
		defer src.Close()

		// 创建原始图片
		img_filename := strings.ReplaceAll(uuid.NewString(), "-", "")
		img_filepath := path.Join(f.base.cfg.UploadDir, img_filename)
		dst, err := os.Create(img_filepath)
		if err != nil {
			f.base.log.Error().Msgf("打开目标图片异常:%s", err)
			return FailRespWithMsg(c, Fail, "上传图片异常")
		}
		defer dst.Close()

		// 保存图片
		if _, err = io.Copy(dst, src); err != nil {
			f.base.log.Error().Msgf("复制图片异常:%s", err)
			return FailRespWithMsg(c, Fail, "上传图片异常")
		}

		// 生成并保存缩略图
		thumb_filename := img_filename + "_thumb"
		thumb_filepath := path.Join(f.base.cfg.UploadDir, thumb_filename)
		if err := CompressImage(f, img_filepath, thumb_filepath, 30); err != nil {
			f.base.log.Error().Msgf("压缩图片异常:%s", err)
		}

		result = append(result, "/upload/"+img_filename)
	}
	return SuccessResp(c, result)
}

// UploadImage godoc
//
//	@Tags		File
//	@Summary	上传图片
//	@Accept		json
//	@Produce	json
//	@Param		x-api-token	header	string	true	"登录TOKEN"
//	@Success	200
//	@Router		/api/file/uploadImage [post]
func (f FileHandler) UploadImage(c echo.Context) error {
	var (
		result []string
		images []db.Image
	)

	form, err := c.MultipartForm()
	if err != nil {
		f.base.log.Error().Msgf("读取上传图片异常:%s", err)
		return FailRespWithMsg(c, Fail, "上传图片异常")
	}

	if err := os.MkdirAll(f.base.cfg.UploadDir, 0755); err != nil {
		f.base.log.Error().Msgf("创建父级目录异常:%s", err)
		return FailRespWithMsg(c, Fail, "创建父级目录异常")
	}

	files := form.File["files"]
	for _, file := range files {
		// 原始图片
		src, err := file.Open()
		if err != nil {
			f.base.log.Error().Msgf("打开上传图片异常:%s", err)
			return FailRespWithMsg(c, Fail, "上传图片异常")
		}
		defer src.Close()

		// 计算图片hash
		hash, err := fs_util.CalHash(src)
		if err != nil {
			f.base.log.Error().Msgf("计算图片hash异常:%s", err)
			return FailRespWithMsg(c, Fail, "上传图片异常")
		}

		var image db.Image
		err = f.base.db.Where("hash = ?", hash).First(&image).Error
		// 如果图片不存在，则创建图片
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建原始图片
			img_filename := strings.ReplaceAll(uuid.NewString(), "-", "")
			img_filepath := path.Join(f.base.cfg.UploadDir, img_filename)
			dst, err := os.Create(img_filepath)
			if err != nil {
				f.base.log.Error().Msgf("打开目标图片异常:%s", err)
				return FailRespWithMsg(c, Fail, "上传图片异常")
			}
			defer dst.Close()

			// 重置文件指针到开头
			if seeker, ok := src.(io.Seeker); ok {
				if _, err := seeker.Seek(0, io.SeekStart); err != nil {
					f.base.log.Error().Msgf("重置文件指针异常:%s", err)
					return FailRespWithMsg(c, Fail, "上传图片异常")
				}
			}

			// 保存图片
			if _, err = io.Copy(dst, src); err != nil {
				f.base.log.Error().Msgf("复制图片异常:%s", err)
				return FailRespWithMsg(c, Fail, "上传图片异常")
			}

			// 生成并保存缩略图
			thumb_filename := img_filename + "_thumb"
			thumb_filepath := path.Join(f.base.cfg.UploadDir, thumb_filename)
			if err := CompressImage(f, img_filepath, thumb_filepath, 30); err != nil {
				f.base.log.Error().Msgf("压缩图片异常:%s", err)
			}

			// 保存图片信息
			image = db.Image{
				Hash: hash,
				Path: "/upload/" + img_filename,
			}
			images = append(images, image)
		}
		result = append(result, image.Path)
	}

	if err := f.base.db.CreateInBatches(images, 100).Error; err != nil {
		f.base.log.Error().Msgf("创建图片信息异常:%s", err)
		return FailRespWithMsg(c, Fail, "上传图片异常")
	}

	return SuccessResp(c, result)
}

// DeleteNoRelImage godoc
//
//	@Tags		File
//	@Summary	删除没有关联的图片
//	@Accept		json
//	@Produce	json
//	@Success	200
//	@Router		/api/file/deleteNoRelImage [post]
func (f FileHandler) DeleteNoRelImage(c echo.Context) error {
	var (
		data []struct {
			Id   int32  `json:"id,omitempty"`
			Path string `json:"path,omitempty"`
		}
		paths    []string
		imageIds []int32
	)

	// 查询没有关联的图片
	if err := f.base.db.Table("Image AS i").
		Joins("LEFT JOIN ImageRel AS ir ON i.id = ir.imageId").
		Select("i.id, i.path, COUNT(ir.memoId) AS relCount").
		Group("i.id").
		Having("relCount = 0").
		Find(&data).Error; err != nil {
		f.base.log.Error().Msgf("查询图片信息异常:%s", err)
		return FailRespWithMsg(c, Fail, "查询图片信息异常")
	}

	if len(data) == 0 {
		return SuccessResp(c, h{
			"num": 0,
		})
	}

	// 删除没有关联的图片文件
	for _, d := range data {
		img := d.Path
		if img == "" || !strings.HasPrefix(img, "/upload/") {
			continue
		}
		paths = append(paths, img)
		imageIds = append(imageIds, d.Id)

		img = strings.ReplaceAll(img, "/upload/", "")
		_ = os.Remove(filepath.Join(f.base.cfg.UploadDir, img))
		thumbImg := strings.ReplaceAll(img+"_thumb", "/upload/", "")
		_ = os.Remove(filepath.Join(f.base.cfg.UploadDir, thumbImg))
	}

	// 删除没有关联的图片信息
	if err := f.base.db.Delete(&db.Image{}, imageIds).Error; err != nil {
		f.base.log.Error().Msgf("删除图片信息异常:%s", err)
		return FailRespWithMsg(c, Fail, "删除图片信息异常")
	}

	f.base.log.Info().Msgf("删除本地图片: %v", paths)

	return SuccessResp(c, h{
		"num": len(paths),
	})
}

type PreSignedReq struct {
	ContentType string `json:"contentType,omitempty"` //图片mime类型
}

type s3PresignedResp struct {
	PreSignedUrl string `json:"preSignedUrl,omitempty"` //S3预签名上传地址
	ImageUrl     string `json:"imageUrl,omitempty"`     //实际的图片地址
}

// S3PreSigned godoc
//
//	@Tags		File
//	@Summary	S3预签名
//	@Accept		json
//	@Produce	json
//	@Param		object		body		PreSignedReq	true	"S3预签名"
//	@Param		x-api-token	header		string			true	"登录TOKEN"
//	@Success	200			{object}	s3PresignedResp
//	@Router		/api/file/s3PreSigned [post]
func (f FileHandler) S3PreSigned(c echo.Context) error {
	var (
		req         PreSignedReq
		sysConfig   db.SysConfig
		sysConfigVo vo.FullSysConfigVO
	)
	if err := c.Bind(&req); err != nil {
		return FailResp(c, ParamError)
	}

	if err := f.base.db.First(&sysConfig).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return FailResp(c, Fail)
	}

	if err := json.Unmarshal([]byte(sysConfig.Content), &sysConfigVo); err != nil {
		f.base.log.Error().Msgf("无法反序列化系统配置, %s", err)
		return FailRespWithMsg(c, Fail, err.Error())
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(sysConfigVo.S3.Region),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(
				func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{URL: sysConfigVo.S3.Endpoint}, nil
				},
			),
		),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				sysConfigVo.S3.AccessKey,
				sysConfigVo.S3.SecretKey,
				"",
			),
		),
	)
	if err != nil {
		f.base.log.Error().Msgf("无法加载SDK配置, %s", err)
		return FailRespWithMsg(c, Fail, err.Error())
	}

	client := s3.NewFromConfig(cfg)
	presignedClient := s3.NewPresignClient(client)

	key := fmt.Sprintf(
		"moments/%s/%s",
		time.Now().Format("2006/01/02"),
		strings.ReplaceAll(uuid.NewString(), "-", ""),
	)
	presignedResult, err := presignedClient.PresignPutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket:      aws.String(sysConfigVo.S3.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String(req.ContentType),
		},
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Minute * 5
		},
	)

	if err != nil {
		f.base.log.Error().Msgf("无法获取预签名URL, %s", err)
		return FailRespWithMsg(c, Fail, fmt.Sprintf("无法获取预签名URL, %s", err))
	}

	return SuccessResp(
		c,
		s3PresignedResp{
			PreSignedUrl: presignedResult.URL,
			ImageUrl:     fmt.Sprintf("%s/%s", sysConfigVo.S3.Domain, key),
		},
	)
}
