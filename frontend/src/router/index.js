import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import MenuView from '../views/MenuView.vue'
import IngredientView from '../views/IngredientView.vue'
import FacilityView from '../views/FacilityView.vue'
import ExpenseView from '../views/ExpenseView.vue'
import OrderView from '../views/OrderView.vue'
import ShiftView from '../views/ShiftView.vue'
import CashierDashboard from '../views/CashierDashboard.vue'
import CashierReports from '../views/CashierReports.vue'
import UserManagementView from '../views/UserManagementView.vue'
import ProfileView from '../views/ProfileView.vue'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { requiresGuest: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { requiresAuth: true }
  },
  {
    path: '/menu',
    name: 'Menu',
    component: MenuView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/ingredients',
    name: 'Ingredients',
    component: IngredientView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/facilities',
    name: 'Facilities',
    component: FacilityView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/expenses',
    name: 'Expenses',
    component: ExpenseView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/orders',
    name: 'Orders',
    component: OrderView,
    meta: { requiresAuth: true }
  },
  {
    path: '/shifts',
    name: 'Shifts',
    component: ShiftView,
    meta: { requiresAuth: true }
  },
  {
    path: '/users',
    name: 'UserManagement',
    component: UserManagementView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: ProfileView,
    meta: { requiresAuth: true }
  },
  {
    path: '/cashier',
    name: 'CashierDashboard',
    component: CashierDashboard,
    meta: { requiresAuth: true, requiresCashier: true }
  },
  {
    path: '/cashier/reports',
    name: 'CashierReports',
    component: CashierReports,
    meta: { requiresAuth: true, requiresCashier: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/dashboard')
  } else if (to.meta.requiresManager && authStore.user?.role !== 'manager') {
    next('/dashboard')
  } else if (to.meta.requiresCashier && !['cashier', 'manager'].includes(authStore.user?.role)) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router