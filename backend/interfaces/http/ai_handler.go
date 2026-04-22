package http

import (
	"net/http"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/infrastructure/mongodb"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	logRepo       *mongodb.AILogRepository
	geminiService *services.GeminiService
}

func NewAIHandler(logRepo *mongodb.AILogRepository, geminiService *services.GeminiService) *AIHandler {
	return &AIHandler{logRepo: logRepo, geminiService: geminiService}
}

type aiParseRequest struct {
	Message             string                   `json:"message" binding:"required"`
	ConversationHistory []map[string]interface{} `json:"conversation_history"`
}

type aiAction struct {
	Type   string                 `json:"type"`
	Fields map[string]interface{} `json:"fields"`
}

type aiParseResponse struct {
	ReplyText string    `json:"reply_text"`
	Action    *aiAction `json:"action,omitempty"`
}

// ParseCommand handles POST /manager/ai/parse
func (h *AIHandler) ParseCommand(c *gin.Context) {
	var req aiParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	_ = h.logRepo.Insert(ctx, &mongodb.AICommandLog{
		Message: req.Message,
		Role:    "user",
	})

	geminiResp, err := h.geminiService.ParseCommand(ctx, req.Message, req.ConversationHistory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service error"})
		return
	}

	resp := aiParseResponse{
		ReplyText: geminiResp.ReplyText,
	}
	if geminiResp.Action != nil {
		resp.Action = &aiAction{
			Type:   geminiResp.Action.Type,
			Fields: geminiResp.Action.Fields,
		}
	}

	actionType := ""
	if resp.Action != nil {
		actionType = resp.Action.Type
	}
	_ = h.logRepo.Insert(ctx, &mongodb.AICommandLog{
		Message:    resp.ReplyText,
		Role:       "agent",
		ActionType: actionType,
	})

	c.JSON(http.StatusOK, resp)
}

// GetHistory handles GET /manager/ai/history
func (h *AIHandler) GetHistory(c *gin.Context) {
	logs, err := h.logRepo.GetRecent(c.Request.Context(), 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
