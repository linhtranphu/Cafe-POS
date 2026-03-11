import api from './api'

// Get balances for all 5 funds from journal (double-entry source of truth)
export const getJournalBalances = async () => {
  const response = await api.get('/manager/fund/journal-balances')
  return response.data
}

// List journal entries with optional filters
// filters: { fund_type, event_type, from_date, to_date, limit, offset }
export const getJournalEntries = async (filters = {}) => {
  const params = new URLSearchParams()
  if (filters.fund_type) params.append('fund_type', filters.fund_type)
  if (filters.event_type) params.append('event_type', filters.event_type)
  if (filters.from_date) params.append('from_date', filters.from_date)
  if (filters.to_date) params.append('to_date', filters.to_date)
  if (filters.limit) params.append('limit', filters.limit)
  if (filters.offset) params.append('offset', filters.offset)
  const response = await api.get(`/manager/fund/journal?${params}`)
  return response.data
}

// Get single journal entry by ID
export const getJournalEntry = async (id) => {
  const response = await api.get(`/manager/fund/journal/${id}`)
  return response.data
}
