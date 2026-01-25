<template>
  <div class="dashboard">
    <Navigation />
    <div class="content">
      <div class="welcome-card">
        <h2>Chào mừng đến với hệ thống POS</h2>
        <p>Bạn đã đăng nhập thành công với quyền {{ user?.role === 'manager' ? 'Quản lý' : user?.role === 'waiter' ? 'Nhân viên' : 'Thu ngân' }}</p>
        
        <div class="role-permissions">
          <h3>Quyền hạn của bạn:</h3>
          <ul>
            <li v-for="permission in permissions" :key="permission">{{ permission }}</li>
          </ul>
          
          <div v-if="user?.role === 'manager'" class="manager-actions">
            <h3>Chức năng quản lý:</h3>
            <button @click="$router.push('/menu')" class="action-btn">🍽️ Quản lý Menu</button>
            <button @click="$router.push('/ingredients')" class="action-btn">🥬 Quản lý Nguyên liệu</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import Navigation from '../components/Navigation.vue'

const authStore = useAuthStore()

const user = computed(() => authStore.user)

const permissions = computed(() => {
  switch (user.value?.role) {
    case 'manager':
      return [
        'Xem tất cả order',
        'Xem báo cáo doanh thu',
        'In lại bill',
        'Chỉnh sửa/hủy order đã thanh toán',
        'Quản lý menu & giá',
        'Quản lý bàn',
        'Quản lý user'
      ]
    case 'waiter':
      return [
        'Tạo order mới',
        'Nhập món, số lượng',
        'Gắn bàn',
        'Xem & thông báo tổng tiền',
        'Chọn phương thức thanh toán',
        'Xác nhận đã thu tiền',
        'In bill'
      ]
    case 'cashier':
      return [
        'Xem order đã tạo',
        'Thu tiền',
        'Xác nhận thanh toán',
        'In/in lại bill'
      ]
    default:
      return []
  }
})
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background: #f5f6fa;
}

.content {
  padding: 30px;
}

.welcome-card {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  max-width: 600px;
}

.welcome-card h2 {
  color: #333;
  margin-bottom: 10px;
  font-size: 24px;
}

.welcome-card p {
  color: #666;
  margin-bottom: 25px;
  line-height: 1.5;
}

.role-permissions h3 {
  color: #333;
  margin-bottom: 15px;
  font-size: 18px;
}

.role-permissions ul {
  list-style: none;
  padding: 0;
}

.role-permissions li {
  padding: 12px 0;
  color: #555;
  border-bottom: 1px solid #eee;
  font-size: 15px;
}

.role-permissions li:before {
  content: "✓ ";
  color: #27ae60;
  font-weight: bold;
  margin-right: 8px;
}

.role-permissions li:last-child {
  border-bottom: none;
}

.manager-actions {
  margin-top: 30px;
}

.manager-actions h3 {
  color: #333;
  margin-bottom: 15px;
  font-size: 18px;
}

.action-btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  margin-right: 10px;
  margin-bottom: 10px;
  transition: all 0.2s;
  min-width: 160px;
}

.action-btn:hover {
  background: #5a6fd8;
  transform: translateY(-1px);
}

/* Mobile responsive styles */
@media (max-width: 768px) {
  .content {
    padding: 15px;
  }
  
  .welcome-card {
    padding: 20px;
    margin: 0;
    max-width: 100%;
    border-radius: 8px;
  }
  
  .welcome-card h2 {
    font-size: 20px;
    text-align: center;
  }
  
  .welcome-card p {
    text-align: center;
    font-size: 14px;
  }
  
  .role-permissions h3 {
    font-size: 16px;
    text-align: center;
  }
  
  .role-permissions li {
    padding: 10px 0;
    font-size: 14px;
  }
  
  .manager-actions {
    margin-top: 25px;
  }
  
  .manager-actions h3 {
    font-size: 16px;
    text-align: center;
  }
  
  .action-btn {
    width: 100%;
    margin-right: 0;
    margin-bottom: 12px;
    padding: 15px 20px;
    font-size: 15px;
    min-width: auto;
  }
}

@media (max-width: 480px) {
  .content {
    padding: 10px;
  }
  
  .welcome-card {
    padding: 15px;
  }
  
  .welcome-card h2 {
    font-size: 18px;
  }
  
  .welcome-card p {
    font-size: 13px;
  }
  
  .role-permissions h3,
  .manager-actions h3 {
    font-size: 15px;
  }
  
  .role-permissions li {
    font-size: 13px;
    padding: 8px 0;
  }
  
  .action-btn {
    padding: 12px 15px;
    font-size: 14px;
  }
}
</style>