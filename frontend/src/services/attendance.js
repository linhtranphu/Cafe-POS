import api from './api'

export const getEntries = async (from, to) => {
  const params = {}
  if (from) params.from = from
  if (to) params.to = to
  const res = await api.get('/manager/attendance', { params })
  return res.data
}

export const addEntry = async (data) => {
  const res = await api.post('/manager/attendance', data)
  return res.data
}

export const updateEntry = async (id, data) => {
  const res = await api.patch(`/manager/attendance/${id}`, data)
  return res.data
}

export const bulkAddEntries = async (entries) => {
  const res = await api.post('/manager/attendance/bulk', entries)
  return res.data
}

export const deleteEntry = async (id) => {
  const res = await api.delete(`/manager/attendance/${id}`)
  return res.data
}
