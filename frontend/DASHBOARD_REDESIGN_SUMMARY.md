# Dashboard & Bottom Nav Redesign for Manager Role

## Thay đổi

### Bottom Nav (BottomNav.vue)
**Trước:** 8 items (scrollable)
- Dashboard, Quản lý ca, Báo cáo, Chi phí món, Lợi nhuận, Nhân viên, Cài đặt, Cá nhân

**Sau:** 5 items (vừa màn hình)
- 🏠 Dashboard
- 📊 Báo cáo (Cashier Reports)
- 📈 Lợi nhuận (Profit Analysis)
- ⚙️ Cài đặt
- 👤 Cá nhân

### Dashboard (DashboardView.vue)
Các mục đã bỏ khỏi bottom nav được tổ chức lại trong dashboard theo nhóm:

#### 📊 Báo cáo & Phân tích (4 items)
- 📈 Phân tích lợi nhuận - `/manager/profit-analysis`
- 💰 Chi phí món - `/manager/menu-costs`
- 📊 Báo cáo thu ngân - `/cashier/reports`
- ⏰ Quản lý ca - `/manager/shifts`

#### 🍽️ Menu & Nguyên liệu (2 items)
- 🍽️ Menu - `/menu`
- 🥬 Nguyên liệu - `/ingredients`

#### 💸 Chi phí & Tài sản (2 items)
- 💸 Chi phí - `/expenses`
- 🏢 Cơ sở vật chất - `/facilities`

#### 👥 Nhân sự (1 item)
- 👥 Nhân viên - `/users`

## Lợi ích

1. **Bottom nav gọn gàng hơn**: Chỉ 5 items quan trọng nhất, không cần scroll
2. **Dashboard có tổ chức**: Các chức năng được nhóm logic theo mục đích sử dụng
3. **Dễ tìm kiếm**: Manager có thể nhanh chóng tìm chức năng theo nhóm
4. **UX tốt hơn**: Giảm cognitive load, tăng hiệu quả sử dụng

## Files đã thay đổi

1. `frontend/src/components/BottomNav.vue` - Giảm items cho manager từ 8 xuống 5
2. `frontend/src/views/DashboardView.vue` - Tổ chức lại quick access theo nhóm chức năng

