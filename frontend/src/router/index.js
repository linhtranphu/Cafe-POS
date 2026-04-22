import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import MenuView from '../views/MenuView.vue'
import IngredientManagementView from '../views/IngredientManagementView.vue'
import FacilityManagementView from '../views/FacilityManagementView.vue'
import FacilityAddEditView from '../views/FacilityAddEditView.vue'
import ExpenseManagementView from '../views/ExpenseManagementView.vue'
import OrderView from '../views/OrderView.vue'
import CreateOrderView from '../views/CreateOrderView.vue'
import BaristaView from '../views/BaristaView.vue'
import ShiftView from '../views/ShiftView.vue'
import ManagerShiftView from '../views/ManagerShiftView.vue'
import CashierDashboard from '../views/CashierDashboard.vue'
import CashierReports from '../views/CashierReports.vue'
import CashierShiftClosureV2 from '../views/CashierShiftClosureV2.vue'
import UserManagementView from '../views/UserManagementView.vue'
import ProfileView from '../views/ProfileView.vue'
import MenuCostView from '../views/MenuCostView.vue'
import ProfitAnalysisView from '../views/ProfitAnalysisView.vue'
import CostAnalysisView from '../views/CostAnalysisView.vue'
import SettingsView from '../views/SettingsView.vue'
import BatchManagementView from '../views/BatchManagementView.vue'
import PrintManagementView from '../views/PrintManagementView.vue'
import InventoryStatsView from '../views/InventoryStatsView.vue'

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
    path: '/facilities/create',
    name: 'FacilityCreate',
    component: FacilityAddEditView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/facilities/edit/:id',
    name: 'FacilityEdit',
    component: FacilityAddEditView,
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
    path: '/create-order',
    name: 'CreateOrder',
    component: CreateOrderView,
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
    component: CashierShiftClosureV2,
    meta: { requiresAuth: true, requiresCashier: true }
  },
  {
    path: '/cashier/handovers',
    name: 'CashierHandovers',
    component: () => import('../views/CashierHandoverView.vue'),
    meta: { requiresAuth: true, requiresCashier: true }
  },
  {
    path: '/manager/menu-costs',
    name: 'MenuCosts',
    component: MenuCostView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/manager/profit-analysis',
    name: 'ProfitAnalysis',
    component: ProfitAnalysisView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/manager/cost-analysis',
    name: 'CostAnalysis',
    component: CostAnalysisView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: SettingsView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/batch',
    name: 'BatchManagement',
    component: BatchManagementView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/batch/definitions',
    name: 'BatchDefinitions',
    component: () => import('../views/BatchDefinitionsView.vue'),
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/batch/definitions/create',
    name: 'BatchDefinitionCreate',
    component: () => import('../views/BatchDefinitionFormView.vue'),
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/batch/definitions/:id',
    name: 'BatchDefinitionEdit',
    component: () => import('../views/BatchDefinitionFormView.vue'),
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/batch/records',
    name: 'BatchRecords',
    component: () => import('../views/BatchRecordsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/batch/records/create',
    name: 'BatchRecordCreate',
    component: () => import('../views/BatchRecordFormView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/batch/records/:id',
    name: 'BatchRecordDetail',
    component: () => import('../views/BatchRecordDetailView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/batch/alerts',
    name: 'BatchAlerts',
    component: () => import('../views/BatchAlertsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/batch/reports',
    name: 'BatchReports',
    component: () => import('../views/BatchReportsView.vue'),
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/print-management',
    name: 'PrintManagement',
    component: PrintManagementView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/manager/fund',
    name: 'FundManagement',
    component: () => import('../views/FundManagementView.vue'),
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/operating-expenses',
    name: 'OperatingExpenses',
    component: () => import('../views/OperatingExpenseView.vue'),
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/manager/inventory-stats',
    name: 'InventoryStats',
    component: InventoryStatsView,
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/manager/schedule',
    name: 'ScheduleManagement',
    component: () => import('../views/ScheduleManagementView.vue'),
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/schedule',
    name: 'Schedule',
    component: () => import('../views/ScheduleView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/manager/attendance',
    component: () => import('../views/AttendanceView.vue'),
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/manager/attendance/summary',
    name: 'AttendanceSummary',
    component: () => import('../views/AttendanceSummaryView.vue'),
    meta: { requiresAuth: true, requiresManager: true }
  },
  {
    path: '/manager/ai-command',
    name: 'AICommand',
    component: () => import('../views/AICommandView.vue'),
    meta: { requiresAuth: true, requiresManager: true }
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    // If using browser back/forward, restore the saved position
    if (savedPosition) {
      return savedPosition
    }
    
    // For hash links (e.g., #section)
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    
    // Default: scroll to top for new navigation
    // But we'll handle this manually in components for better control
    return { top: 0 }
  }
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Ensure auth is initialized before checking
  if (!authStore._initialized) {
    authStore.initAuth()
    authStore._initialized = true
  }
  
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