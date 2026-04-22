# AI Command Screen Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build an admin-only chat screen at `/manager/ai-command` where the admin types natural language commands to add ingredients or record expenses, and the AI returns inline action cards for review before saving.

**Architecture:** A new `POST /api/ai/parse` backend endpoint (stub initially) parses intent and returns structured action fields. The frontend renders inline action cards in a chat thread. On confirm, the frontend calls existing APIs (`/manager/ingredients`, `/manager/ingredients/:id/restock/from-fund`, `/manager/expenses/from-fund`). Conversation history is persisted in a new `ai_command_logs` MongoDB collection.

**Tech Stack:** Go (Gin, MongoDB), Vue 3 (Composition API, local state only), Tailwind CSS, existing `api.js` axios instance.

---

## File Map

**New backend files:**
- `backend/infrastructure/mongodb/ai_log_repository.go` — `ai_command_logs` collection CRUD
- `backend/interfaces/http/ai_handler.go` — `POST /api/ai/parse`, `GET /api/ai/history`

**Modified backend files:**
- `backend/main.go` — init `AILogRepository`, init `AIHandler`, register routes under manager group

**New frontend files:**
- `src/services/aiService.js` — `parseCommand()`, `getHistory()`
- `src/components/ai/FieldRow.vue` — shared field display row
- `src/components/ai/IngredientActionCard.vue` — create new ingredient card
- `src/components/ai/RestockActionCard.vue` — adjust stock (stock-in) card
- `src/components/ai/ExpenseActionCard.vue` — expense from-fund card
- `src/views/AICommandView.vue` — chat thread + sticky input bar

**Modified frontend files:**
- `src/router/index.js` — add `/manager/ai-command` route
- `src/components/Navigation.vue` — add AI Command nav link for manager role

---

## Task 1: MongoDB log repository

**Files:**
- Create: `backend/infrastructure/mongodb/ai_log_repository.go`

- [ ] **Step 1: Create the repository file**

```go
package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AICommandLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Message    string             `bson:"message" json:"message"`
	Role       string             `bson:"role" json:"role"` // "user" or "agent"
	ActionType string             `bson:"action_type,omitempty" json:"action_type,omitempty"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
}

type AILogRepository struct {
	collection *mongo.Collection
}

func NewAILogRepository(db *mongo.Database) *AILogRepository {
	return &AILogRepository{
		collection: db.Collection("ai_command_logs"),
	}
}

func (r *AILogRepository) Insert(ctx context.Context, log *AICommandLog) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	log.ID = primitive.NewObjectID()
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}
	_, err := r.collection.InsertOne(ctx, log)
	return err
}

func (r *AILogRepository) GetRecent(ctx context.Context, limit int) ([]AICommandLog, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var logs []AICommandLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
```

- [ ] **Step 2: Verify it compiles**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend" && go build ./...
```
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add backend/infrastructure/mongodb/ai_log_repository.go
git commit -m "feat(ai): add AILogRepository for ai_command_logs collection"
```

---

## Task 2: AI handler (stubbed)

**Files:**
- Create: `backend/interfaces/http/ai_handler.go`

- [ ] **Step 1: Create the handler**

```go
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
```

- [ ] **Step 2: Verify compilation**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend" && go build ./...
```
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add backend/interfaces/http/ai_handler.go
git commit -m "feat(ai): add stubbed AI handler with parse and history endpoints"
```

---

## Task 3: Register routes in main.go

**Files:**
- Modify: `backend/main.go`

- [ ] **Step 1: Add AILogRepository init**

After the existing repo initialization block (near line 50), add:
```go
aiLogRepo := mongodb.NewAILogRepository(db)
```

- [ ] **Step 2: Add AIHandler init**

After the existing handler initialization block, add:
```go
aiHandler := http.NewAIHandler(aiLogRepo)
```

- [ ] **Step 3: Add routes inside the manager group**

After the last existing manager route, add:
```go
manager.POST("/ai/parse", aiHandler.ParseCommand)
manager.GET("/ai/history", aiHandler.GetHistory)
```

- [ ] **Step 4: Verify compilation**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend" && go build ./...
```
Expected: no errors

- [ ] **Step 5: Commit**

```bash
git add backend/main.go
git commit -m "feat(ai): register AI parse and history routes under manager group"
```

---

## Task 4: Frontend AI service

**Files:**
- Create: `frontend/src/services/aiService.js`

- [ ] **Step 1: Create the service**

```js
import api from './api'

export const aiService = {
  async parseCommand(message, conversationHistory = []) {
    const response = await api.post('/manager/ai/parse', {
      message,
      conversation_history: conversationHistory,
    })
    return response.data // { reply_text, action: { type, fields } | null }
  },

  async getHistory() {
    const response = await api.get('/manager/ai/history')
    return response.data.logs // AICommandLog[]
  },
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/services/aiService.js
git commit -m "feat(ai): add frontend aiService for parse and history"
```

---

## Task 5: FieldRow shared component + IngredientActionCard

**Files:**
- Create: `frontend/src/components/ai/FieldRow.vue`
- Create: `frontend/src/components/ai/IngredientActionCard.vue`

- [ ] **Step 1: Create FieldRow.vue**

```vue
<template>
  <div class="flex justify-between items-center text-sm py-1">
    <span class="text-gray-500">{{ label }}</span>
    <span :class="highlight ? 'font-semibold text-green-700' : 'font-medium text-gray-900'">{{ value }}</span>
  </div>
</template>
<script setup>
defineProps({ label: String, value: [String, Number], highlight: Boolean })
</script>
```

- [ ] **Step 2: Create IngredientActionCard.vue**

```vue
<template>
  <div class="border-2 border-amber-400 rounded-2xl overflow-hidden bg-white shadow-sm w-full max-w-sm">
    <div class="bg-amber-50 px-4 py-3 flex items-center gap-2 border-b border-amber-200">
      <span class="text-base">📦</span>
      <span class="font-semibold text-sm text-amber-900">Thêm nguyên liệu mới</span>
    </div>
    <div class="px-4 py-3 flex flex-col gap-2">
      <div v-if="!editing">
        <FieldRow label="Tên" :value="fields.name" />
        <FieldRow label="Số lượng" :value="`${fields.quantity} ${fields.unit}`" />
        <FieldRow label="Đơn giá" :value="formatPrice(fields.cost_per_unit)" />
        <FieldRow label="Tổng giá trị" :value="formatPrice(fields.quantity * fields.cost_per_unit)" :highlight="true" />
      </div>
      <div v-else class="flex flex-col gap-2">
        <label class="text-xs text-gray-500">Tên</label>
        <input v-model="edit.name" class="border rounded-lg px-3 py-2 text-sm w-full" />
        <label class="text-xs text-gray-500">Số lượng</label>
        <div class="flex gap-2">
          <input v-model.number="edit.quantity" type="number" min="0" class="border rounded-lg px-3 py-2 text-sm flex-1" />
          <input v-model="edit.unit" class="border rounded-lg px-3 py-2 text-sm w-20" placeholder="đơn vị" />
        </div>
        <label class="text-xs text-gray-500">Đơn giá (VNĐ)</label>
        <input v-model.number="edit.cost_per_unit" type="number" min="0" class="border rounded-lg px-3 py-2 text-sm w-full" />
      </div>
      <p v-if="error" class="text-xs text-red-600 mt-1">{{ error }}</p>
    </div>
    <div v-if="!confirmed" class="flex border-t border-gray-100">
      <button @click="toggleEdit" class="flex-1 py-3 text-sm text-gray-600 font-medium border-r border-gray-100 active:bg-gray-50">
        {{ editing ? 'Xem lại' : '✏️ Sửa' }}
      </button>
      <button @click="confirm" :disabled="loading" class="flex-1 py-3 text-sm text-white font-semibold bg-amber-500 active:bg-amber-600 disabled:bg-gray-300">
        {{ loading ? 'Đang lưu...' : '✓ Xác nhận' }}
      </button>
    </div>
    <div v-else class="px-4 py-3 text-sm text-green-700 font-medium text-center bg-green-50">✅ Đã thêm thành công</div>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
import { ingredientService } from '../../services/ingredient'
import FieldRow from './FieldRow.vue'

const props = defineProps({ fields: { type: Object, required: true } })
const emit = defineEmits(['confirmed'])

const editing = ref(false)
const loading = ref(false)
const confirmed = ref(false)
const error = ref('')
const edit = reactive({ ...props.fields })

function toggleEdit() {
  if (!editing.value) Object.assign(edit, props.fields)
  editing.value = !editing.value
}
function formatPrice(v) {
  return new Intl.NumberFormat('vi-VN').format(v || 0) + 'đ'
}
async function confirm() {
  error.value = ''
  loading.value = true
  try {
    const data = editing.value ? edit : props.fields
    await ingredientService.createIngredient({
      name: data.name,
      quantity: data.quantity,
      unit: data.unit,
      cost_per_unit: data.cost_per_unit,
    })
    confirmed.value = true
    emit('confirmed', { type: 'add_ingredient', name: data.name })
  } catch (e) {
    error.value = e.message || 'Không thể thêm nguyên liệu'
  } finally {
    loading.value = false
  }
}
</script>
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/ai/FieldRow.vue frontend/src/components/ai/IngredientActionCard.vue
git commit -m "feat(ai): add FieldRow and IngredientActionCard components"
```

---

## Task 6: RestockActionCard component

**Files:**
- Create: `frontend/src/components/ai/RestockActionCard.vue`

- [ ] **Step 1: Create the component**

```vue
<template>
  <div class="border-2 border-blue-400 rounded-2xl overflow-hidden bg-white shadow-sm w-full max-w-sm">
    <div class="bg-blue-50 px-4 py-3 flex items-center gap-2 border-b border-blue-200">
      <span class="text-base">🔄</span>
      <span class="font-semibold text-sm text-blue-900">Nhập kho nguyên liệu</span>
    </div>
    <div class="px-4 py-3 flex flex-col gap-2">
      <div v-if="!editing">
        <FieldRow label="Tên" :value="fields.ingredient_name" />
        <FieldRow label="Tồn kho hiện tại" :value="`${fields.current_stock} ${fields.unit}`" />
        <FieldRow label="Số lượng nhập" :value="`${fields.quantity} ${fields.unit}`" />
        <FieldRow label="Đơn giá" :value="formatPrice(fields.cost_per_unit)" />
        <FieldRow label="Phương thức" :value="fields.money_type === 'cash' ? '💵 Tiền mặt' : '🏦 Chuyển khoản'" />
        <FieldRow v-if="fields.reason" label="Lý do" :value="fields.reason" />
      </div>
      <div v-else class="flex flex-col gap-2">
        <label class="text-xs text-gray-500">Số lượng nhập</label>
        <input v-model.number="edit.quantity" type="number" min="0" class="border rounded-lg px-3 py-2 text-sm w-full" />
        <label class="text-xs text-gray-500">Đơn giá (VNĐ)</label>
        <input v-model.number="edit.cost_per_unit" type="number" min="0" class="border rounded-lg px-3 py-2 text-sm w-full" />
        <label class="text-xs text-gray-500">Phương thức</label>
        <div class="flex gap-2">
          <button @click="edit.money_type = 'cash'"
            :class="edit.money_type === 'cash' ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700'"
            class="flex-1 py-2 rounded-lg text-sm font-medium">💵 Tiền mặt</button>
          <button @click="edit.money_type = 'transfer'"
            :class="edit.money_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700'"
            class="flex-1 py-2 rounded-lg text-sm font-medium">🏦 CK</button>
        </div>
        <label class="text-xs text-gray-500">Lý do (tuỳ chọn)</label>
        <input v-model="edit.reason" class="border rounded-lg px-3 py-2 text-sm w-full" />
      </div>
      <p v-if="error" class="text-xs text-red-600 mt-1">{{ error }}</p>
    </div>
    <div v-if="!confirmed" class="flex border-t border-gray-100">
      <button @click="toggleEdit" class="flex-1 py-3 text-sm text-gray-600 font-medium border-r border-gray-100 active:bg-gray-50">
        {{ editing ? 'Xem lại' : '✏️ Sửa' }}
      </button>
      <button @click="confirm" :disabled="loading" class="flex-1 py-3 text-sm text-white font-semibold bg-blue-500 active:bg-blue-600 disabled:bg-gray-300">
        {{ loading ? 'Đang lưu...' : '✓ Xác nhận' }}
      </button>
    </div>
    <div v-else class="px-4 py-3 text-sm text-green-700 font-medium text-center bg-green-50">✅ Đã nhập kho thành công</div>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
import { fundIngredientService } from '../../services/fundIngredientService'
import FieldRow from './FieldRow.vue'

const props = defineProps({ fields: { type: Object, required: true } })
const emit = defineEmits(['confirmed'])

const editing = ref(false)
const loading = ref(false)
const confirmed = ref(false)
const error = ref('')
const edit = reactive({ ...props.fields })

function toggleEdit() {
  if (!editing.value) Object.assign(edit, props.fields)
  editing.value = !editing.value
}
function formatPrice(v) {
  return new Intl.NumberFormat('vi-VN').format(v || 0) + 'đ'
}
async function confirm() {
  error.value = ''
  loading.value = true
  try {
    const data = editing.value ? edit : props.fields
    await fundIngredientService.restockIngredientFromFund(data.ingredient_id, {
      quantity: data.quantity,
      cost_per_unit: data.cost_per_unit,
      money_type: data.money_type,
      reason: data.reason || '',
    })
    confirmed.value = true
    emit('confirmed', { type: 'restock_ingredient', name: data.ingredient_name })
  } catch (e) {
    error.value = e.message || 'Không thể nhập kho'
  } finally {
    loading.value = false
  }
}
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/ai/RestockActionCard.vue
git commit -m "feat(ai): add RestockActionCard component"
```

---

## Task 7: ExpenseActionCard component

**Files:**
- Create: `frontend/src/components/ai/ExpenseActionCard.vue`

- [ ] **Step 1: Create the component**

```vue
<template>
  <div class="border-2 border-purple-400 rounded-2xl overflow-hidden bg-white shadow-sm w-full max-w-sm">
    <div class="bg-purple-50 px-4 py-3 flex items-center gap-2 border-b border-purple-200">
      <span class="text-base">💸</span>
      <span class="font-semibold text-sm text-purple-900">Ghi nhận chi phí</span>
    </div>
    <div class="px-4 py-3 flex flex-col gap-2">
      <div v-if="!editing">
        <FieldRow label="Mô tả" :value="fields.description" />
        <FieldRow label="Số tiền" :value="formatPrice(fields.amount)" :highlight="true" />
        <FieldRow label="Phương thức" :value="fields.money_type === 'cash' ? '💵 Tiền mặt' : '🏦 Chuyển khoản'" />
        <FieldRow label="Danh mục" :value="categoryName" />
        <FieldRow label="Ngày" :value="fields.date || today" />
      </div>
      <div v-else class="flex flex-col gap-2">
        <label class="text-xs text-gray-500">Mô tả</label>
        <input v-model="edit.description" class="border rounded-lg px-3 py-2 text-sm w-full" />
        <label class="text-xs text-gray-500">Số tiền (VNĐ)</label>
        <input v-model.number="edit.amount" type="number" min="0" class="border rounded-lg px-3 py-2 text-sm w-full" />
        <label class="text-xs text-gray-500">Phương thức</label>
        <div class="flex gap-2">
          <button @click="edit.money_type = 'cash'"
            :class="edit.money_type === 'cash' ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700'"
            class="flex-1 py-2 rounded-lg text-sm font-medium">💵 Tiền mặt</button>
          <button @click="edit.money_type = 'transfer'"
            :class="edit.money_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700'"
            class="flex-1 py-2 rounded-lg text-sm font-medium">🏦 CK</button>
        </div>
        <label class="text-xs text-gray-500">Danh mục</label>
        <select v-model="edit.category_id" class="border rounded-lg px-3 py-2 text-sm w-full">
          <option value="">-- Chọn danh mục --</option>
          <option v-for="cat in categories" :key="cat._id" :value="cat._id">{{ cat.name }}</option>
        </select>
        <label class="text-xs text-gray-500">Ngày</label>
        <input v-model="edit.date" type="date" class="border rounded-lg px-3 py-2 text-sm w-full" />
      </div>
      <p v-if="error" class="text-xs text-red-600 mt-1">{{ error }}</p>
    </div>
    <div v-if="!confirmed" class="flex border-t border-gray-100">
      <button @click="toggleEdit" class="flex-1 py-3 text-sm text-gray-600 font-medium border-r border-gray-100 active:bg-gray-50">
        {{ editing ? 'Xem lại' : '✏️ Sửa' }}
      </button>
      <button @click="confirm" :disabled="loading" class="flex-1 py-3 text-sm text-white font-semibold bg-purple-500 active:bg-purple-600 disabled:bg-gray-300">
        {{ loading ? 'Đang lưu...' : '✓ Xác nhận' }}
      </button>
    </div>
    <div v-else class="px-4 py-3 text-sm text-green-700 font-medium text-center bg-green-50">✅ Đã ghi nhận chi phí thành công</div>
  </div>
</template>
<script setup>
import { ref, reactive, computed } from 'vue'
import { fundExpenseService } from '../../services/fundExpenseService'
import FieldRow from './FieldRow.vue'

const props = defineProps({
  fields: { type: Object, required: true },
  categories: { type: Array, default: () => [] },
})
const emit = defineEmits(['confirmed'])

const editing = ref(false)
const loading = ref(false)
const confirmed = ref(false)
const error = ref('')
const edit = reactive({ ...props.fields })
const today = new Date().toISOString().slice(0, 10)

const categoryName = computed(() => {
  const cat = props.categories.find(c => c._id === props.fields.category_id)
  return cat ? cat.name : '—'
})
function toggleEdit() {
  if (!editing.value) Object.assign(edit, props.fields)
  editing.value = !editing.value
}
function formatPrice(v) {
  return new Intl.NumberFormat('vi-VN').format(v || 0) + 'đ'
}
async function confirm() {
  error.value = ''
  loading.value = true
  try {
    const data = editing.value ? edit : props.fields
    await fundExpenseService.createExpenseFromFund({
      description: data.description,
      amount: data.amount,
      money_type: data.money_type,
      category_id: data.category_id,
      date: data.date || today,
    })
    confirmed.value = true
    emit('confirmed', { type: 'add_expense', description: data.description })
  } catch (e) {
    error.value = e.message || 'Không thể ghi nhận chi phí'
  } finally {
    loading.value = false
  }
}
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/ai/ExpenseActionCard.vue
git commit -m "feat(ai): add ExpenseActionCard component"
```

---

## Task 8: AICommandView — main chat page

**Files:**
- Create: `frontend/src/views/AICommandView.vue`

- [ ] **Step 1: Create the view**

```vue
<template>
  <div class="flex flex-col h-screen bg-gray-50">
    <div class="bg-white border-b border-gray-200 px-4 py-3 flex items-center gap-2 shrink-0">
      <span class="text-lg">🤖</span>
      <h1 class="font-semibold text-gray-900">AI Command</h1>
    </div>

    <div ref="threadRef" class="flex-1 overflow-y-auto px-4 py-4 flex flex-col gap-3">
      <div v-if="loading" class="text-center text-sm text-gray-400 py-8">Đang tải lịch sử...</div>

      <template v-for="(msg, i) in messages" :key="i">
        <div v-if="msg.role === 'user'" class="flex justify-end">
          <div class="bg-blue-500 text-white px-4 py-2 rounded-2xl rounded-br-sm text-sm max-w-xs">
            {{ msg.message }}
          </div>
        </div>
        <div v-else class="flex flex-col gap-2 items-start">
          <div class="bg-white border border-gray-200 px-4 py-2 rounded-2xl rounded-bl-sm text-sm max-w-xs text-gray-800">
            {{ msg.message }}
          </div>
          <IngredientActionCard
            v-if="msg.action_type === 'add_ingredient' && msg.action"
            :fields="msg.action.fields"
            @confirmed="onConfirmed"
          />
          <RestockActionCard
            v-else-if="msg.action_type === 'restock_ingredient' && msg.action"
            :fields="msg.action.fields"
            @confirmed="onConfirmed"
          />
          <ExpenseActionCard
            v-else-if="msg.action_type === 'add_expense' && msg.action"
            :fields="msg.action.fields"
            :categories="categories"
            @confirmed="onConfirmed"
          />
        </div>
      </template>

      <div v-if="thinking" class="flex items-center gap-1 text-gray-400 text-sm">
        <span class="animate-bounce">●</span><span class="animate-bounce" style="animation-delay:0.1s">●</span><span class="animate-bounce" style="animation-delay:0.2s">●</span>
      </div>
    </div>

    <div class="bg-white border-t border-gray-200 px-4 py-3 flex gap-2 items-center shrink-0">
      <input
        v-model="input"
        @keydown.enter.prevent="send"
        class="flex-1 bg-gray-50 border border-gray-200 rounded-full px-4 py-2 text-sm outline-none focus:border-blue-400"
        placeholder="Nhập lệnh cho AI..."
        :disabled="thinking"
      />
      <button
        @click="send"
        :disabled="!input.trim() || thinking"
        class="w-10 h-10 bg-blue-500 disabled:bg-gray-300 rounded-full text-white text-lg flex items-center justify-center shrink-0"
      >↑</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { aiService } from '../services/aiService'
import { expenseService } from '../services/expense'
import IngredientActionCard from '../components/ai/IngredientActionCard.vue'
import RestockActionCard from '../components/ai/RestockActionCard.vue'
import ExpenseActionCard from '../components/ai/ExpenseActionCard.vue'

const messages = ref([])
const input = ref('')
const thinking = ref(false)
const loading = ref(true)
const threadRef = ref(null)
const categories = ref([])

onMounted(async () => {
  try {
    const [logs, cats] = await Promise.all([
      aiService.getHistory(),
      expenseService.getCategories(),
    ])
    messages.value = logs.map(l => ({ ...l, action: null }))
    categories.value = cats || []
  } catch (e) {
    console.error('Failed to load AI history or categories', e)
  } finally {
    loading.value = false
    scrollToBottom()
  }
})

async function send() {
  const text = input.value.trim()
  if (!text || thinking.value) return
  input.value = ''
  messages.value.push({ role: 'user', message: text, action_type: null, action: null })
  thinking.value = true
  scrollToBottom()
  try {
    const history = messages.value.slice(-20).map(m => ({ role: m.role, message: m.message }))
    const resp = await aiService.parseCommand(text, history)
    messages.value.push({
      role: 'agent',
      message: resp.reply_text,
      action_type: resp.action?.type || null,
      action: resp.action || null,
    })
  } catch {
    messages.value.push({ role: 'agent', message: 'Có lỗi xảy ra. Vui lòng thử lại.', action_type: null, action: null })
  } finally {
    thinking.value = false
    scrollToBottom()
  }
}

function onConfirmed(payload) {
  const label = payload.type === 'add_expense' ? 'ghi nhận chi phí' : 'cập nhật nguyên liệu'
  messages.value.push({ role: 'agent', message: `✅ Đã ${label} thành công! Bạn cần làm gì thêm không?`, action_type: null, action: null })
  scrollToBottom()
}

function scrollToBottom() {
  nextTick(() => { if (threadRef.value) threadRef.value.scrollTop = threadRef.value.scrollHeight })
}
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/AICommandView.vue
git commit -m "feat(ai): add AICommandView chat page"
```

---

## Task 9: Router + Navigation

**Files:**
- Modify: `frontend/src/router/index.js`
- Modify: `frontend/src/components/Navigation.vue`

- [ ] **Step 1: Add import to router/index.js**

In the imports block at the top of the file, add:
```js
import AICommandView from '../views/AICommandView.vue'
```

- [ ] **Step 2: Add route in the routes array** (after the last manager route):

```js
{
  path: '/manager/ai-command',
  name: 'AICommand',
  component: AICommandView,
  meta: { requiresAuth: true, requiresManager: true }
},
```

- [ ] **Step 3: Add nav link in Navigation.vue**

Inside the `v-if="userRole === 'manager'"` grid div, after the last existing `router-link`, add:
```html
<router-link to="/manager/ai-command" @click="handleNavClick"
  class="flex flex-col items-center gap-1 p-3 rounded-xl bg-white shadow-sm hover:bg-blue-50 transition-colors">
  <span class="text-2xl">🤖</span>
  <span class="text-xs font-medium text-gray-700">AI Command</span>
</router-link>
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/router/index.js frontend/src/components/Navigation.vue
git commit -m "feat(ai): add AI Command route and navigation link"
```

---

## Task 10: Manual smoke test

- [ ] **Step 1: Start servers**

```bash
# Terminal 1 — backend
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS" && ./start-backend.sh

# Terminal 2 — frontend
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/frontend" && npm run dev
```

- [ ] **Step 2: Test add ingredient**
  1. Log in as manager → navigate to AI Command via nav
  2. Type: `Thêm nguyên liệu sữa tươi`
  3. Expected: amber action card appears below agent reply
  4. Click **Sửa**, fill in correct values, click **Xác nhận**
  5. Expected: "✅ Đã cập nhật nguyên liệu thành công!" in chat
  6. Verify ingredient appears in `/ingredients`

- [ ] **Step 3: Test expense**
  1. Type: `Chi phí điện tháng 4`
  2. Expected: purple expense card appears
  3. Click **Sửa**, set amount + category, click **Xác nhận**
  4. Expected: success message in chat
  5. Verify expense appears in `/expenses`

- [ ] **Step 4: Test history persistence**
  1. Reload the page
  2. Expected: previous conversation messages visible in thread

- [ ] **Step 5: Commit any fixes**

```bash
git add -p
git commit -m "fix(ai): smoke test fixes"
```
