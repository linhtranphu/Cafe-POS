# Quản lý Quỹ Tiền (Fund Management) — Requirements

## Tổng quan

Hệ thống quản lý 4 quỹ tiền tách biệt theo chức năng, cho phép manager theo dõi dòng tiền chi tiết, xem báo cáo theo ngày, và chuyển tiền giữa các quỹ.

---

## Mô hình 4 quỹ

| Quỹ | ID | Mục đích | Nguồn vào | Nguồn ra |
|-----|----|----------|-----------|----------|
| 🗄️ Ngăn kéo tiền | `cash_drawer` | Tiền mặt tại quầy thu ngân | Bàn giao từ waiter (cash_handover) | Tiền đầu ca (starting_float), chuyển về quỹ vận hành |
| ⚙️ Quỹ vận hành | `operating` | Chi phí vận hành hàng ngày | Bàn giao ca (fund_handover), nạp thủ công | Chi phí vận hành, chuyển quỹ |
| 📦 Quỹ nguyên liệu | `inventory` | Mua nguyên liệu / hàng tồn kho | Chuyển từ quỹ vận hành | Mua nguyên liệu |
| 💎 Quỹ lợi nhuận | `profit` | Lợi nhuận tích lũy | Chuyển từ quỹ vận hành | Rút thủ công |

---

## User Stories

### US1: Xem tổng quan 4 quỹ
**Là** manager
**Tôi muốn** xem số dư của từng quỹ
**Để** biết tiền đang nằm ở đâu

**Acceptance Criteria:**
- Hiển thị 4 thẻ quỹ theo lưới 2×2
- Mỗi thẻ hiển thị: tên quỹ, icon, tổng số dư, phân tách cash / transfer
- Hiển thị thêm tổng tất cả quỹ ở đầu trang
- Load dữ liệu realtime khi refresh

---

### US2: Báo cáo ngày
**Là** manager
**Tôi muốn** xem báo cáo thu chi trong ngày
**Để** đối chiếu cuối ngày

**Acceptance Criteria:**
- Hiển thị số dư đầu ngày (opening balance) của từng quỹ
- Tổng tiền vào hôm nay (inflow) phân theo nguồn:
  - Bàn giao ca (handover)
  - Nạp thủ công (deposit)
- Tổng tiền ra hôm nay (outflow) phân theo nguồn:
  - Chi phí vận hành (expense)
  - Mua nguyên liệu (ingredient)
  - Rút thủ công (withdrawal)
- Số dư cuối ngày (= đầu ngày + vào - ra)

---

### US3: Lịch sử giao dịch với filter
**Là** manager
**Tôi muốn** xem và lọc lịch sử giao dịch
**Để** tra cứu và kiểm soát dòng tiền

**Acceptance Criteria:**
- Filter theo loại giao dịch: tất cả / nạp / rút / bàn giao / chuyển quỹ
- Filter theo quỹ: tất cả / từng quỹ cụ thể
- Filter theo nguồn: tất cả / expense / nguyên liệu / bàn giao / thủ công
- Filter theo loại tiền: tất cả / tiền mặt / chuyển khoản
- Filter theo ngày: hôm nay / hôm qua / tuần này / tháng này / tùy chỉnh
- Phân trang (load more)

---

### US4: Chi tiết nguồn giao dịch
**Là** manager
**Tôi muốn** xem giao dịch quỹ liên kết đến record gốc
**Để** hiểu tại sao tiền ra/vào

**Acceptance Criteria:**
- Mỗi giao dịch có badge hiển thị nguồn (expense / ingredient / handover / manual / fund_transfer)
- Nhấn vào giao dịch → xem modal chi tiết
- Modal chi tiết hiển thị:
  - Thông tin giao dịch (type, amount, reason, người thực hiện, thời gian)
  - Số dư trước / sau
  - Link đến record gốc: expense → xem expense; ingredient → xem restock; handover → xem ca
  - Với giao dịch chuyển quỹ: hiển thị quỹ nguồn → quỹ đích

---

### US5: Bàn giao ca → Quỹ
**Là** manager
**Tôi muốn** xem giao dịch bàn giao ca phản ánh đúng trong quỹ
**Để** đối chiếu với thực tế

**Acceptance Criteria:**
- Khi cashier kết thúc ca và bàn giao tiền, tạo FundTransaction với type=`fund_handover`, fund_type=`cash_drawer`
- Hiển thị trong lịch sử quỹ với badge "Bàn giao ca"
- Link tới cashier shift record tương ứng

---

### US6: Nạp tiền / Rút tiền thủ công
**Là** manager
**Tôi muốn** nạp hoặc rút tiền thủ công từ bất kỳ quỹ nào
**Để** điều chỉnh số dư quỹ

**Acceptance Criteria:**
- Chọn quỹ muốn nạp/rút
- Chọn loại tiền: tiền mặt / chuyển khoản
- Nhập số tiền (rút: validate ≤ số dư loại tiền đó trong quỹ)
- Nhập lý do (tối thiểu 10 ký tự)
- Hiển thị số dư quỹ hiện tại trước khi xác nhận

---

### US7: Chuyển tiền giữa quỹ
**Là** manager
**Tôi muốn** chuyển tiền từ quỹ này sang quỹ khác
**Để** phân bổ vốn theo nhu cầu

**Acceptance Criteria:**
- Chọn quỹ nguồn và quỹ đích (không được trùng nhau)
- Chọn loại tiền: tiền mặt / chuyển khoản
- Nhập số tiền (validate ≤ số dư loại tiền đó của quỹ nguồn)
- Nhập lý do
- Tạo atomic: 1 withdrawal từ quỹ nguồn + 1 deposit vào quỹ đích, cross-reference lẫn nhau
- Hiển thị 2 giao dịch liên kết trong lịch sử

---

## Data Flow

```
Order completed
    ↓
Waiter bàn giao cho Cashier (cash_handover)
    ↓ tiền mặt vào ngăn kéo
Cashier kết ca → bàn giao cho Manager (fund_handover)
    ↓ FundTransaction(type=fund_handover, fund_type=cash_drawer)

Manager phân bổ:
    ├─ Chuyển sang Quỹ vận hành   → fund_transfer
    ├─ Chuyển sang Quỹ nguyên liệu → fund_transfer
    └─ Chuyển sang Quỹ lợi nhuận  → fund_transfer

Chi phí phát sinh:
    ├─ Expense "Chi từ quỹ" → withdrawal, fund_type=operating, source_type=expense
    └─ Mua nguyên liệu      → withdrawal, fund_type=inventory, source_type=ingredient
```
