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
	"slices"
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
//	@Summary	上传文件
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
		f.base.log.Error().Msgf("读取上传文件异常: %v", err)
		return FailRespWithMsg(c, Fail, "上传文件异常")
	}

	if err := os.MkdirAll(f.base.cfg.UploadDir, 0755); err != nil {
		f.base.log.Error().Msgf("创建父级目录异常: %v", err)
		return FailRespWithMsg(c, Fail, "上传文件异常")
	}

	files := form.File["files"]
	for _, file := range files {
		// 原始文件数据
		reader, err := file.Open()
		if err != nil {
			f.base.log.Error().Msgf("打开上传文件异常: %v", err)
			return FailRespWithMsg(c, Fail, "上传文件异常")
		}
		defer reader.Close()

		// 计算文件 hash
		sha256, err := fs_util.Sha256(reader)
		if err != nil {
			f.base.log.Error().Msgf("计算文件 hash 异常: %v", err)
			return FailRespWithMsg(c, Fail, "上传文件异常")
		}

		// 计算文件后缀
		ext := filepath.Ext(file.Filename)

		// 计算文件本地路径
		filename := fmt.Sprintf("%s%s", sha256, ext)
		filePath := path.Join(f.base.cfg.UploadDir, filename)

		// 添加到结果中
		result = append(result, "/upload/"+filename)

		// 如果文件存在，则跳过保存操作
		if fs_util.Exists(filePath) {
			continue
		}

		// 创建原始文件
		dst, err := os.Create(filePath)
		if err != nil {
			f.base.log.Error().Msgf("打开目标文件异常: %v", err)
			return FailRespWithMsg(c, Fail, "上传文件异常")
		}
		defer dst.Close()

		// 重置文件指针到开头
		if seeker, ok := reader.(io.Seeker); ok {
			if _, err := seeker.Seek(0, io.SeekStart); err != nil {
				f.base.log.Error().Msgf("重置文件指针异常: %v", err)
				return FailRespWithMsg(c, Fail, "上传文件异常")
			}
		}

		// 保存文件数据
		if _, err = io.Copy(dst, reader); err != nil {
			f.base.log.Error().Msgf("复制文件异常: %v", err)
			return FailRespWithMsg(c, Fail, "上传文件异常")
		}

		// 生成并保存缩略图文件
		if SupportCompress(filename) {
			thumb_filename := fmt.Sprintf("%s_thumb%s", sha256, ext)
			thumb_filepath := path.Join(f.base.cfg.UploadDir, thumb_filename)
			if err := CompressImage(f, filePath, thumb_filepath, 30); err != nil {
				f.base.log.Error().Msgf("压缩文件异常: %v", err)
			}
		}
	}

	return SuccessResp(c, result)
}

func (f FileHandler) Exist(c echo.Context) error {
	filename := c.QueryParam("filename")
	if strings.Contains(filename, "..") || strings.HasPrefix(filename, "/") {
		return FailRespWithMsg(c, Fail, "文件名包含非法字符")
	}

	filePath := filepath.Join(f.base.cfg.UploadDir, filename)
	return SuccessResp(c, h{
		"exist": fs_util.Exists(filePath),
		"path":  "/upload/" + filename,
	})
}

// Clean godoc
//
//	@Tags		File
//	@Summary	将没有关联的文件移动到 {uploadDir}/removed 目录下
//	@Accept		json
//	@Produce	json
//	@Success	200
//	@Router		/api/file/clean [post]
func (f FileHandler) Clean(c echo.Context) error {
	var (
		sysConfig   db.SysConfig
		sysConfigVo vo.FullSysConfigVO
	)

	if err := f.base.db.First(&sysConfig).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return FailResp(c, Fail)
	}

	if err := json.Unmarshal([]byte(sysConfig.Content), &sysConfigVo); err != nil {
		f.base.log.Error().Msgf("无法反序列化系统配置, %s", err)
		return FailRespWithMsg(c, Fail, err.Error())
	}

	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil || currentUser.Id != 1 {
		return FailRespWithMsg(c, Fail, "需要先登录")
	}

	uploadDir := f.base.cfg.UploadDir
	removedDir := filepath.Join(uploadDir, "removed")
	if err := os.MkdirAll(removedDir, 0755); err != nil {
		f.base.log.Error().Msgf("创建删除文件目录异常: %v", err)
		return FailRespWithMsg(c, Fail, "创建删除文件目录异常")
	}

	// uploaded files
	uploadedFiles, err := fs_util.GetFileList(uploadDir)
	if err != nil {
		f.base.log.Error().Msgf("获取上传文件列表异常: %v", err)
		return FailRespWithMsg(c, Fail, "获取上传文件列表异常")
	}

	// add prefix /upload/ to uploaded files
	for i, file := range uploadedFiles {
		uploadedFiles[i] = fmt.Sprintf("/upload/%s", file)
	}

	// if no uploaded files, return success
	if len(uploadedFiles) == 0 {
		return SuccessResp(c, h{
			"num": 0,
		})
	}

	// used files
	usedFiles := make([]string, 0)

	// memo imgs
	memoImgFiles := make([]string, 0)
	f.base.db.Model(&db.Memo{}).Where("imgs is not null and imgs != ''").Pluck("imgs", &memoImgFiles)
	for _, file := range memoImgFiles {
		files := strings.Split(file, ",")
		for _, file := range files {
			if strings.HasPrefix(file, "/upload/") {
				usedFiles = append(usedFiles, file)

				// add thumb filename to used files
				ext := filepath.Ext(file)
				filenameWithoutExt := strings.TrimSuffix(file, ext)
				thumbFilename := fmt.Sprintf("%s_thumb%s", filenameWithoutExt, ext)
				usedFiles = append(usedFiles, thumbFilename)
			}
		}
	}

	// memo video
	memoVideoFiles := make([]string, 0)
	f.base.db.Model(&db.Memo{}).Where("ext->>'video.value' is not null and ext->>'video.value' != ''").Pluck("ext->>'video.value'", &memoVideoFiles)
	for _, file := range memoVideoFiles {
		if strings.HasPrefix(file, "/upload/") {
			usedFiles = append(usedFiles, file)
		}
	}

	// sys config favicon
	if sysConfigVo.Favicon != "" && strings.HasPrefix(sysConfigVo.Favicon, "/upload/") {
		usedFiles = append(usedFiles, sysConfigVo.Favicon)
	}

	// user avatar url
	userAvatarFiles := make([]string, 0)
	f.base.db.Model(&db.User{}).Where("avatarUrl is not null and avatarUrl != ''").Pluck("avatarUrl", &userAvatarFiles)
	for _, file := range userAvatarFiles {
		if strings.HasPrefix(file, "/upload/") {
			usedFiles = append(usedFiles, file)
		}
	}

	// user cover url
	userCoverFiles := make([]string, 0)
	f.base.db.Model(&db.User{}).Where("coverUrl is not null and coverUrl != ''").Pluck("coverUrl", &userCoverFiles)
	for _, file := range userCoverFiles {
		if strings.HasPrefix(file, "/upload/") {
			usedFiles = append(usedFiles, file)
		}
	}

	// unused files
	unusedFiles := make([]string, 0)
	for _, file := range uploadedFiles {
		if !slices.Contains(usedFiles, file) {
			unusedFiles = append(unusedFiles, file)
		}
	}

	// move unused files to removed dir
	for _, filePath := range unusedFiles {
		if filePath == "" || !strings.HasPrefix(filePath, "/upload/") {
			continue
		}

		filename := strings.TrimPrefix(filePath, "/upload/")
		os.Rename(filepath.Join(uploadDir, filename), filepath.Join(removedDir, filename))
	}

	return SuccessResp(c, h{
		"num": len(unusedFiles),
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
