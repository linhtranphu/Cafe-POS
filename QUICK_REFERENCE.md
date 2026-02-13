# 🚀 Quick Reference - Cost Analysis Testing

## Backend Info
- **Port:** 3000
- **URL:** http://localhost:3000
- **Run:** `cd backend && go run main.go`

## Quick Commands

### 1. Start Backend
```bash
cd backend
go run main.go
```

### 2. Setup Everything
```bash
./setup-and-calculate-costs.sh
```

### 3. Check Backend
```bash
curl http://localhost:3000/api/menu
```

### 4. Get Menu Item ID
```bash
curl -s http://localhost:3000/api/menu | grep -o '"id":"[^"]*"' | head -1
```

### 5. View Cost Breakdown
```bash
MENU_ID="<paste_id_here>"
curl http://localhost:3000/api/menu/$MENU_ID/cost-breakdown | python3 -m json.tool
```

### 6. View Profit Analysis
```bash
curl http://localhost:3000/api/menu/$MENU_ID/profit-analysis | python3 -m json.tool
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/menu` | List all menu items |
| GET | `/api/menu/:id` | Get menu item detail |
| POST | `/api/menu/:id/calculate-cost` | Calculate cost |
| GET | `/api/menu/:id/cost-breakdown` | Cost breakdown |
| GET | `/api/menu/:id/profit-analysis` | Profit analysis |

## Troubleshooting

### Backend not responding?
```bash
# Check if running
lsof -i :3000

# Check MongoDB
ps aux | grep mongod
```

### Need to reseed?
```bash
./seed-menu-variants-auto.sh
```

### Recalculate costs?
```bash
./calculate-costs-simple.sh
```

## Expected Results

### Cà phê sữa đá
- Size M: Cost 13,800đ, Margin 44.8%
- Size L: Cost 20,700đ, Margin 31.0%
- Size XL: Cost 27,600đ, Margin 21.1%

## Frontend
```
http://localhost:5173/cost-analysis
```

## Files
- `START_HERE.md` - Detailed guide
- `COST_ANALYSIS_SETUP_COMPLETE.md` - Complete info
- `HUONG_DAN_CHAY_DAY_DU.md` - Vietnamese guide
