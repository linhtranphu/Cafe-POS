// User Constants
// Must match backend: backend/domain/user/user.go

// User Roles
export const USER_ROLES = {
  ADMIN: 'admin',
  MANAGER: 'manager',
  CASHIER: 'cashier',
  WAITER: 'waiter',
  BARISTA: 'barista'
}

// User Role Options for UI
export const USER_ROLE_OPTIONS = [
  { value: USER_ROLES.ADMIN, label: 'Admin', icon: '👑', badge: 'bg-red-100 text-red-800' },
  { value: USER_ROLES.MANAGER, label: 'Quản lý', icon: '👨‍💼', badge: 'bg-purple-100 text-purple-800' },
  { value: USER_ROLES.CASHIER, label: 'Thu ngân', icon: '💰', badge: 'bg-green-100 text-green-800' },
  { value: USER_ROLES.WAITER, label: 'Phục vụ', icon: '👨‍🍳', badge: 'bg-blue-100 text-blue-800' },
  { value: USER_ROLES.BARISTA, label: 'Pha chế', icon: '☕', badge: 'bg-yellow-100 text-yellow-800' }
]

// Helper Functions
export function getUserRoleLabel(role) {
  const option = USER_ROLE_OPTIONS.find(opt => opt.value === role)
  return option ? option.label : role
}

export function getUserRoleIcon(role) {
  const option = USER_ROLE_OPTIONS.find(opt => opt.value === role)
  return option ? option.icon : '👤'
}

export function getUserRoleBadge(role) {
  const option = USER_ROLE_OPTIONS.find(opt => opt.value === role)
  return option ? option.badge : 'bg-gray-100 text-gray-800'
}

// Role Permissions (for frontend display logic)
export const ROLE_PERMISSIONS = {
  [USER_ROLES.ADMIN]: {
    canManageUsers: true,
    canManageMenu: true,
    canManageIngredients: true,
    canManageFacilities: true,
    canManageExpenses: true,
    canViewReports: true,
    canManageShifts: true
  },
  [USER_ROLES.MANAGER]: {
    canManageUsers: true,
    canManageMenu: true,
    canManageIngredients: true,
    canManageFacilities: true,
    canManageExpenses: true,
    canViewReports: true,
    canManageShifts: true
  },
  [USER_ROLES.CASHIER]: {
    canManageUsers: false,
    canManageMenu: false,
    canManageIngredients: false,
    canManageFacilities: false,
    canManageExpenses: false,
    canViewReports: false,
    canManageShifts: true
  },
  [USER_ROLES.WAITER]: {
    canManageUsers: false,
    canManageMenu: false,
    canManageIngredients: false,
    canManageFacilities: false,
    canManageExpenses: false,
    canViewReports: false,
    canManageShifts: true
  },
  [USER_ROLES.BARISTA]: {
    canManageUsers: false,
    canManageMenu: false,
    canManageIngredients: false,
    canManageFacilities: false,
    canManageExpenses: false,
    canViewReports: false,
    canManageShifts: true
  }
}

export function hasPermission(role, permission) {
  const permissions = ROLE_PERMISSIONS[role]
  return permissions ? permissions[permission] : false
}
