# Gemini AI Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the keyword-matching stub in `ai_handler.go` with a real Google Gemini API call that interprets Vietnamese natural-language commands and returns structured JSON actions for ingredient and expense management.

**Architecture:** A new `GeminiService` in `backend/application/services/gemini_service.go` reads the API key + model from `ShopSettings` (MongoDB), builds a Vietnamese system prompt embedding the current ingredient list and expense categories, then calls the Gemini Go SDK with JSON mode. `AIHandler` is updated to call `GeminiService.ParseCommand()` instead of `stubParse()`. `SettingsView.vue` gains an "AI Settings" card so managers can configure the key.

**Tech Stack:** Go 1.25, `google.golang.org/genai` (Gemini Go SDK), Gin, MongoDB, Vue 3 + Tailwind

---

## File Map

| File | Action | Responsibility |
|------|--------|---------------|
| `backend/domain/settings/shop_settings.go` | Modify | Add `GeminiAPIKey` + `GeminiModel` fields |
| `backend/application/services/gemini_service.go` | Create | SDK init, system prompt build, Gemini call, JSON parse |
| `backend/interfaces/http/ai_handler.go` | Modify | Accept `*services.GeminiService`, call it instead of `stubParse` |
| `backend/main.go` | Modify | Init `GeminiService`, inject into `AIHandler` |
| `backend/go.mod` + `backend/go.sum` | Modify | Add `google.golang.org/genai` dependency |
| `frontend/src/views/SettingsView.vue` | Modify | Add "AI Settings" card with key + model inputs |

---

## Task 1: Add Gemini fields to ShopSettings

**Files:**
- Modify: `backend/domain/settings/shop_settings.go`

- [ ] **Step 1: Add two fields to the ShopSettings struct**

Open `backend/domain/settings/shop_settings.go`. After the `LabelHeight` field (line ~31), add:

```go
// AI Settings
GeminiAPIKey string `bson:"gemini_api_key" json:"gemini_api_key"`
GeminiModel  string `bson:"gemini_model" json:"gemini_model"`
```

The block should look like:

```go
LabelWidth          int    `bson:"label_width" json:"label_width"`   // mm
LabelHeight         int    `bson:"label_height" json:"label_height"` // mm

// AI Settings
GeminiAPIKey string `bson:"gemini_api_key" json:"gemini_api_key"`
GeminiModel  string `bson:"gemini_model" json:"gemini_model"`

CreatedAt time.Time `bson:"created_at" json:"created_at"`
```

- [ ] **Step 2: Verify build is clean**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend"
go build .
```

Expected: no output (clean build).

- [ ] **Step 3: Commit**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend"
git add domain/settings/shop_settings.go
git commit -m "feat(settings): add GeminiAPIKey and GeminiModel fields to ShopSettings"
```

---

## Task 2: Add Gemini Go SDK dependency

**Files:**
- Modify: `backend/go.mod`, `backend/go.sum`

- [ ] **Step 1: Add the SDK**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend"
go get google.golang.org/genai
```

Expected output: lines like `go: added google.golang.org/genai v...`

- [ ] **Step 2: Verify build still clean**

```bash
go build .
```

Expected: no output.

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: add google.golang.org/genai SDK dependency"
```

---

## Task 3: Create GeminiService

**Files:**
- Create: `backend/application/services/gemini_service.go`

- [ ] **Step 1: Create the file**

Create `backend/application/services/gemini_service.go` with the full implementation:

```go
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
		return errorResponse(), nil
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
		return errorResponse(), nil
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
		return errorResponse(), nil
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
		return errorResponse(), nil
	}

	// 6. Parse JSON response
	if resp == nil || len(resp.Candidates) == 0 {
		return errorResponse(), nil
	}
	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return errorResponse(), nil
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

func errorResponse() *GeminiParseResponse {
	return &GeminiParseResponse{
		ReplyText: "Có lỗi kết nối AI. Vui lòng thử lại.",
	}
}

func invalidJSONResponse() *GeminiParseResponse {
	return &GeminiParseResponse{
		ReplyText: "AI trả về kết quả không hợp lệ. Vui lòng thử lại.",
	}
}
```

- [ ] **Step 2: Verify build**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend"
go build .
```

Expected: no output.

- [ ] **Step 3: Commit**

```bash
git add application/services/gemini_service.go
git commit -m "feat(ai): create GeminiService with system prompt and SDK call"
```

---

## Task 4: Wire GeminiService into AIHandler

**Files:**
- Modify: `backend/interfaces/http/ai_handler.go`
- Modify: `backend/main.go`

- [ ] **Step 1: Replace ai_handler.go**

Replace `backend/interfaces/http/ai_handler.go` entirely with:

```go
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
```

- [ ] **Step 2: Update main.go**

In `backend/main.go`:

1. Delete the line (around line 247):
   ```go
   aiHandler := http.NewAIHandler(aiLogRepo)
   ```

2. After the `expenseRepo` block (around line 352), add:
   ```go
   geminiService := services.NewGeminiService(shopSettingsRepo, ingredientRepo, expenseRepo)
   aiHandler := http.NewAIHandler(aiLogRepo, geminiService)
   ```

- [ ] **Step 3: Verify build**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend"
go build .
```

Expected: no output.

- [ ] **Step 4: Commit**

```bash
git add interfaces/http/ai_handler.go main.go
git commit -m "feat(ai): wire GeminiService into AIHandler, remove stub"
```

---

## Task 5: Add AI Settings card to SettingsView

**Files:**
- Modify: `frontend/src/views/SettingsView.vue`

- [ ] **Step 1: Find the form object in SettingsView**

Open `frontend/src/views/SettingsView.vue`. Find the reactive form/settings object and note its variable name (likely `form` or `settings`). Find where the object is initialized from the API response (in `onMounted` or a load function).

- [ ] **Step 2: Add gemini fields to the form initialization**

In the object where settings fields are assigned from the API response, add:

```js
gemini_api_key: data.gemini_api_key || '',
gemini_model: data.gemini_model || '',
```

(Use the same pattern as existing fields like `print_bridge_url` or `shop_name` in the same block.)

- [ ] **Step 3: Add the AI Settings card to the template**

Find the last settings card in the template. After it, add:

```html
<!-- AI Settings -->
<div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
  <h3 class="text-base font-semibold text-gray-900 mb-4">🤖 AI Settings</h3>
  <div class="space-y-4">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Gemini API Key</label>
      <input
        v-model="form.gemini_api_key"
        type="password"
        placeholder="Nhập Gemini API key..."
        class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
    </div>
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Gemini Model</label>
      <input
        v-model="form.gemini_model"
        type="text"
        placeholder="gemini-2.5-flash"
        class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
    </div>
  </div>
</div>
```

> Replace `form.gemini_api_key` / `form.gemini_model` with the actual variable name found in Step 1 if different.

- [ ] **Step 4: Verify frontend build**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/frontend"
npm run build 2>&1 | tail -5
```

Expected: `✓ built in ...` with no errors.

- [ ] **Step 5: Commit**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS"
git add frontend/src/views/SettingsView.vue
git commit -m "feat(settings): add AI Settings card for Gemini API key and model"
```

---

## Task 6: Manual smoke test

- [ ] **Step 1: Start backend and frontend**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS"
bash RESTART_SERVERS.sh
```

- [ ] **Step 2: Test no-key fallback**

Navigate to `/#/manager/ai-command`, send any message.

Expected reply: "Chưa cấu hình Gemini API key. Vui lòng vào Cài đặt → AI Settings để thêm."

- [ ] **Step 3: Configure key and test real call**

Navigate to `/#/settings`, scroll to "AI Settings", enter a valid Gemini API key, save.

Then go to `/#/manager/ai-command` and send:

> "Nhập kho sữa tươi 10 lít giá 15000 một lít"

Expected: Gemini returns a `restock_ingredient` action card with ingredient name, quantity, unit, cost pre-filled.

- [ ] **Step 4: Test expense command**

Send:

> "Chi tiền điện tháng này 500000 tiền mặt"

Expected: `add_expense` action card with description, amount=500000, money_type=cash.

---

## Self-Review

- ✅ All 6 spec flow steps covered in Task 3 (settings → ingredients → categories → prompt → SDK call → JSON parse)
- ✅ `GeminiAPIKey`/`GeminiModel` fields added to `ShopSettings` (Task 1)
- ✅ Dependency added before service is written (Task 2 before Task 3)
- ✅ `aiHandler` moved after `expenseRepo` in `main.go` so all deps are in scope
- ✅ Frontend card uses `type="password"` for key, placeholder `gemini-2.5-flash` for model
- ✅ All 3 error scenarios handled: no key, API error, invalid JSON
- ✅ Default model `gemini-2.5-flash` when `GeminiModel` field is empty
- ✅ Type names consistent: `GeminiParseResponse`/`GeminiAction` in service → mapped to `aiParseResponse`/`aiAction` in handler
