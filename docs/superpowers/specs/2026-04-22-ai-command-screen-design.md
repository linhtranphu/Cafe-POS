# AI Command Screen — Design Spec
**Date:** 2026-04-22
**Status:** Approved

---

## Overview

An admin-only chat interface at `/admin/ai-command` that lets the admin type natural language commands to add ingredients or record expenses. The AI parses the intent, extracts structured fields, and presents an inline action card for review. The admin confirms before any data is written.

---

## Architecture

### Flow

```
Admin types command
  → POST /api/ai/parse  (new, AI-only)
  → returns { reply_text, action: { type, fields } | null }
  → frontend shows agent reply + action card
  → admin confirms
  → frontend calls existing API (ingredient or expense)
```

### Backend — new endpoint only

**`POST /api/ai/parse`**
- File: `backend/interfaces/http/ai_handler.go`
- Request: `{ message: string, conversation_history: [...] }`
- Response: `{ reply_text: string, action: { type, fields } | null }`
- Action types: `"add_ingredient"` | `"restock_ingredient"` | `"add_expense"`
- Does zero CRUD — parsing only
- Stubbed initially; AI provider wired in later

**`GET /api/ai/history`**
- Returns paginated `ai_command_logs` entries (shared, no per-user scoping)

**New DB collection: `ai_command_logs`**
```json
{
  "message": "string",
  "role": "user | agent",
  "action_type": "add_ingredient | restock_ingredient | add_expense | null",
  "timestamp": "ISODate"
}
```

### Frontend — new files

| File | Purpose |
|---|---|
| `src/views/AICommandView.vue` | Main page — chat thread + sticky input bar |
| `src/components/ai/IngredientActionCard.vue` | Create new ingredient card |
| `src/components/ai/RestockActionCard.vue` | Adjust stock (stock-in) card |
| `src/components/ai/ExpenseActionCard.vue` | Record expense card |

### Frontend — reused

- Existing ingredient service (`POST /api/ingredients`)
- Existing restock endpoint (`POST /api/ingredients/:id/restock/from-fund`)
- Existing expense service (`POST /api/expenses/from-fund`)
- Existing fund store — loaded on page open to resolve fund name → fund_id

---

## UI

**Layout:** Single-column chat thread, sticky input bar at bottom.

- User messages: right-aligned, blue bubble
- Agent messages: left-aligned, white bubble with border
- Action cards: appear inline below the agent's text reply, left-aligned

**Conversation history:** Loaded from `GET /api/ai/history` on page open. Shared across all admin sessions (no per-user scoping).

---

## Action Card Definitions

### Create Ingredient Card (`add_ingredient`)
Calls: `POST /api/ingredients`

| Field | Type | Notes |
|---|---|---|
| Tên | text | Required |
| Số lượng | number + unit | Required |
| Đơn giá | number | Required |

### Adjust Stock Card (`restock_ingredient`)
Calls: `POST /api/ingredients/:id/restock/from-fund`

| Field | Type | Notes |
|---|---|---|
| Tên | text (locked) | Matched from existing ingredient |
| Tồn kho hiện tại | read-only | Shown for context |
| Số lượng nhập | number | Required |
| Đơn giá | number | cost_per_unit, required |
| Phương thức | cash / transfer | Defaults to cash |
| Lý do | text | Optional |

### Expense Card (`add_expense`)
Calls: `POST /api/expenses/from-fund`

| Field | Type | Notes |
|---|---|---|
| Mô tả | text | Required |
| Số tiền | number | Required |
| Phương thức | cash / transfer | Required |
| Loại quỹ | select from fund list | Resolved by name match |
| Ngày | date | Defaults to today |

---

## Ingredient Existence Check

Before returning an action card, `ai_handler.go` checks if the ingredient name already exists:

| Scenario | Result |
|---|---|
| Not found | `add_ingredient` card |
| Exact match found | `restock_ingredient` card with current stock shown |
| Multiple matches | No card — agent lists matches and asks to clarify |

---

## Error Handling

| Scenario | Behavior |
|---|---|
| AI returns incomplete fields | Agent asks clarifying question in chat, no card shown |
| Fund name not matched | Agent lists available funds, asks to clarify |
| Multiple ingredient name matches | Agent lists matches, asks to clarify |
| Confirm → API save error | Error shown below card buttons, card stays open for retry |
| Unknown intent | Agent replies with supported commands list |
| Stub mode (no AI wired) | Handler returns hardcoded mock response |

---

## Routing & Access

- Route: `/admin/ai-command`
- Guard: admin role only (same as other admin routes)
- Added to admin navigation

---

## Out of Scope (v1)

- Per-user conversation history
- Query/reporting commands ("revenue today?")
- Menu item management
- Editing or deleting existing records via AI
