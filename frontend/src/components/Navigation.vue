<template>
  <nav class="navbar">
    <div class="nav-brand">
      <router-link to="/dashboard" class="brand-link">
        ☕ Café POS
      </router-link>
    </div>
    
    <!-- Mobile hamburger button -->
    <button class="mobile-menu-btn" @click="toggleMobileMenu" :class="{ active: isMobileMenuOpen }">
      <span></span>
      <span></span>
      <span></span>
    </button>
    
    <!-- Navigation overlay for mobile -->
    <div class="nav-overlay" :class="{ active: isMobileMenuOpen }" @click="closeMobileMenu"></div>
    
    <!-- Navigation menu -->
    <div class="nav-menu" :class="{ active: isMobileMenuOpen }">
      <div class="nav-links">
        <router-link to="/dashboard" class="nav-link" @click="closeMobileMenu">
          🏠 Dashboard
        </router-link>
        
        <div v-if="userRole === 'manager'" class="manager-links">
          <router-link to="/menu" class="nav-link" @click="closeMobileMenu">
            🍽️ Menu
          </router-link>
          <router-link to="/ingredients" class="nav-link" @click="closeMobileMenu">
            🥬 Nguyên liệu
          </router-link>
          <router-link to="/facilities" class="nav-link" @click="closeMobileMenu">
            🏢 Cơ sở vật chất
          </router-link>
          <router-link to="/expenses" class="nav-link" @click="closeMobileMenu">
            💰 Chi phí
          </router-link>
        </div>
      </div>

      <div class="nav-user">
        <span class="user-name">{{ userName }}</span>
        <button @click="logout" class="logout-btn">Đăng xuất</button>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const isMobileMenuOpen = ref(false)

const userRole = computed(() => authStore.user?.role)
const userName = computed(() => authStore.user?.name)

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}

const logout = () => {
  authStore.logout()
  router.push('/login')
  closeMobileMenu()
}

// Close mobile menu on window resize
const handleResize = () => {
  if (window.innerWidth > 768) {
    closeMobileMenu()
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.navbar {
  background: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  position: relative;
  z-index: 1000;
}

.nav-brand .brand-link {
  font-size: 20px;
  font-weight: bold;
  color: #667eea;
  text-decoration: none;
}

/* Mobile hamburger button */
.mobile-menu-btn {
  display: none;
  flex-direction: column;
  background: none;
  border: none;
  cursor: pointer;
  padding: 5px;
  z-index: 1001;
}

.mobile-menu-btn span {
  width: 25px;
  height: 3px;
  background: #333;
  margin: 3px 0;
  transition: 0.3s;
  border-radius: 2px;
}

.mobile-menu-btn.active span:nth-child(1) {
  transform: rotate(-45deg) translate(-5px, 6px);
}

.mobile-menu-btn.active span:nth-child(2) {
  opacity: 0;
}

.mobile-menu-btn.active span:nth-child(3) {
  transform: rotate(45deg) translate(-5px, -6px);
}

/* Navigation overlay */
.nav-overlay {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  z-index: 998;
}

.nav-overlay.active {
  display: block;
}

/* Navigation menu */
.nav-menu {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex: 1;
  margin: 0 20px;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 20px;
}

.manager-links {
  display: flex;
  gap: 20px;
}

.nav-link {
  color: #666;
  text-decoration: none;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.2s;
  white-space: nowrap;
}

.nav-link:hover {
  background: #f5f6fa;
  color: #333;
}

.nav-link.router-link-active {
  background: #667eea;
  color: white;
}

.nav-user {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-name {
  color: #666;
  font-weight: 500;
}

.logout-btn {
  background: #e74c3c;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  white-space: nowrap;
}

.logout-btn:hover {
  background: #c0392b;
}

/* Mobile styles */
@media (max-width: 768px) {
  .navbar {
    padding: 0 15px;
  }
  
  .mobile-menu-btn {
    display: flex;
  }
  
  .nav-menu {
    position: fixed;
    top: 60px;
    right: -100%;
    width: 280px;
    height: calc(100vh - 60px);
    background: white;
    box-shadow: -2px 0 10px rgba(0,0,0,0.1);
    transition: right 0.3s ease;
    flex-direction: column;
    justify-content: flex-start;
    padding: 20px;
    margin: 0;
    z-index: 999;
  }
  
  .nav-menu.active {
    right: 0;
  }
  
  .nav-links {
    flex-direction: column;
    align-items: stretch;
    gap: 0;
    width: 100%;
    margin-bottom: 30px;
  }
  
  .manager-links {
    flex-direction: column;
    gap: 0;
    width: 100%;
  }
  
  .nav-link {
    padding: 15px 20px;
    border-radius: 8px;
    margin-bottom: 5px;
    font-size: 16px;
    border: 1px solid #eee;
  }
  
  .nav-user {
    flex-direction: column;
    align-items: stretch;
    gap: 15px;
    width: 100%;
    padding-top: 20px;
    border-top: 1px solid #eee;
  }
  
  .user-name {
    text-align: center;
    font-size: 16px;
    padding: 10px;
    background: #f8f9fa;
    border-radius: 8px;
  }
  
  .logout-btn {
    padding: 12px 20px;
    font-size: 16px;
    width: 100%;
  }
}

@media (max-width: 480px) {
  .nav-brand .brand-link {
    font-size: 18px;
  }
  
  .nav-menu {
    width: 100%;
    right: -100%;
  }
  
  .nav-menu.active {
    right: 0;
  }
}
</style>