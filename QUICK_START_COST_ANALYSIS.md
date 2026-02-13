# Quick Start: Test Chi Phí & Lợi Nhuận

## 🚀 Bắt Đầu Nhanh (3 Bước)

### 1️⃣ Seed Menu Data

```bash
go run backend/cmd/seed-menu-variants/main.go
```

Tạo 10 món:
- 6 món có nhiều size (M, L, XL)
- 4 món size đơn

### 2️⃣ Tính Chi Phí

```bash
./calculate-all-menu-costs.sh
```

Tự động tính chi phí cho tất cả món.

### 3️⃣ Xem Phân Tích

```bash
./view-menu-cost-analysis.sh
```

Hiển thị chi phí, lợi nhuận, và margin cho từng món/size.

## 📊 Kết Quả Mẫu

```
🔄 Cà phê sữa đá (Cà phê) - MULTI-SIZE

Variant    | Price        | Cost         | Profit       | Margin
-----------|--------------|--------------|--------------|------------
Size M     |     25000 đ |     13800 đ |     11200 đ |      44.8%
Size L     |     30000 đ |     20700 đ |      9300 đ |      31.0%
Size XL    |     35000 đ |     27600 đ |      7400 đ |      21.1%

⭐ Most profitable: Size M (44.8% margin)
```

## 🌐 Test Trên Frontend

```
http://localhost:5173/cost-analysis
```

## 📖 Hướng Dẫn Chi Tiết

Xem file: `HUONG_DAN_TEST_CHI_PHI_LOI_NHUAN.md`

## 🔗 API Endpoints

```bash
# Tính chi phí
POST /api/menu/:id/calculate-cost

# Xem chi tiết chi phí
GET /api/menu/:id/cost-breakdown

# Phân tích lợi nhuận
GET /api/menu/:id/profit-analysis
```

## ✅ Checklist

- [ ] Seed menu data
- [ ] Seed ingredients (nếu chưa có)
- [ ] Tính chi phí tất cả món
- [ ] Xem phân tích trên terminal
- [ ] Test trên frontend
- [ ] Test các API endpoints

Chúc bạn test thành công! 🎉
