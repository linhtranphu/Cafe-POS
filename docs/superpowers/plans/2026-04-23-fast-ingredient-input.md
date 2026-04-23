# Fast Ingredient Input Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a "Nhập nhanh" dashboard tile and dedicated screen that lets managers restock any ingredient in a few taps, using the last restock's quantity as the step size and cost as the default price.

**Architecture:** New `GET /api/manager/ingredients/recent-restocks` endpoint aggregates the `ingredient_restock_history` collection to return the 30-day ingredient list enriched with last-restock metadata. The frontend `FastIngredientInputView` consumes this endpoint and calls the existing `POST .../restock/from-fund` endpoint for submission — no new submit path needed.

**Tech Stack:** Go (Gin, MongoDB driver), Vue 3 Composition API, Vitest + @vue/test-utils, Tailwind CSS.

**Spec:** `docs/superpowers/specs/2026-04-23-fast-ingredient-input-design.md`

---

## File Map

| File | Action | Responsibility |
|------|--------|----------------|
| `backend/infrastructure/mongodb/ingredient_restock_repository.go` | Modify | Add `FindLatestByIngredientSince` aggregation method |
| `backend/application/services/ingredient.go` | Modify | Add `RestockRepository` interface, `RecentRestockedIngredient` DTO, `GetRecentRestockedIngredients` service method |
| `backend/interfaces/http/ingredient_handler.go` | Modify | Add `GetRecentRestockedIngredients` handler method |
| `backend/main.go` | Modify | Wire `SetRestockRepository` and register the new route |
| `frontend/src/services/ingredient.js` | Modify | Add `getRecentRestocked(days)` API call |
| `frontend/src/router/index.js` | Modify | Add `/ingredients/fast-input` route |
| `frontend/src/views/DashboardView.vue` | Modify | Add "Nhập nhanh" tile in Menu & Nguyên liệu section |
| `frontend/src/views/FastIngredientInputView.vue` | Create | New fast-input screen |
| `frontend/src/views/__tests__/FastIngredientInputView.test.js` | Create | Component tests |

---

## Task 1: Repository — `FindLatestByIngredientSince` aggregation

**Files:**
- Modify: `backend/infrastructure/mongodb/ingredient_restock_repository.go`

This method aggregates `ingredient_restock_history`, groups by `ingredient_id`, and returns the most recent restock record per ingredient within the given time window.

- [ ] **Step 1: Add the method to `IngredientRestockRepository`**

Open `backend/infrastructure/mongodb/ingredient_restock_repository.go` and append after the `CountByIngredientID` method (after line 127):

```go
// FindLatestByIngredientSince returns the most recent restock record per ingredient
// for restocks created on or after `since`. Results are sorted by created_at DESC.
func (r *IngredientRestockRepository) FindLatestByIngredientSince(
	ctx context.Context,
	since time.Time,
) ([]*ingredient.IngredientRestockRecord, error) {
	pipeline := mongo.Pipeline{
		// Stage 1: filter by date window
		{{Key: "$match", Value: bson.M{
			"created_at": bson.M{"$gte": since},
		}}},
		// Stage 2: sort descending so $first picks the latest
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		// Stage 3: group by ingredient_id, keep latest record's fields
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$ingredient_id"},
			{Key: "ingredient_id", Value: bson.M{"$first": "$ingredient_id"}},
			{Key: "quantity", Value: bson.M{"$first": "$quantity"}},
			{Key: "cost_per_unit", Value: bson.M{"$first": "$cost_per_unit"}},
			{Key: "total_cost", Value: bson.M{"$first": "$total_cost"}},
			{Key: "performed_by", Value: bson.M{"$first": "$performed_by"}},
			{Key: "performed_by_name", Value: bson.M{"$first": "$performed_by_name"}},
			{Key: "reason", Value: bson.M{"$first": "$reason"}},
			{Key: "created_at", Value: bson.M{"$first": "$created_at"}},
		}}},
		// Stage 4: sort final results by latest restock desc
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode into bson.M because $group overwrites _id
	var rawResults []bson.M
	if err := cursor.All(ctx, &rawResults); err != nil {
		return nil, err
	}

	records := make([]*ingredient.IngredientRestockRecord, 0, len(rawResults))
	for _, raw := range rawResults {
		rec := &ingredient.IngredientRestockRecord{}
		if id, ok := raw["ingredient_id"].(primitive.ObjectID); ok {
			rec.IngredientID = id
		}
		if q, ok := raw["quantity"].(float64); ok {
			rec.Quantity = q
		}
		if c, ok := raw["cost_per_unit"].(float64); ok {
			rec.CostPerUnit = c
		}
		if t, ok := raw["total_cost"].(float64); ok {
			rec.TotalCost = t
		}
		if pb, ok := raw["performed_by"].(string); ok {
			rec.PerformedBy = pb
		}
		if pbn, ok := raw["performed_by_name"].(string); ok {
			rec.PerformedByName = pbn
		}
		if reason, ok := raw["reason"].(string); ok {
			rec.Reason = reason
		}
		if ca, ok := raw["created_at"].(primitive.DateTime); ok {
			rec.CreatedAt = ca.Time()
		}
		records = append(records, rec)
	}

	return records, nil
}
```

- [ ] **Step 2: Verify the file compiles**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend"
go build ./infrastructure/mongodb/...
```

Expected: no output (clean build).

- [ ] **Step 3: Commit**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS"
git add backend/infrastructure/mongodb/ingredient_restock_repository.go
git commit -m "feat(ingredient): add FindLatestByIngredientSince aggregation to restock repo"
```

---

## Task 2: Service — DTO + `GetRecentRestockedIngredients`

**Files:**
- Modify: `backend/application/services/ingredient.go`

- [ ] **Step 1: Add `RestockRepository` interface**

In `backend/application/services/ingredient.go`, after the `StockHistoryRepository` interface (around line 26), add:

```go
type RestockRepository interface {
	FindLatestByIngredientSince(ctx context.Context, since time.Time) ([]*ingredient.IngredientRestockRecord, error)
}
```

- [ ] **Step 2: Add DTO types**

After the `StockReportItem` struct (around line 42), add:

```go
// LastRestockInfo holds the most recent restock metadata for a single ingredient.
type LastRestockInfo struct {
	Quantity    float64   `json:"quantity"`
	CostPerUnit float64   `json:"cost_per_unit"`
	CreatedAt   time.Time `json:"created_at"`
}

// RecentRestockedIngredient is the response item for GET /ingredients/recent-restocks.
type RecentRestockedIngredient struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Unit        string           `json:"unit"`
	Quantity    float64          `json:"quantity"`
	CostPerUnit float64          `json:"cost_per_unit"`
	LastRestock *LastRestockInfo `json:"last_restock,omitempty"`
}
```

- [ ] **Step 3: Add `restockRepo` field and `SetRestockRepository` setter**

Update the `IngredientService` struct to add the field (insert after `stockHistoryRepo`):

```go
type IngredientService struct {
	ingredientRepo        IngredientRepository
	stockHistoryRepo      StockHistoryRepository
	restockRepo           RestockRepository
	autoExpenseService    *AutoExpenseService
	costCalculatorService *CostCalculatorService
}
```

Add the setter method after `SetCostCalculatorService`:

```go
// SetRestockRepository injects the restock repository (called after initialization).
func (s *IngredientService) SetRestockRepository(repo RestockRepository) {
	s.restockRepo = repo
}
```

- [ ] **Step 4: Add `GetRecentRestockedIngredients` method**

Append to the end of `backend/application/services/ingredient.go`:

```go
// GetRecentRestockedIngredients returns ingredients restocked within the past `days` days,
// sorted by most-recent-restock descending, enriched with last-restock metadata.
// days is clamped to [1, 365].
func (s *IngredientService) GetRecentRestockedIngredients(ctx context.Context, days int) ([]*RecentRestockedIngredient, error) {
	if days < 1 {
		days = 1
	}
	if days > 365 {
		days = 365
	}

	since := time.Now().AddDate(0, 0, -days)

	restocks, err := s.restockRepo.FindLatestByIngredientSince(ctx, since)
	if err != nil {
		return nil, err
	}

	allIngredients, err := s.ingredientRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	ingMap := make(map[primitive.ObjectID]*ingredient.Ingredient, len(allIngredients))
	for _, ing := range allIngredients {
		ingMap[ing.ID] = ing
	}

	results := make([]*RecentRestockedIngredient, 0, len(restocks))
	for _, r := range restocks {
		ing, ok := ingMap[r.IngredientID]
		if !ok {
			continue // ingredient deleted — skip
		}
		results = append(results, &RecentRestockedIngredient{
			ID:          ing.ID.Hex(),
			Name:        ing.Name,
			Unit:        string(ing.Unit),
			Quantity:    ing.Quantity,
			CostPerUnit: ing.CostPerUnit,
			LastRestock: &LastRestockInfo{
				Quantity:    r.Quantity,
				CostPerUnit: r.CostPerUnit,
				CreatedAt:   r.CreatedAt,
			},
		})
	}

	return results, nil
}
```

- [ ] **Step 5: Verify package compiles**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend"
go build ./application/services/...
```

Expected: no output.

- [ ] **Step 6: Commit**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS"
git add backend/application/services/ingredient.go
git commit -m "feat(ingredient): add RecentRestockedIngredient DTO and GetRecentRestockedIngredients service method"
```

---

## Task 3: Handler + route registration

**Files:**
- Modify: `backend/interfaces/http/ingredient_handler.go`
- Modify: `backend/main.go`

- [ ] **Step 1: Add `GetRecentRestockedIngredients` handler**

In `backend/interfaces/http/ingredient_handler.go`, append before the final closing brace of the file:

```go
// GetRecentRestockedIngredients handles GET /api/manager/ingredients/recent-restocks?days=30
func (h *IngredientHandler) GetRecentRestockedIngredients(c *gin.Context) {
	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		var parsed int
		if _, err := fmt.Sscanf(daysStr, "%d", &parsed); err == nil {
			days = parsed
		}
	}

	items, err := h.ingredientService.GetRecentRestockedIngredients(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}
```

- [ ] **Step 2: Wire `SetRestockRepository` in `main.go`**

In `backend/main.go`, find the block around line 343:
```go
ingredientService := services.NewIngredientService(ingredientRepo, stockHistoryRepo)
```

Immediately after that line, add:
```go
ingredientService.SetRestockRepository(ingredientRestockRepo)
```

(`ingredientRestockRepo` is already created on line 342.)

- [ ] **Step 3: Register the route in `main.go`**

In `backend/main.go`, find the manager ingredient routes block (around line 694):
```go
manager.GET("/ingredients/summary", ingredientHandler.GetStockSummary)
```

Add the new route immediately after it:
```go
manager.GET("/ingredients/recent-restocks", ingredientHandler.GetRecentRestockedIngredients)
```

- [ ] **Step 4: Verify full backend compiles**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/backend"
go build ./...
```

Expected: no output.

- [ ] **Step 5: Commit**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS"
git add backend/interfaces/http/ingredient_handler.go backend/main.go
git commit -m "feat(ingredient): add GetRecentRestockedIngredients handler and route"
```

---

## Task 4: Frontend service method

**Files:**
- Modify: `frontend/src/services/ingredient.js`

- [ ] **Step 1: Add `getRecentRestocked` to `ingredientService`**

In `frontend/src/services/ingredient.js`, add the following method inside the `ingredientService` object after `getStockSummary` (before the closing `}`):

```js
async getRecentRestocked(days = 30) {
  const response = await api.get('/manager/ingredients/recent-restocks', { params: { days } })
  return response.data
},
```

- [ ] **Step 2: Commit**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS"
git add frontend/src/services/ingredient.js
git commit -m "feat(ingredient): add getRecentRestocked service method"
```

---

## Task 5: Router entry

**Files:**
- Modify: `frontend/src/router/index.js`

- [ ] **Step 1: Add import**

In `frontend/src/router/index.js`, after line 6:
```js
import IngredientManagementView from '../views/IngredientManagementView.vue'
```

Add:
```js
import FastIngredientInputView from '../views/FastIngredientInputView.vue'
```

- [ ] **Step 2: Add route**

After the `/ingredients` route block (lines 52–56):
```js
{
  path: '/ingredients',
  name: 'Ingredients',
  component: IngredientManagementView,
  meta: { requiresAuth: true, requiresManager: true }
},
```

Add:
```js
{
  path: '/ingredients/fast-input',
  name: 'FastIngredientInput',
  component: FastIngredientInputView,
  meta: { requiresAuth: true, requiresManager: true }
},
```

**Note:** Do not commit the router alone — the import will be broken until `FastIngredientInputView.vue` exists. Commit router + view together in Task 6 Step 2.

---

## Task 6: `FastIngredientInputView.vue`

**Files:**
- Create: `frontend/src/views/FastIngredientInputView.vue`

- [ ] **Step 1: Create the file**

Create `frontend/src/views/FastIngredientInputView.vue`:

```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white border-b border-gray-200 px-4 py-3 flex items-center gap-3 sticky top-0 z-10">
      <button @click="$router.back()" class="p-2 rounded-lg hover:bg-gray-100 active:bg-gray-200 transition-colors">
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <h1 class="text-lg font-bold text-gray-800">Nhập nhanh nguyên liệu</h1>
    </div>

    <div class="p-4 space-y-3">
      <!-- Search + toggle -->
      <div class="flex gap-2">
        <input
          v-model="search"
          type="text"
          placeholder="Tìm nguyên liệu..."
          class="flex-1 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-lime-500 focus:border-transparent text-base"
        />
        <button
          @click="toggleShowAll"
          :class="[
            'px-3 py-2 rounded-xl text-sm font-medium border transition-colors whitespace-nowrap',
            showAll
              ? 'bg-lime-500 text-white border-lime-500'
              : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
          ]"
        >
          Hiện tất cả
        </button>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="text-center py-8 text-gray-400">Đang tải...</div>

      <!-- Error -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-4 text-red-700 text-sm">
        {{ error }}
      </div>

      <!-- Empty -->
      <div v-else-if="displayList.length === 0" class="text-center py-8 text-gray-400 text-sm">
        {{ search ? 'Không tìm thấy nguyên liệu phù hợp' : 'Chưa có nguyên liệu nào được nhập trong 30 ngày qua' }}
      </div>

      <!-- Ingredient list -->
      <div v-else class="space-y-2">
        <div
          v-for="ing in displayList"
          :key="ing.id"
          class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm"
        >
          <!-- Collapsed row -->
          <button
            @click="selectIngredient(ing)"
            class="w-full px-4 py-3 flex items-center justify-between text-left active:bg-gray-50 transition-colors"
          >
            <div>
              <div class="font-semibold text-gray-800">{{ ing.name }}</div>
              <div class="text-xs text-gray-500 mt-0.5">
                Tồn: {{ ing.quantity }} {{ ing.unit }}
                <span v-if="ing.last_restock" class="ml-2 text-lime-600">
                  • Lần cuối: {{ ing.last_restock.quantity }} {{ ing.unit }}
                </span>
              </div>
            </div>
            <svg
              :class="['w-5 h-5 text-gray-400 transition-transform', activeId === ing.id ? 'rotate-180' : '']"
              fill="none" stroke="currentColor" viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>

          <!-- Expanded form -->
          <div v-if="activeId === ing.id" class="border-t border-gray-100 px-4 pb-4 pt-3 space-y-4">
            <!-- Quantity counter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Số lượng ({{ ing.unit }})</label>
              <div class="flex items-center gap-2">
                <button
                  @click="decrement"
                  class="w-12 h-12 rounded-xl bg-gray-100 text-gray-700 text-xl font-bold flex items-center justify-center active:bg-gray-200 transition-colors"
                >−</button>
                <input
                  v-model.number="quantity"
                  type="number"
                  min="0"
                  step="any"
                  class="flex-1 h-12 text-center text-lg font-bold border-2 border-gray-300 rounded-xl focus:ring-2 focus:ring-lime-500 focus:border-lime-500"
                />
                <button
                  @click="increment"
                  class="w-12 h-12 rounded-xl bg-lime-500 text-white text-xl font-bold flex items-center justify-center active:bg-lime-600 transition-colors"
                >+</button>
              </div>
            </div>

            <!-- Cost per unit -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Đơn giá (₫ / {{ ing.unit }})</label>
              <input
                v-model.number="costPerUnit"
                type="number"
                min="0"
                class="w-full h-12 px-4 text-base border-2 border-gray-300 rounded-xl focus:ring-2 focus:ring-lime-500 focus:border-lime-500"
              />
            </div>

            <!-- Total (read-only) -->
            <div class="bg-lime-50 rounded-xl px-4 py-3 flex justify-between items-center">
              <span class="text-sm text-gray-600">Tổng tiền</span>
              <span class="text-lg font-bold text-lime-700">{{ formatCurrency(quantity * costPerUnit) }}</span>
            </div>

            <!-- Payment method -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Hình thức thanh toán</label>
              <div class="flex gap-2">
                <button
                  v-for="opt in paymentOptions"
                  :key="opt.value"
                  @click="moneyType = opt.value"
                  :class="[
                    'flex-1 py-3 rounded-xl font-medium text-sm border-2 transition-colors',
                    moneyType === opt.value
                      ? 'bg-lime-500 border-lime-500 text-white'
                      : 'bg-white border-gray-300 text-gray-700'
                  ]"
                >
                  {{ opt.label }}
                </button>
              </div>
            </div>

            <!-- Submit -->
            <button
              @click="submit(ing)"
              :disabled="submitting || quantity <= 0"
              class="w-full py-4 bg-lime-500 text-white rounded-xl font-bold text-base active:bg-lime-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {{ submitting ? 'Đang lưu...' : 'Lưu' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Toast -->
    <div
      v-if="toast.show"
      :class="[
        'fixed bottom-6 left-4 right-4 p-4 rounded-xl text-white font-medium text-center shadow-lg transition-all z-50',
        toast.success ? 'bg-green-500' : 'bg-red-500'
      ]"
    >
      {{ toast.message }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ingredientService } from '../services/ingredient'
import { fundIngredientService } from '../services/fundIngredientService'

// --- State ---
const ingredients = ref([])
const allIngredients = ref([])
const allFetched = ref(false)
const showAll = ref(false)
const search = ref('')
const loading = ref(false)
const error = ref(null)

const activeId = ref(null)
const step = ref(1)
const quantity = ref(0)
const costPerUnit = ref(0)
const moneyType = ref('cash')
const submitting = ref(false)

const toast = ref({ show: false, success: true, message: '' })

const paymentOptions = [
  { value: 'cash', label: 'Tiền mặt' },
  { value: 'bank', label: 'Chuyển khoản' }
]

// --- Derived list ---
const displayList = computed(() => {
  const source = showAll.value ? allIngredients.value : ingredients.value
  const q = search.value.toLowerCase().trim()
  return q ? source.filter(i => i.name.toLowerCase().includes(q)) : source
})

// --- Fetch ---
async function fetchRecent() {
  loading.value = true
  error.value = null
  try {
    ingredients.value = await ingredientService.getRecentRestocked(30)
  } catch {
    error.value = 'Không thể tải danh sách. Vui lòng thử lại.'
  } finally {
    loading.value = false
  }
}

async function fetchAll() {
  if (allFetched.value) return
  try {
    const data = await ingredientService.getIngredients()
    allIngredients.value = [...data].sort((a, b) => {
      const ta = a.created_at ? new Date(a.created_at).getTime() : 0
      const tb = b.created_at ? new Date(b.created_at).getTime() : 0
      return tb - ta
    })
    allFetched.value = true
  } catch {
    // Non-fatal: keep allIngredients empty
  }
}

async function toggleShowAll() {
  showAll.value = !showAll.value
  if (showAll.value) await fetchAll()
}

onMounted(fetchRecent)

// --- Ingredient selection ---
function selectIngredient(ing) {
  if (activeId.value === ing.id) {
    activeId.value = null
    return
  }
  activeId.value = ing.id
  step.value = ing.last_restock?.quantity ?? 1
  quantity.value = 0
  costPerUnit.value = ing.last_restock?.cost_per_unit ?? ing.cost_per_unit ?? 0
  moneyType.value = 'cash'
}

// --- Counter ---
function increment() {
  quantity.value = Math.max(0, +(quantity.value + step.value).toFixed(6))
}

function decrement() {
  quantity.value = Math.max(0, +(quantity.value - step.value).toFixed(6))
}

// --- Submit ---
async function submit(ing) {
  if (quantity.value <= 0) {
    showToast('Vui lòng nhập số lượng hợp lệ', false)
    return
  }
  submitting.value = true
  try {
    await fundIngredientService.restockIngredientFromFund(ing.id, {
      quantity: quantity.value,
      cost_per_unit: costPerUnit.value,
      reason: 'Nhập nhanh',
      money_type: moneyType.value
    })
    showToast(`Đã nhập ${quantity.value} ${ing.unit} ${ing.name}`, true)
    activeId.value = null
    await fetchRecent()
  } catch (e) {
    showToast(e?.response?.data?.error || e.message || 'Có lỗi xảy ra', false)
  } finally {
    submitting.value = false
  }
}

// --- Helpers ---
function formatCurrency(value) {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value || 0)
}

function showToast(message, success) {
  toast.value = { show: true, success, message }
  setTimeout(() => { toast.value.show = false }, 2500)
}
</script>
```

- [ ] **Step 2: Commit view + router together**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS"
git add frontend/src/views/FastIngredientInputView.vue frontend/src/router/index.js
git commit -m "feat(ingredient): add FastIngredientInputView and /ingredients/fast-input route"
```

---

## Task 7: Dashboard tile

**Files:**
- Modify: `frontend/src/views/DashboardView.vue`

- [ ] **Step 1: Add the "Nhập nhanh" tile**

In `frontend/src/views/DashboardView.vue`, find the "Tồn kho" button (around line 114–119) inside the "Menu & Nguyên liệu" grid. After its closing `</button>` tag and before the closing `</div>` of the grid, add:

```html
<button @click="$router.push('/ingredients/fast-input')"
  class="bg-gradient-to-br from-lime-500 to-green-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
  <div class="text-4xl mb-2">➕</div>
  <div class="font-bold text-base">Nhập nhanh</div>
  <div class="text-xs opacity-80 mt-1">Nhập kho siêu nhanh</div>
</button>
```

- [ ] **Step 2: Start the dev server and verify visually**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/frontend"
npm run dev
```

Open the dashboard. Confirm:
1. "Nhập nhanh" tile appears in the "Menu & Nguyên liệu" section.
2. Tapping it navigates to `/ingredients/fast-input`.
3. The fast-input screen renders the ingredient list (or empty state).
4. Selecting an ingredient expands the form with counter, cost, payment, and Lưu button.
5. Tapping "+" increments by `last_restock.quantity`.
6. Tapping "Lưu" submits and shows the success toast.

- [ ] **Step 3: Commit**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS"
git add frontend/src/views/DashboardView.vue
git commit -m "feat(dashboard): add Nhập nhanh tile to Menu & Nguyên liệu section"
```

---

## Task 8: Frontend component tests

**Files:**
- Create: `frontend/src/views/__tests__/FastIngredientInputView.test.js`

- [ ] **Step 1: Write the tests**

Create `frontend/src/views/__tests__/FastIngredientInputView.test.js`:

```js
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import FastIngredientInputView from '../FastIngredientInputView.vue'

vi.mock('../../services/ingredient', () => ({
  ingredientService: {
    getRecentRestocked: vi.fn(),
    getIngredients: vi.fn()
  }
}))

vi.mock('../../services/fundIngredientService', () => ({
  fundIngredientService: {
    restockIngredientFromFund: vi.fn()
  }
}))

import { ingredientService } from '../../services/ingredient'
import { fundIngredientService } from '../../services/fundIngredientService'

const mockIngredient = {
  id: 'abc123',
  name: 'Sữa tươi',
  unit: 'L',
  quantity: 5.0,
  cost_per_unit: 28000,
  created_at: '2026-04-01T00:00:00Z',
  last_restock: {
    quantity: 10,
    cost_per_unit: 27500,
    created_at: '2026-04-21T08:00:00Z'
  }
}

const mountView = () => mount(FastIngredientInputView, {
  global: {
    stubs: { RouterLink: true },
    mocks: { $router: { back: vi.fn(), push: vi.fn() } }
  }
})

describe('FastIngredientInputView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    ingredientService.getRecentRestocked.mockResolvedValue([mockIngredient])
    ingredientService.getIngredients.mockResolvedValue([mockIngredient])
    fundIngredientService.restockIngredientFromFund.mockResolvedValue({ success: true })
  })

  it('renders the recent-restocked list', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    expect(ingredientService.getRecentRestocked).toHaveBeenCalledWith(30)
    expect(wrapper.text()).toContain('Sữa tươi')
  })

  it('fetches all ingredients when Hiện tất cả is toggled', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const toggleBtn = wrapper.findAll('button').find(b => b.text().includes('Hiện tất cả'))
    await toggleBtn.trigger('click')
    await nextTick()
    await nextTick()
    expect(ingredientService.getIngredients).toHaveBeenCalled()
  })

  it('filters display list by search query', async () => {
    ingredientService.getRecentRestocked.mockResolvedValue([
      mockIngredient,
      { ...mockIngredient, id: 'xyz', name: 'Đường', last_restock: null }
    ])
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    await wrapper.find('input[placeholder*="Tìm"]').setValue('Sữa')
    await nextTick()
    expect(wrapper.text()).toContain('Sữa tươi')
    expect(wrapper.text()).not.toContain('Đường')
  })

  it('initializes quantity=0 and cost from last_restock when ingredient selected', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    const inputs = wrapper.findAll('input[type="number"]')
    expect(Number(inputs[0].element.value)).toBe(0)          // quantity
    expect(Number(inputs[1].element.value)).toBe(27500)      // costPerUnit from last_restock
  })

  it('increments quantity by last_restock.quantity per tap', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    const plusBtn = wrapper.findAll('button').find(b => b.text() === '+')
    await plusBtn.trigger('click')
    await nextTick()
    const qtyInput = wrapper.find('input[type="number"]')
    expect(Number(qtyInput.element.value)).toBe(10)
  })

  it('decrement does not go below 0', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    const minusBtn = wrapper.findAll('button').find(b => b.text() === '−')
    await minusBtn.trigger('click')
    await nextTick()
    expect(Number(wrapper.find('input[type="number"]').element.value)).toBe(0)
  })

  it('submit calls restockIngredientFromFund with correct payload', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    await wrapper.findAll('button').find(b => b.text() === '+').trigger('click')
    await nextTick()
    await wrapper.findAll('button').find(b => b.text() === 'Lưu').trigger('click')
    await nextTick()
    expect(fundIngredientService.restockIngredientFromFund).toHaveBeenCalledWith(
      'abc123',
      expect.objectContaining({
        quantity: 10,
        cost_per_unit: 27500,
        reason: 'Nhập nhanh',
        money_type: 'cash'
      })
    )
  })

  it('submit success clears active selection and refreshes list', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    await wrapper.findAll('button').find(b => b.text() === '+').trigger('click')
    await nextTick()
    await wrapper.findAll('button').find(b => b.text() === 'Lưu').trigger('click')
    await nextTick()
    await nextTick()
    expect(ingredientService.getRecentRestocked).toHaveBeenCalledTimes(2)
    // Expanded form gone
    expect(wrapper.findAll('input[type="number"]').length).toBe(0)
  })

  it('Lưu button is disabled when quantity is 0', async () => {
    const wrapper = mountView()
    await nextTick()
    await nextTick()
    const rowBtn = wrapper.findAll('button').find(b => b.text().includes('Sữa tươi'))
    await rowBtn.trigger('click')
    await nextTick()
    const saveBtn = wrapper.findAll('button').find(b => b.text() === 'Lưu')
    expect(saveBtn.attributes('disabled')).toBeDefined()
    await saveBtn.trigger('click')
    await nextTick()
    expect(fundIngredientService.restockIngredientFromFund).not.toHaveBeenCalled()
  })
})
```

- [ ] **Step 2: Run the tests**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS/frontend"
npx vitest run src/views/__tests__/FastIngredientInputView.test.js
```

Expected: all 8 tests pass. If any fail, fix the view or test before continuing.

- [ ] **Step 3: Commit**

```bash
cd "/Volumes/Linh-DAT/Misc Projects/Cafe POS"
git add frontend/src/views/__tests__/FastIngredientInputView.test.js
git commit -m "test(ingredient): add FastIngredientInputView component tests"
```

---

## Self-Review

**Spec coverage:**
- [x] `GET /api/manager/ingredients/recent-restocks?days=30` → Tasks 1–3
- [x] `days` clamped to [1, 365] → Task 2 Step 4
- [x] `getRecentRestocked` frontend service method → Task 4
- [x] `/ingredients/fast-input` route with manager guard → Task 5
- [x] "Nhập nhanh" dashboard tile → Task 7
- [x] Search box + "Hiện tất cả" toggle (sorted by `created_at DESC`) → Task 6
- [x] Counter step = `last_restock.quantity` (fallback 1) → Task 6 `selectIngredient`
- [x] Cost default = `last_restock.cost_per_unit` (fallback `ing.cost_per_unit`) → Task 6
- [x] Payment: Tiền mặt / Chuyển khoản → Task 6
- [x] Submit calls `fundIngredientService.restockIngredientFromFund` with `reason: 'Nhập nhanh'` → Task 6
- [x] Success: toast + clear active + refresh list → Task 6
- [x] Error: toast with server message, keep form → Task 6
- [x] Lưu disabled when quantity = 0 → Task 6
- [x] Tests for all behaviors above → Task 8

**Placeholders:** None.

**Type consistency:** `ing.id` (string hex) used in `selectIngredient`, `submit`, and tests. `ing.last_restock?.quantity` and `ing.last_restock?.cost_per_unit` used consistently in `selectIngredient` and test assertions. `RecentRestockedIngredient.ID` serialized as `ing.id` via `.Hex()` in the Go service.
