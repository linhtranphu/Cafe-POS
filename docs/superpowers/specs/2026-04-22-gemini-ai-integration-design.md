# Gemini AI Integration Design

**Date:** 2026-04-22
**Status:** Approved
**Scope:** Replace the keyword-matching stub in `ai_handler.go` with a real Google Gemini API call. API key stored in MongoDB `shop_settings`, configurable via existing Settings UI.

---

## 1. Architecture

### New files
- `backend/services/gemini_service.go` — SDK client init, system prompt builder, Gemini call, JSON response parser

### Modified files
| File | Change |
|------|--------|
| `backend/domain/settings/shop_settings.go` | Add `GeminiAPIKey`, `GeminiModel` fields |
| `backend/interfaces/http/ai_handler.go` | Accept `*services.GeminiService`, call it instead of `stubParse` |
| `backend/main.go` | Init `GeminiService` with `shopSettingsRepo` + `ingredientRepo` + `expenseCategoryRepo`, inject into `AIHandler` |
| `backend/go.mod` / `go.sum` | Add `google.golang.org/genai` |
| `frontend/src/views/SettingsView.vue` | Add "AI Settings" card with `gemini_api_key` (password input) and `gemini_model` (text input) |

### No new env vars
The API key lives entirely in MongoDB `shop_settings`. No env var fallback.

---

## 2. ShopSettings Schema

Two new optional fields added to the `ShopSettings` struct:

```go
GeminiAPIKey string `bson:"gemini_api_key" json:"gemini_api_key"`
GeminiModel  string `bson:"gemini_model" json:"gemini_model"`
```

- No migration required — existing documents read missing fields as empty string.
- Empty `GeminiAPIKey` triggers a user-friendly fallback message.
- Empty `GeminiModel` defaults to `"gemini-2.5-flash"` inside `GeminiService`.

---

## 3. GeminiService — Request Flow

```
ParseCommand(ctx, message, conversationHistory)
  │
  ├─ 1. GetShopSettings() → read GeminiAPIKey, GeminiModel
  │      if GeminiAPIKey == "" → return friendly error message (no action)
  │
  ├─ 2. IngredientRepo.FindAll() → []{ id, name, unit, quantity }
  │
  ├─ 3. ExpenseCategoryRepo.FindAll() → []{ id, name }
  │
  ├─ 4. Build system prompt (Vietnamese):
  │      - Describe 3 action types + required fields
  │      - Embed ingredient list as JSON
  │      - Embed expense category list as JSON
  │      - Instruct: respond ONLY with JSON matching aiParseResponse schema
  │
  ├─ 5. Call Gemini SDK:
  │      model: GeminiModel (default "gemini-2.5-flash")
  │      response_mime_type: "application/json"
  │      contents: conversationHistory + current message
  │
  └─ 6. Unmarshal JSON → aiParseResponse{ reply_text, action? }
         on unmarshal error → return generic Vietnamese error message
```

---

## 4. System Prompt Design

```
Bạn là trợ lý quản lý quán cà phê. Nhiệm vụ của bạn là phân tích lệnh của người dùng
và trả về JSON theo đúng schema sau. Không trả về gì ngoài JSON.

Schema:
{
  "reply_text": "<câu trả lời ngắn bằng tiếng Việt>",
  "action": {          // null nếu không rõ ý định
    "type": "add_ingredient" | "restock_ingredient" | "add_expense",
    "fields": { ... }  // xem bên dưới
  }
}

Action fields:

add_ingredient:
  name (string), quantity (number), unit (string), cost_per_unit (number)

restock_ingredient:
  ingredient_id (string - từ danh sách bên dưới),
  ingredient_name (string), current_stock (number), unit (string),
  quantity (number), cost_per_unit (number), money_type ("cash"|"transfer"), reason (string)

add_expense:
  description (string), amount (number), money_type ("cash"|"transfer"),
  category_id (string - từ danh sách bên dưới), date (string YYYY-MM-DD)

Danh sách nguyên liệu hiện có:
<JSON array: [{id, name, unit, quantity}]>

Danh sách danh mục chi phí:
<JSON array: [{id, name}]>
```

---

## 5. Frontend — Settings UI Change

Add a new section card to `SettingsView.vue` inside the existing manager settings form:

```
┌─ AI Settings ──────────────────────────────┐
│ Gemini API Key  [••••••••••••••] (password) │
│ Gemini Model    [gemini-2.5-flash]          │
└─────────────────────────────────────────────┘
```

- Saved via the existing `saveShopSettings()` call (no new endpoint).
- `gemini_api_key` uses `type="password"` to avoid exposing the key on screen.
- `gemini_model` placeholder: `gemini-2.5-flash`.

---

## 6. Error Handling

| Scenario | Behaviour |
|----------|-----------|
| `GeminiAPIKey` not set | Return `reply_text`: "Chưa cấu hình Gemini API key. Vui lòng vào Cài đặt → AI Settings để thêm." — no action |
| Gemini API error (network, quota) | Return `reply_text`: "Có lỗi kết nối AI. Vui lòng thử lại." — no action |
| Gemini returns invalid JSON | Return `reply_text`: "AI trả về kết quả không hợp lệ. Vui lòng thử lại." — no action |
| `ingredient_id` not resolved | Gemini returns empty string; frontend `RestockActionCard` still shows, user edits manually |

---

## 7. Dependencies

- Go SDK: `google.golang.org/genai` (official Gemini Go SDK)
- Model: `gemini-2.5-flash` (default, overridable via DB setting)
- No new frontend dependencies
