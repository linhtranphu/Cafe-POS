import api from './api'

// Get current fund balance (operating fund — for backward compat)
export const getBalance = async () => {
  const response = await api.get('/manager/fund/journal-balances')
  return response.data?.operating || { cash: 0, transfer: 0, total: 0 }
}

// Get transaction history
// Updated to include source info (Requirements: 7.7, 7.8)
export const getTransactions = async (filters = {}) => {
  const params = new URLSearchParams()
  
  if (filters.type) params.append('type', filters.type)
  if (filters.money_type) params.append('money_type', filters.money_type)
  if (filters.from_date) params.append('from_date', filters.from_date)
  if (filters.to_date) params.append('to_date', filters.to_date)
  if (filters.limit) params.append('limit', filters.limit)
  if (filters.offset) params.append('offset', filters.offset)
  if (filters.source_type) params.append('source_type', filters.source_type)
  if (filters.fund_type) params.append('fund_type', filters.fund_type)

  const response = await api.get(`/manager/fund/transactions?${params}`)
  return response.data
}

// Create deposit
export const createDeposit = async (data) => {
  const response = await api.post('/manager/fund/deposit', data)
  return response.data
}

// Create withdrawal
export const createWithdrawal = async (data) => {
  const response = await api.post('/manager/fund/withdraw', data)
  return response.data
}

// Get transaction detail
// Updated to include source details (Requirements: 7.7, 7.8)
export const getTransactionDetail = async (id) => {
  const response = await api.get(`/manager/fund/transactions/${id}`)
  return response.data
}

// Get all fund balances from journal (double-entry source of truth)
export const getAllBalances = async () => {
  const response = await api.get('/manager/fund/journal-balances')
  return response.data
}

// Transfer between fund accounts
export const transferBetweenFunds = async (data) => {
  const response = await api.post('/manager/fund/transfer', data)
  return response.data
}

// Get daily fund report
export const getDailyReport = async (date) => {
  const response = await api.get(`/manager/fund/daily-report?date=${date}`)
  return response.data
}
