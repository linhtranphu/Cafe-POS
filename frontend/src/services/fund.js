import api from './api'

// Get current fund balance
export const getBalance = async () => {
  const response = await api.get('/manager/fund/balance')
  return response.data
}

// Get transaction history
export const getTransactions = async (filters = {}) => {
  const params = new URLSearchParams()
  
  if (filters.type) params.append('type', filters.type)
  if (filters.money_type) params.append('money_type', filters.money_type)
  if (filters.from_date) params.append('from_date', filters.from_date)
  if (filters.to_date) params.append('to_date', filters.to_date)
  if (filters.limit) params.append('limit', filters.limit)
  if (filters.offset) params.append('offset', filters.offset)
  
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
export const getTransactionDetail = async (id) => {
  const response = await api.get(`/manager/fund/transactions/${id}`)
  return response.data
}
