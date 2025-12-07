package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type AIHandler struct {
	base BaseHandler
}

func NewAIHandler(injector do.Injector) *AIHandler {
	return &AIHandler{do.MustInvoke[BaseHandler](injector)}
}

// AIPolish godoc
//
//	@Tags			AI
//	@Summary		使用AI润色文本
//	@Description	使用AI润色文本内容
//	@Accept			json
//	@Produce		json
//	@Param			request	body	vo.AIPolishRequest	true	"润色请求"
//	@Param			x-api-token	header	string				true	"登录TOKEN"
//	@Success		200		{object}	vo.AIPolishResponse
//	@Router			/api/ai/polish [post]
func (a AIHandler) AIPolish(c echo.Context) error {
	var (
		request vo.AIPolishRequest
		config  db.SysConfig
	)

	// 检查用户是否登录
	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil {
		return FailRespWithMsg(c, Fail, "需要先登录")
	}

	// 绑定请求数据
	if err := c.Bind(&request); err != nil {
		return FailResp(c, ParamError)
	}

	// 检查内容是否为空
	if strings.TrimSpace(request.Content) == "" {
		return FailRespWithMsg(c, ParamError, "内容不能为空")
	}

	// 获取系统配置
	if err := a.base.db.First(&config).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return FailRespWithMsg(c, Fail, "系统配置不存在")
	}

	// 解析系统配置
	var sysConfig vo.FullSysConfigVO
	if err := json.Unmarshal([]byte(config.Content), &sysConfig); err != nil {
		return FailRespWithMsg(c, Fail, "系统配置解析失败")
	}

	// 检查AI功能是否启用
	if !sysConfig.EnableAI {
		return FailRespWithMsg(c, Fail, "AI润色功能未启用")
	}

	// 检查API Key是否配置
	if sysConfig.AIApiKey == "" {
		return FailRespWithMsg(c, Fail, "AI API Key未配置")
	}

	// 设置默认API URL
	apiURL := sysConfig.AIApiUrl
	if apiURL == "" {
		apiURL = "https://open.bigmodel.cn/api/paas/v4/chat/completions" // 智谱AI默认API地址
	}

	// 构建系统提示词
	systemPrompt := "你是一个专业的内容润色助手，负责优化用户输入的文本，使其更加流畅、生动、有感染力，同时保持原意不变。请直接返回润色后的文本，不要添加任何解释或说明。"
	if strings.TrimSpace(request.Prompt) != "" {
		systemPrompt = request.Prompt
	}

	// 构建请求体
	reqBody := map[string]interface{}{
		"model": "glm-4-flash", // 使用智谱AI的模型
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("请润色以下内容：\n%s", request.Content),
			},
		},
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return FailRespWithMsg(c, Fail, "构建请求失败")
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return FailRespWithMsg(c, Fail, "创建请求失败")
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sysConfig.AIApiKey))

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return FailRespWithMsg(c, Fail, "请求AI服务失败")
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return FailRespWithMsg(c, Fail, "读取响应失败")
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return FailRespWithMsg(c, Fail, fmt.Sprintf("AI服务返回错误: %s", string(respBody)))
	}

	// 解析响应
	var aiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &aiResponse); err != nil {
		return FailRespWithMsg(c, Fail, "解析AI响应失败")
	}

	// 检查AI返回的错误
	if aiResponse.Error != nil {
		return FailRespWithMsg(c, Fail, fmt.Sprintf("AI服务错误: %s", aiResponse.Error.Message))
	}

	// 检查是否有润色结果
	if len(aiResponse.Choices) == 0 || aiResponse.Choices[0].Message.Content == "" {
		return FailRespWithMsg(c, Fail, "AI未返回润色结果")
	}

	// 返回润色结果
	return SuccessResp(c, vo.AIPolishResponse{
		OriginalContent: request.Content,
		PolishedContent: strings.TrimSpace(aiResponse.Choices[0].Message.Content),
	})
}

// AIChat godoc
//
//	@Tags			AI
//	@Summary		AI对话
//	@Description	与AI进行对话交互
//	@Accept			json
//	@Produce		json
//	@Param			request	body	vo.AIChatRequest	true	"对话请求"
//	@Param			x-api-token	header	string				true	"登录TOKEN"
//	@Success		200		{object}	vo.AIChatResponse
//	@Router			/api/ai/chat [post]
func (a AIHandler) AIChat(c echo.Context) error {
	var (
		request vo.AIChatRequest
		config  db.SysConfig
	)

	// 检查用户是否登录
	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil {
		return FailRespWithMsg(c, Fail, "需要先登录")
	}

	// 绑定请求数据
	if err := c.Bind(&request); err != nil {
		return FailResp(c, ParamError)
	}

	// 检查消息列表是否为空
	if len(request.Messages) == 0 {
		return FailRespWithMsg(c, ParamError, "消息列表不能为空")
	}

	// 获取系统配置
	if err := a.base.db.First(&config).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return FailRespWithMsg(c, Fail, "系统配置不存在")
	}

	// 解析系统配置
	var sysConfig vo.FullSysConfigVO
	if err := json.Unmarshal([]byte(config.Content), &sysConfig); err != nil {
		return FailRespWithMsg(c, Fail, "系统配置解析失败")
	}

	// 检查AI功能是否启用
	if !sysConfig.EnableAI {
		return FailRespWithMsg(c, Fail, "AI功能未启用")
	}

	// 检查API Key是否配置
	if sysConfig.AIApiKey == "" {
		return FailRespWithMsg(c, Fail, "AI API Key未配置")
	}

	// 设置默认API URL
	apiURL := sysConfig.AIApiUrl
	if apiURL == "" {
		apiURL = "https://open.bigmodel.cn/api/paas/v4/chat/completions" // 智谱AI默认API地址
	}

	// 构建请求体
	messages := make([]map[string]interface{}, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	reqBody := map[string]interface{}{
		"model":    "glm-4-flash", // 使用智谱AI的模型
		"messages": messages,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return FailRespWithMsg(c, Fail, "构建请求失败")
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return FailRespWithMsg(c, Fail, "创建请求失败")
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sysConfig.AIApiKey))

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return FailRespWithMsg(c, Fail, "请求AI服务失败")
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return FailRespWithMsg(c, Fail, "读取响应失败")
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return FailRespWithMsg(c, Fail, fmt.Sprintf("AI服务返回错误: %s", string(respBody)))
	}

	// 解析响应
	var aiResponse struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &aiResponse); err != nil {
		return FailRespWithMsg(c, Fail, "解析AI响应失败")
	}

	// 检查AI返回的错误
	if aiResponse.Error != nil {
		return FailRespWithMsg(c, Fail, fmt.Sprintf("AI服务错误: %s", aiResponse.Error.Message))
	}

	// 检查是否有回复结果
	if len(aiResponse.Choices) == 0 || aiResponse.Choices[0].Message.Content == "" {
		return FailRespWithMsg(c, Fail, "AI未返回回复结果")
	}

	// 返回对话结果
	return SuccessResp(c, vo.AIChatResponse{
		Message: vo.AIChatMessage{
			Role:    aiResponse.Choices[0].Message.Role,
			Content: strings.TrimSpace(aiResponse.Choices[0].Message.Content),
		},
	})
}
