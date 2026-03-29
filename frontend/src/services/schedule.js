import api from './api'

// ─── Templates (manager) ──────────────────────────────────────────────────────

export const getTemplates = async () => {
  const res = await api.get('/manager/schedule/templates')
  return res.data
}

export const createTemplate = async (data) => {
  const res = await api.post('/manager/schedule/templates', data)
  return res.data
}

export const updateTemplate = async (id, data) => {
  const res = await api.put(`/manager/schedule/templates/${id}`, data)
  return res.data
}

export const deleteTemplate = async (id) => {
  const res = await api.delete(`/manager/schedule/templates/${id}`)
  return res.data
}

// ─── Periods (manager) ────────────────────────────────────────────────────────

export const getPeriods = async () => {
  const res = await api.get('/manager/schedule/periods')
  return res.data
}

export const createPeriod = async (data) => {
  const res = await api.post('/manager/schedule/periods', data)
  return res.data
}

export const setPeriodStatus = async (id, status) => {
  const res = await api.patch(`/manager/schedule/periods/${id}/status`, { status })
  return res.data
}

export const getPeriodDetail = async (id) => {
  const res = await api.get(`/manager/schedule/periods/${id}`)
  return res.data
}

// ─── Slots (manager) ─────────────────────────────────────────────────────────

export const addSlot = async (periodId, templateId, date) => {
  const res = await api.post(`/manager/schedule/periods/${periodId}/slots`, {
    template_id: templateId,
    date,
  })
  return res.data
}

export const removeSlot = async (periodId, slotId) => {
  const res = await api.delete(`/manager/schedule/periods/${periodId}/slots/${slotId}`)
  return res.data
}

export const cancelRegistration = async (regId) => {
  const res = await api.delete(`/manager/schedule/registrations/${regId}`)
  return res.data
}

// ─── Staff ────────────────────────────────────────────────────────────────────

export const getCurrentPeriod = async () => {
  const res = await api.get('/schedule/current')
  return res.data
}

export const getMySchedule = async () => {
  const res = await api.get('/schedule/my')
  return res.data
}

export const registerSlot = async (slotId) => {
  const res = await api.post(`/schedule/slots/${slotId}/register`)
  return res.data
}

export const unregisterSlot = async (slotId) => {
  const res = await api.delete(`/schedule/slots/${slotId}/register`)
  return res.data
}
