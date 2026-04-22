package http

import (
	"net/http"
	"strings"

	"cafe-pos/backend/infrastructure/mongodb"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	logRepo *mongodb.AILogRepository
}

func NewAIHandler(logRepo *mongodb.AILogRepository) *AIHandler {
	return &AIHandler{logRepo: logRepo}
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
// Stub: returns a mock action based on keywords. Replace inner block with real AI call later.
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

	// --- STUB: replace this block with real AI call ---
	resp := stubParse(req.Message)
	// --- END STUB ---

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

func stubParse(message string) aiParseResponse {
	lower := strings.ToLower(message)

	if strings.Contains(lower, "chi phí") || strings.Contains(lower, "expense") {
		return aiParseResponse{
			ReplyText: "Đã hiểu! Xác nhận ghi nhận chi phí sau:",
			Action: &aiAction{
				Type: "add_expense",
				Fields: map[string]interface{}{
					"description": "Chi phí (stub - hãy sửa lại)",
					"amount":      0,
					"money_type":  "cash",
					"category_id": "",
					"date":        "",
				},
			},
		}
	}

	if strings.Contains(lower, "nhập kho") || strings.Contains(lower, "restock") {
		return aiParseResponse{
			ReplyText: "Đã tìm thấy nguyên liệu! Xác nhận nhập kho:",
			Action: &aiAction{
				Type: "restock_ingredient",
				Fields: map[string]interface{}{
					"ingredient_id":   "",
					"ingredient_name": "Nguyên liệu (stub)",
					"current_stock":   0,
					"unit":            "kg",
					"quantity":        0,
					"cost_per_unit":   0,
					"money_type":      "cash",
					"reason":          "",
				},
			},
		}
	}

	if strings.Contains(lower, "thêm") || strings.Contains(lower, "nguyên liệu") || strings.Contains(lower, "ingredient") {
		return aiParseResponse{
			ReplyText: "Đã hiểu! Xác nhận thêm nguyên liệu mới:",
			Action: &aiAction{
				Type: "add_ingredient",
				Fields: map[string]interface{}{
					"name":          "Nguyên liệu (stub - hãy sửa lại)",
					"quantity":      0,
					"unit":          "kg",
					"cost_per_unit": 0,
				},
			},
		}
	}

	return aiParseResponse{
		ReplyText: "Xin lỗi, tôi chưa hiểu lệnh này. Hiện tại tôi có thể: thêm nguyên liệu mới, nhập kho nguyên liệu, ghi nhận chi phí.",
	}
}
