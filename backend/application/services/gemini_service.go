package services

import (
	"context"
	"encoding/json"
	"fmt"

	"cafe-pos/backend/infrastructure/mongodb"
	"google.golang.org/genai"
)

const defaultGeminiModel = "gemini-2.5-flash"

type GeminiService struct {
	shopSettingsRepo *mongodb.ShopSettingsRepository
	ingredientRepo   *mongodb.IngredientRepository
	expenseRepo      *mongodb.ExpenseRepository
}

func NewGeminiService(
	shopSettingsRepo *mongodb.ShopSettingsRepository,
	ingredientRepo *mongodb.IngredientRepository,
	expenseRepo *mongodb.ExpenseRepository,
) *GeminiService {
	return &GeminiService{
		shopSettingsRepo: shopSettingsRepo,
		ingredientRepo:   ingredientRepo,
		expenseRepo:      expenseRepo,
	}
}

type GeminiParseResponse struct {
	ReplyText string        `json:"reply_text"`
	Action    *GeminiAction `json:"action,omitempty"`
}

type GeminiAction struct {
	Type   string                 `json:"type"`
	Fields map[string]interface{} `json:"fields"`
}

// ParseCommand sends the user message + history to Gemini and returns a structured response.
// conversationHistory is a slice of {"role": "user"|"agent", "message": "..."} maps.
func (s *GeminiService) ParseCommand(ctx context.Context, message string, conversationHistory []map[string]interface{}) (*GeminiParseResponse, error) {
	// 1. Read settings
	shopSettings, err := s.shopSettingsRepo.GetSettings(ctx)
	if err != nil || shopSettings == nil {
		return noKeyResponse(), nil
	}
	if shopSettings.GeminiAPIKey == "" {
		return noKeyResponse(), nil
	}

	apiKey := shopSettings.GeminiAPIKey
	modelName := shopSettings.GeminiModel
	if modelName == "" {
		modelName = defaultGeminiModel
	}

	// 2. Fetch ingredients
	ingredients, err := s.ingredientRepo.FindAll(ctx)
	if err != nil {
		return geminiErrorResponse(), nil
	}
	type ingredientSummary struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Unit     string  `json:"unit"`
		Quantity float64 `json:"quantity"`
	}
	ingList := make([]ingredientSummary, 0, len(ingredients))
	for _, ing := range ingredients {
		ingList = append(ingList, ingredientSummary{
			ID:       ing.ID.Hex(),
			Name:     ing.Name,
			Unit:     string(ing.Unit),
			Quantity: ing.Quantity,
		})
	}
	ingJSON, _ := json.Marshal(ingList)

	// 3. Fetch expense categories
	categories, err := s.expenseRepo.GetCategories(ctx, "")
	if err != nil {
		return geminiErrorResponse(), nil
	}
	type categorySummary struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	catList := make([]categorySummary, 0, len(categories))
	for _, cat := range categories {
		catList = append(catList, categorySummary{
			ID:   cat.ID.Hex(),
			Name: cat.Name,
		})
	}
	catJSON, _ := json.Marshal(catList)

	// 4. Build system prompt
	systemPrompt := fmt.Sprintf(`Bạn là trợ lý quản lý quán cà phê. Nhiệm vụ của bạn là phân tích lệnh của người dùng và trả về JSON theo đúng schema sau. Không trả về gì ngoài JSON.

Schema:
{
  "reply_text": "<câu trả lời ngắn bằng tiếng Việt>",
  "action": {
    "type": "add_ingredient" | "restock_ingredient" | "add_expense",
    "fields": { ... }
  }
}

Nếu không rõ ý định, trả về action = null.

Fields theo từng action type:

add_ingredient:
  name (string), quantity (number), unit (string), cost_per_unit (number)

restock_ingredient:
  ingredient_id (string - từ danh sách bên dưới, hoặc "" nếu không tìm thấy),
  ingredient_name (string), current_stock (number), unit (string),
  quantity (number), cost_per_unit (number), money_type ("cash"|"transfer"), reason (string)

add_expense:
  description (string), amount (number), money_type ("cash"|"transfer"),
  category_id (string - từ danh sách bên dưới, hoặc "" nếu không tìm thấy),
  date (string YYYY-MM-DD, hoặc "" để dùng ngày hôm nay)

Danh sách nguyên liệu hiện có:
%s

Danh sách danh mục chi phí:
%s`, string(ingJSON), string(catJSON))

	// 5. Call Gemini SDK
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return geminiErrorResponse(), nil
	}

	// Build contents: history + current message
	var contents []*genai.Content
	for _, h := range conversationHistory {
		role, _ := h["role"].(string)
		msg, _ := h["message"].(string)
		if role == "" || msg == "" {
			continue
		}
		geminiRole := "user"
		if role == "agent" {
			geminiRole = "model"
		}
		contents = append(contents, &genai.Content{
			Role:  geminiRole,
			Parts: []*genai.Part{{Text: msg}},
		})
	}
	contents = append(contents, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: message}},
	})

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
		ResponseMIMEType: "application/json",
	})
	if err != nil {
		return geminiErrorResponse(), nil
	}

	// 6. Parse JSON response
	if resp == nil || len(resp.Candidates) == 0 {
		return geminiErrorResponse(), nil
	}
	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return geminiErrorResponse(), nil
	}
	rawText := candidate.Content.Parts[0].Text

	var parsed GeminiParseResponse
	if err := json.Unmarshal([]byte(rawText), &parsed); err != nil {
		return invalidJSONResponse(), nil
	}
	return &parsed, nil
}

func noKeyResponse() *GeminiParseResponse {
	return &GeminiParseResponse{
		ReplyText: "Chưa cấu hình Gemini API key. Vui lòng vào Cài đặt → AI Settings để thêm.",
	}
}

func geminiErrorResponse() *GeminiParseResponse {
	return &GeminiParseResponse{
		ReplyText: "Có lỗi kết nối AI. Vui lòng thử lại.",
	}
}

func invalidJSONResponse() *GeminiParseResponse {
	return &GeminiParseResponse{
		ReplyText: "AI trả về kết quả không hợp lệ. Vui lòng thử lại.",
	}
}
