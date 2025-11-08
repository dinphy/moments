package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// WebhookMessage 企业微信Webhook消息结构
type WebhookMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

// SendWebhookNotification 发送企业微信Webhook通知
func SendWebhookNotification(webhookURL, title, content string) error {
	if webhookURL == "" {
		return fmt.Errorf("企业微信Webhook URL不能为空")
	}

	// 构建消息内容
	message := WebhookMessage{
		MsgType: "markdown",
	}
	message.Markdown.Content = fmt.Sprintf("### %s\n\n%s", title, content)

	// 将消息转换为JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("构建企业微信消息失败: %v", err)
	}

	// 创建HTTP请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建企业微信请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送企业微信通知失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("企业微信通知响应异常, 状态码: %d, 内容: %s", resp.StatusCode, string(body))
	}

	log.Info().Msg("成功发送企业微信通知")
	return nil
}
