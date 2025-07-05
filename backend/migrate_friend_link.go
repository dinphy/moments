package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func migrateFriendLink(tx *gorm.DB, log zerolog.Logger) {
	var sysConfig string
	tx.Raw("SELECT content FROM SysConfig").Scan(&sysConfig)

	var config map[string]any
	err := json.Unmarshal([]byte(sysConfig), &config)
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
		tx.Exec("INSERT INTO Friend(name, url, icon, createdAt, updatedAt) VALUES(?, ?, ?, ?, ?)", name, url, icon, now, now)
	}

	log.Info().Msgf("迁移友情链接完成")

	config["friendLinks"] = nil
	content, err := json.Marshal(config)
	if err != nil {
		log.Error().Msgf("序列化新配置失败, err: %v", err)
	}

	tx.Exec("UPDATE SysConfig SET content = ?", string(content))
}
