<template>
  <nav class="bg-white shadow-md px-4 lg:px-5 flex items-center justify-between h-16 relative z-50">
    <div class="nav-brand">
      <router-link to="/dashboard" class="text-xl font-bold text-blue-600 no-underline">
        ☕ Café POS
      </router-link>
    </div>
    
    <!-- Mobile hamburger button -->
    <button 
      class="lg:hidden flex flex-col bg-transparent border-none cursor-pointer p-1 z-50" 
      @click="toggleMobileMenu" 
      :class="{ 'active': isMobileMenuOpen }"
    >
      <span class="w-6 h-0.5 bg-gray-800 my-0.5 transition-all duration-300 rounded-sm" :class="{ 'rotate-45 translate-y-2': isMobileMenuOpen }"></span>
      <span class="w-6 h-0.5 bg-gray-800 my-0.5 transition-all duration-300 rounded-sm" :class="{ 'opacity-0': isMobileMenuOpen }"></span>
      <span class="w-6 h-0.5 bg-gray-800 my-0.5 transition-all duration-300 rounded-sm" :class="{ '-rotate-45 -translate-y-2': isMobileMenuOpen }"></span>
    </button>
    
    <!-- Navigation overlay for mobile -->
    <div 
      class="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden transition-opacity duration-300" 
      :class="{ 'opacity-100 pointer-events-auto': isMobileMenuOpen, 'opacity-0 pointer-events-none': !isMobileMenuOpen }" 
      @click="closeMobileMenu"
    ></div>
    
    <!-- Navigation menu -->
    <div class="hidden lg:flex items-center justify-between flex-1 mx-5 nav-menu" :class="{ 'active': isMobileMenuOpen }">
      <div class="flex items-center gap-5">
        <router-link to="/dashboard" class="text-gray-600 no-underline px-4 py-2 rounded-md transition-all duration-200 hover:bg-gray-100 hover:text-gray-800" @click="closeMobileMenu">
          🏠 Dashboard
        </router-link>
        
        <div v-if="userRole === 'manager'" class="flex gap-5">
          <router-link to="/menu" class="text-gray-600 no-underline px-4 py-2 rounded-md transition-all duration-200 hover:bg-gray-100 hover:text-gray-800" @click="closeMobileMenu">
            🍽️ Menu
          </router-link>
          <router-link to="/ingredients" class="text-gray-600 no-underline px-4 py-2 rounded-md transition-all duration-200 hover:bg-gray-100 hover:text-gray-800" @click="closeMobileMenu">
            🥬 Nguyên liệu
          </router-link>
          <router-link to="/facilities" class="text-gray-600 no-underline px-4 py-2 rounded-md transition-all duration-200 hover:bg-gray-100 hover:text-gray-800" @click="closeMobileMenu">
            🏢 Cơ sở vật chất
          </router-link>
          <router-link to="/expenses" class="text-gray-600 no-underline px-4 py-2 rounded-md transition-all duration-200 hover:bg-gray-100 hover:text-gray-800" @click="closeMobileMenu">
            💰 Chi phí
          </router-link>
        </div>
      </div>

      <div class="flex items-center gap-4">
        <span class="text-gray-600 font-medium">{{ userName }}</span>
        <button @click="logout" class="bg-red-600 text-white border-none px-4 py-2 rounded-md cursor-pointer text-sm hover:bg-red-700 transition-colors duration-200">
          Đăng xuất
        </button>
      </div>
    </div>
    
    <!-- Mobile menu -->
    <div class="fixed top-16 right-0 w-72 h-screen bg-white shadow-xl transition-transform duration-300 flex flex-col justify-start p-5 z-40 lg:hidden" :class="{ 'translate-x-0': isMobileMenuOpen, 'translate-x-full': !isMobileMenuOpen }">
      <div class="flex flex-col gap-0 w-full mb-8">
        <router-link to="/dashboard" class="text-gray-600 no-underline py-4 px-5 rounded-lg mb-1 text-base border border-gray-200 hover:bg-gray-100 hover:text-gray-800 transition-all duration-200" @click="closeMobileMenu">
          🏠 Dashboard
        </router-link>
        
        <div v-if="userRole === 'manager'" class="flex flex-col gap-0 w-full">
          <router-link to="/menu" class="text-gray-600 no-underline py-4 px-5 rounded-lg mb-1 text-base border border-gray-200 hover:bg-gray-100 hover:text-gray-800 transition-all duration-200" @click="closeMobileMenu">
            🍽️ Menu
          </router-link>
          <router-link to="/ingredients" class="text-gray-600 no-underline py-4 px-5 rounded-lg mb-1 text-base border border-gray-200 hover:bg-gray-100 hover:text-gray-800 transition-all duration-200" @click="closeMobileMenu">
            🥬 Nguyên liệu
          </router-link>
          <router-link to="/facilities" class="text-gray-600 no-underline py-4 px-5 rounded-lg mb-1 text-base border border-gray-200 hover:bg-gray-100 hover:text-gray-800 transition-all duration-200" @click="closeMobileMenu">
            🏢 Cơ sở vật chất
          </router-link>
          <router-link to="/expenses" class="text-gray-600 no-underline py-4 px-5 rounded-lg mb-1 text-base border border-gray-200 hover:bg-gray-100 hover:text-gray-800 transition-all duration-200" @click="closeMobileMenu">
            💰 Chi phí
          </router-link>
        </div>
      </div>
      
      <div class="flex flex-col gap-4 w-full pt-5 border-t border-gray-200">
        <div class="text-center text-base p-2 bg-gray-100 rounded-lg">
          {{ userName }}
        </div>
        <button @click="logout" class="bg-red-600 text-white border-none py-3 px-5 rounded-lg cursor-pointer text-base w-full hover:bg-red-700 transition-colors duration-200">
          Đăng xuất
        </button>
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
.router-link-active {
  @apply bg-blue-600 text-white;
}
</style>