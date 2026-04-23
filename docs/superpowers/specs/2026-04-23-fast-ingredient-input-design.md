# Fast Ingredient Input — Design Spec

**Date:** 2026-04-23
**Status:** Draft — pending implementation plan

## Goal

Give managers a dedicated screen to restock a commonly-used ingredient in a few taps: pick ingredient → tap `+`/`−` to set quantity → choose cash/bank → submit. Entry point is a new tile on the dashboard.

## User Flow

1. On the dashboard, in the **"Menu & Nguyên liệu"** section, user taps a new **"Nhập nhanh"** tile.
2. Screen opens showing ingredients restocked in the last 30 days, most-recently-restocked first.
3. User can type in the search box or toggle **"Hiện tất cả"** to see the full ingredient catalog (sorted by `created_at DESC`).
4. User taps an ingredient. The row expands to show:
   - Current stock
   - Counter: `[−]  [quantity field]  [+]` — quantity starts at 0; each `+`/`−` step equals the quantity of the ingredient's most recent restock (falls back to 1 if no history)
   - Cost per unit — pre-filled from the most recent restock's `cost_per_unit` (falls back to the ingredient's current `cost_per_unit`), editable
   - Total = quantity × cost (read-only)
   - Payment method: segmented control **Tiền mặt / Chuyển khoản** (default: Tiền mặt)
   - **[Lưu]** button
5. On submit: call the existing fund-restock endpoint, show success toast, clear the active selection, refresh the list. User stays on the screen ready for the next input.

## Non-goals

- Batch submission of multiple ingredients in one payment. (Explicitly single-ingredient per submission.)
- Gifted / free restocks (no-fund path). Fast input always deducts from the fund.
- New permission scopes — this screen uses the same manager role guard as `/ingredients`.

## Backend

### New endpoint

`GET /api/manager/ingredients/recent-restocks?days=30`

Response: array of ingredients that have at least one restock in the last `days` days, enriched with the most-recent restock's quantity and cost.

```json
[
  {
    "id": "65f…",
    "name": "Sữa tươi",
    "unit": "L",
    "quantity": 5.2,
    "cost_per_unit": 28000,
    "last_restock": {
      "quantity": 10,
      "cost_per_unit": 27500,
      "created_at": "2026-04-21T08:30:00Z"
    }
  }
]
```

- Sorted by `last_restock.created_at DESC`.
- Default `days` = 30; clamp to `[1, 365]`.
- Manager-only (same auth group as existing `/manager/ingredients` routes).

### Changes

- **Handler** (`backend/interfaces/http/ingredient_handler.go`): new `GetRecentRestockedIngredients` handler method wired under the manager group.
- **Service** (`backend/application/services/ingredient_service.go` or a new service file): new `GetRecentRestockedIngredients(ctx, days int)` orchestrating repo calls and returning the enriched DTO.
- **Repository** (`backend/infrastructure/mongodb/ingredient_restock_repository.go`): new `FindLatestByIngredientSince(ctx, since time.Time)` that aggregates `ingredient_restock_history` grouped by `ingredient_id`, keeping the latest record per group. The service then joins with the ingredients collection to build the response.

### Submit reuses existing endpoint

```
POST /api/manager/ingredients/:id/restock/from-fund
{
  "quantity": <number>,
  "cost_per_unit": <number>,
  "reason": "Nhập nhanh",
  "money_type": "cash" | "bank"
}
```

No backend changes needed for the submit path — this endpoint already handles fund deduction, expense creation, and restock-record insertion atomically.

### "Hiện tất cả" fallback

Uses the existing `GET /manager/ingredients` endpoint. Client sorts by `created_at DESC`.

## Frontend

### New files

- `frontend/src/views/FastIngredientInputView.vue`
- `frontend/src/views/__tests__/FastIngredientInputView.test.js`

### Modified files

- `frontend/src/services/ingredient.js` — add `getRecentRestocked(days = 30)` method hitting the new endpoint.
- `frontend/src/router/index.js` — add route `/ingredients/fast-input` with manager guard, pointing to `FastIngredientInputView`.
- `frontend/src/views/DashboardView.vue` — add a 5th tile inside the "Menu & Nguyên liệu" grid. Gradient `from-lime-500 to-green-500`, icon `➕`, label `Nhập nhanh`, sub-label `Nhập kho siêu nhanh`, routing to `/ingredients/fast-input`. Grid stays `grid-cols-2`; the 5th tile sits alone on the last row.

### View structure (`FastIngredientInputView.vue`)

**Local state:**

- `ingredients` — recent-restocked list (from new endpoint)
- `allIngredients` — full catalog (fetched lazily when `showAll` turns on)
- `showAll` — boolean toggle
- `search` — filter string (case-insensitive, matches name substring)
- `activeId` — currently expanded ingredient
- `quantity`, `costPerUnit`, `moneyType` — form state for active ingredient
- `submitting` — submit-in-progress flag

**Behavior:**

- On mount → fetch recent-restocked list.
- When `showAll` toggles on → fetch full catalog if not already cached; switch list source.
- Display list = `(showAll ? allIngredients : ingredients)` filtered by `search`, sorted as specified above.
- `selectIngredient(ing)`:
  - `activeId = ing.id`
  - `step = ing.last_restock?.quantity || 1`
  - `quantity = 0`
  - `costPerUnit = ing.last_restock?.cost_per_unit ?? ing.cost_per_unit`
  - `moneyType = 'cash'`
- `increment()` / `decrement()` → `quantity = max(0, quantity ± step)`.
- Quantity field is editable — user can type a custom number directly (minimum 0, decimals allowed).
- `submit()` → validates `quantity > 0`, then calls `fundIngredientService.restockIngredientFromFund(activeId, { quantity, cost_per_unit: costPerUnit, reason: 'Nhập nhanh', money_type: moneyType })`. On success: show toast, refresh recent list, clear active state. On error: toast/alert with message, keep form state.

### State ownership

Purely local to the view. No Pinia store changes. Existing `ingredientStore` is not required by this screen.

### Styling

Match existing ingredient management view patterns (Tailwind, large tap targets, rounded-xl cards) so the screen feels native to the app.

## Testing

### Frontend

Single component test at `frontend/src/views/__tests__/FastIngredientInputView.test.js` covering:

- Renders the recent-restocked list from the service
- Toggling "Hiện tất cả" switches the list source and sort order
- Search filters by name
- Selecting an ingredient initializes counter step + cost from `last_restock`
- `+` / `−` buttons step by `last_restock.quantity`
- Submit posts the correct payload to `fundIngredientService.restockIngredientFromFund`
- Submit success clears the active state and refreshes the list
- Submit blocked when quantity is 0

### Backend

Unit tests for the new service method covering:

- Only ingredients with restocks in the window are returned
- Sort order is `last_restock.created_at DESC`
- `days` parameter is clamped to `[1, 365]`
- Empty window returns an empty array

Repository aggregation tested against a seeded collection: multiple restocks per ingredient return only the latest per `ingredient_id`.

## Error handling

- Backend returns `500` on repository failure; handler wraps it in the standard error response.
- Frontend shows an alert/toast on fetch failure; screen continues to render prior cached data if present.
- Submit errors mirror the existing quick-stock-in modal: show the server error message, keep form state so the user can retry.

## Open questions

None at spec time.
