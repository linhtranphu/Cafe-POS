<template>
  <div class="expense-management">
    <Navigation />
    <div class="content">
      <div class="header">
        <h2>💰 Quản lý Chi phí</h2>
        <button class="btn-primary">+ Thêm chi phí</button>
      </div>
      
      <div class="expense-grid">
        <div class="expense-card">
          <h3>Đang phát triển...</h3>
          <p>Chức năng quản lý chi phí đang được phát triển</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import Navigation from '../components/Navigation.vue'
</script>

<style scoped>
.expense-management {
  min-height: 100vh;
  background: #f5f6fa;
}

.content {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.header h2 {
  color: #333;
  margin: 0;
}

.btn-primary {
  background: #667eea;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 6px;
  cursor: pointer;
}

.expense-grid {
  display: grid;
  gap: 20px;
}

.expense-card {
  background: white;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  text-align: center;
}

.expense-card h3 {
  color: #333;
  margin-bottom: 10px;
}

.expense-card p {
  color: #666;
}
</style>