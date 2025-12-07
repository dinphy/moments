package vo

// AIPolishRequest AI润色请求
type AIPolishRequest struct {
	Content string `json:"content" binding:"required"` // 需要润色的内容
	Prompt  string `json:"prompt"`  // 自定义润色指令，为空则使用默认指令
}

// AIPolishResponse AI润色响应
type AIPolishResponse struct {
	OriginalContent string `json:"originalContent"` // 原始内容
	PolishedContent string `json:"polishedContent"` // 润色后的内容
}

// AIChatRequest AI对话请求
type AIChatRequest struct {
	Messages []AIChatMessage `json:"messages" binding:"required"` // 对话消息列表
}

// AIChatMessage AI对话消息
type AIChatMessage struct {
	Role    string `json:"role"`    // 角色：system, user, assistant
	Content string `json:"content"` // 消息内容
}

// AIChatResponse AI对话响应
type AIChatResponse struct {
	Message AIChatMessage `json:"message"` // AI回复消息
}
