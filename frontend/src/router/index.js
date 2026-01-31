import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import MenuView from '../views/MenuView.vue'
import IngredientManagementView from '../views/IngredientManagementView.vue'
import FacilityManagementView from '../views/FacilityManagementView.vue'
import ExpenseManagementView from '../views/ExpenseManagementView.vue'
import OrderView from '../views/OrderView.vue'
import BaristaView from '../views/BaristaView.vue'
import ShiftView from '../views/ShiftView.vue'
import ManagerShiftView from '../views/ManagerShiftView.vue'
import CashierDashboard from '../views/CashierDashboard.vue'
import CashierReports from '../views/CashierReports.vue'
import CashierShiftClosure from '../views/CashierShiftClosure.vue'
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
    component: IngredientManagementView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/facilities',
    name: 'Facilities',
    component: FacilityManagementView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/expenses',
    name: 'Expenses',
    component: ExpenseManagementView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/orders',
    name: 'Orders',
    component: OrderView,
    meta: { requiresAuth: true, requiresNotBarista: true }
  },
  {
    path: '/barista',
    name: 'Barista',
    component: BaristaView,
    meta: { requiresAuth: true, requiresBarista: true }
  },
  {
    path: '/shifts',
    name: 'Shifts',
    component: ShiftView,
    meta: { requiresAuth: true }
  },
  {
    path: '/manager/shifts',
    name: 'ManagerShifts',
    component: ManagerShiftView,
    meta: { requiresAuth: true, requiresManager: true }
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
  },
  {
    path: '/cashier/shift-closure/:id',
    name: 'CashierShiftClosure',
    component: CashierShiftClosure,
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
  } else if (to.meta.requiresBarista && authStore.user?.role !== 'barista') {
    next('/dashboard')
  } else if (to.meta.requiresNotBarista && authStore.user?.role === 'barista') {
    next('/dashboard')
  } else {
    next()
  }
})

export default router