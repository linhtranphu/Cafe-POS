// Facility Status Constants
// Must match backend: backend/domain/facility/facility.go

export const FACILITY_STATUS = {
  IN_USE: 'Đang sử dụng',
  BROKEN: 'Hỏng',
  REPAIRING: 'Đang sửa',
  INACTIVE: 'Ngừng sử dụng',
  DISPOSED: 'Thanh lý'
}

export const FACILITY_STATUS_OPTIONS = [
  { value: FACILITY_STATUS.IN_USE, label: 'Đang sử dụng', color: 'green' },
  { value: FACILITY_STATUS.REPAIRING, label: 'Đang sửa', color: 'yellow' },
  { value: FACILITY_STATUS.BROKEN, label: 'Hỏng', color: 'red' },
  { value: FACILITY_STATUS.INACTIVE, label: 'Ngừng sử dụng', color: 'gray' },
  { value: FACILITY_STATUS.DISPOSED, label: 'Thanh lý', color: 'gray' }
]

export const FACILITY_STATUS_CLASSES = {
  [FACILITY_STATUS.IN_USE]: 'bg-green-100 text-green-800',
  [FACILITY_STATUS.REPAIRING]: 'bg-yellow-100 text-yellow-800',
  [FACILITY_STATUS.BROKEN]: 'bg-red-100 text-red-800',
  [FACILITY_STATUS.INACTIVE]: 'bg-gray-100 text-gray-800',
  [FACILITY_STATUS.DISPOSED]: 'bg-gray-100 text-gray-800'
}

// Facility Types
export const FACILITY_TYPES = {
  FURNITURE: 'Bàn ghế',
  MACHINE: 'Máy móc',
  UTENSIL: 'Dụng cụ',
  ELECTRIC: 'Điện tử',
  OTHER: 'Khác'
}

export const FACILITY_TYPE_OPTIONS = [
  { value: FACILITY_TYPES.FURNITURE, label: 'Bàn ghế' },
  { value: FACILITY_TYPES.MACHINE, label: 'Máy móc' },
  { value: FACILITY_TYPES.UTENSIL, label: 'Dụng cụ' },
  { value: FACILITY_TYPES.ELECTRIC, label: 'Điện tử' },
  { value: FACILITY_TYPES.OTHER, label: 'Khác' }
]

// Facility Areas
export const FACILITY_AREAS = {
  DINING_ROOM: 'Phòng khách',
  KITCHEN: 'Bếp',
  COUNTER: 'Quầy bar',
  STORAGE: 'Kho',
  OFFICE: 'Văn phòng',
  OTHER: 'Khác'
}

export const FACILITY_AREA_OPTIONS = [
  { value: FACILITY_AREAS.DINING_ROOM, label: 'Phòng khách' },
  { value: FACILITY_AREAS.KITCHEN, label: 'Bếp' },
  { value: FACILITY_AREAS.COUNTER, label: 'Quầy bar' },
  { value: FACILITY_AREAS.STORAGE, label: 'Kho' },
  { value: FACILITY_AREAS.OFFICE, label: 'Văn phòng' },
  { value: FACILITY_AREAS.OTHER, label: 'Khác' }
]

// Maintenance Types
export const MAINTENANCE_TYPES = {
  SCHEDULED: 'scheduled',
  EMERGENCY: 'emergency',
  PREVENTIVE: 'preventive',
  CORRECTIVE: 'corrective'
}

export const MAINTENANCE_TYPE_OPTIONS = [
  { value: MAINTENANCE_TYPES.SCHEDULED, label: 'Định kỳ' },
  { value: MAINTENANCE_TYPES.EMERGENCY, label: 'Khẩn cấp' },
  { value: MAINTENANCE_TYPES.PREVENTIVE, label: 'Phòng ngừa' },
  { value: MAINTENANCE_TYPES.CORRECTIVE, label: 'Sửa chữa' }
]

// Issue Severity
export const ISSUE_SEVERITY = {
  LOW: 'low',
  MEDIUM: 'medium',
  HIGH: 'high',
  CRITICAL: 'critical'
}

export const ISSUE_SEVERITY_OPTIONS = [
  { value: ISSUE_SEVERITY.LOW, label: 'Thấp', color: 'blue' },
  { value: ISSUE_SEVERITY.MEDIUM, label: 'Trung bình', color: 'yellow' },
  { value: ISSUE_SEVERITY.HIGH, label: 'Cao', color: 'orange' },
  { value: ISSUE_SEVERITY.CRITICAL, label: 'Nghiêm trọng', color: 'red' }
]

export const ISSUE_SEVERITY_CLASSES = {
  [ISSUE_SEVERITY.LOW]: 'bg-blue-100 text-blue-800',
  [ISSUE_SEVERITY.MEDIUM]: 'bg-yellow-100 text-yellow-800',
  [ISSUE_SEVERITY.HIGH]: 'bg-orange-100 text-orange-800',
  [ISSUE_SEVERITY.CRITICAL]: 'bg-red-100 text-red-800'
}

// Issue Status
export const ISSUE_STATUS = {
  OPEN: 'open',
  IN_PROGRESS: 'in_progress',
  RESOLVED: 'resolved'
}

export const ISSUE_STATUS_OPTIONS = [
  { value: ISSUE_STATUS.OPEN, label: 'Chờ xử lý', color: 'yellow' },
  { value: ISSUE_STATUS.IN_PROGRESS, label: 'Đang xử lý', color: 'blue' },
  { value: ISSUE_STATUS.RESOLVED, label: 'Đã giải quyết', color: 'green' }
]

export const ISSUE_STATUS_CLASSES = {
  [ISSUE_STATUS.OPEN]: 'bg-yellow-100 text-yellow-800',
  [ISSUE_STATUS.IN_PROGRESS]: 'bg-blue-100 text-blue-800',
  [ISSUE_STATUS.RESOLVED]: 'bg-green-100 text-green-800'
}

// Helper functions
export function getFacilityStatusClass(status) {
  return FACILITY_STATUS_CLASSES[status] || 'bg-gray-100 text-gray-800'
}

export function getIssueStatusClass(status) {
  return ISSUE_STATUS_CLASSES[status] || 'bg-gray-100 text-gray-800'
}

export function getIssueSeverityClass(severity) {
  return ISSUE_SEVERITY_CLASSES[severity] || 'bg-gray-100 text-gray-800'
}
