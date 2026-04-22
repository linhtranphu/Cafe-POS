import api from './api'

export const aiService = {
  async parseCommand(message, conversationHistory = []) {
    const response = await api.post('/manager/ai/parse', {
      message,
      conversation_history: conversationHistory,
    })
    return response.data // { reply_text, action: { type, fields } | null }
  },

  async getHistory() {
    const response = await api.get('/manager/ai/history')
    return response.data.logs // AICommandLog[]
  },
}
