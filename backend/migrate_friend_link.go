package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	fs_util "github.com/kingwrcy/moments/util"
	"github.com/kingwrcy/moments/vo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func migrateFriendLink(log zerolog.Logger, cfg *vo.AppConfig) {
	if !fs_util.Exists(cfg.DB) {
		log.Debug().Msgf("原数据库不存在, 所以无需迁移友情链接")
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

	friendLinks, ok := config["friendLinks"].(string)
	if !ok {
		return
	}

	log.Info().Msgf("开始迁移友情链接, friendLinks: %s", friendLinks)

	for _, friendLink := range strings.Split(friendLinks, "\n") {
		if friendLink == "" {
			continue
		}

		itemList := strings.Split(friendLink, "|")
		if len(itemList) != 3 {
			continue
		}

		now := time.Now()
		name := itemList[0]
		url := itemList[1]
		icon := itemList[2]
		log.Info().Msgf("创建友情链接, name: %s, url: %s, icon: %s", name, url, icon)
		db.Exec("INSERT INTO Friend(name, url, icon, createdAt, updatedAt) VALUES(?, ?, ?, ?, ?)", name, url, icon, now, now)
	}

	log.Info().Msgf("迁移友情链接完成")

	config["friendLinks"] = nil
	content, err := json.Marshal(config)
	if err != nil {
		log.Error().Msgf("序列化新配置失败, err: %v", err)
	}

	db.Exec("UPDATE SysConfig SET content = ?", string(content))
}
