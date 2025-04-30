package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/glebarez/sqlite"
	fs_util "github.com/kingwrcy/moments/util"
	"github.com/kingwrcy/moments/vo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func backupDatabase(log zerolog.Logger, cfg *vo.AppConfig) {
	if !fs_util.Exists(cfg.DB) {
		log.Debug().Msgf("原数据库不存在, 所以无需备份数据库")
		return
	}

	db, err := gorm.Open(sqlite.Open(cfg.DB))
	if err != nil {
		log.Fatal().Msgf("打开原数据库出错, err: %v", err)
		return
	}

	var sysConfig string
	db.Raw("SELECT content FROM SysConfig").Scan(&sysConfig)

	var config map[string]any
	err = json.Unmarshal([]byte(sysConfig), &config)
	if err != nil {
		log.Fatal().Msgf("反序列化原配置出错, err: %v", err)
		return
	}

	version, ok := config["version"].(string)
	if !ok {
		version = "unknown"
	}

	commitId, ok = config["commitId"].(string)
	if !ok {
		commitId = "unknown"
	}

	if version == cfg.Version && commitId == cfg.CommitId {
		log.Debug().Msgf("当前版本和上次启动时的版本一致, 所以无需备份数据库, version: %s, commitId: %s", version, commitId)
		return
	}

	log.Info().Msgf("开始备份数据库")

	sourceFile, err := os.Open(cfg.DB)
	if err != nil {
		log.Error().Msgf("打开原数据库文件出错, err: %v", err)
		return
	}
	defer sourceFile.Close()

	now := time.Now()
	datetime := now.Format("2006-01-02-15-04-05")
	destinationFileName := fmt.Sprintf("backup-%s-%s-%s.sqlite3", datetime, version, commitId)

	destinationFile, err := os.Create(path.Join(path.Dir(cfg.DB), destinationFileName))
	if err != nil {
		log.Error().Msgf("创建备份数据库文件失败, err: %v", err)
		return
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		log.Error().Msgf("复制数据库文件失败, err: %v", err)
		return
	}

	log.Info().Msgf("数据库备份完成")

	config["version"] = cfg.Version
	config["commitId"] = cfg.CommitId
	content, err := json.Marshal(config)
	if err != nil {
		log.Error().Msgf("序列化新配置失败, err: %v", err)
	}

	db.Exec("UPDATE SysConfig SET content = ?", string(content))

	cleanBackupFiles(path.Dir(cfg.DB))
}

func cleanBackupFiles(dir string) error {
	re := regexp.MustCompile(`^backup-\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}-dev-.{7}\.sqlite3$`)

	var files []os.FileInfo
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && re.MatchString(info.Name()) {
			files = append(files, info)
		}
		return nil
	})
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	for i, file := range files {
		if i >= 5 {
			filePath := filepath.Join(dir, file.Name())
			os.Remove(filePath)
		}
	}

	return nil
}
