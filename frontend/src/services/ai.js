import apiClient from './deploy';

export const aiService = {
  // 发送聊天消息
  chat: async (data) => {
    const response = await apiClient.post('/ai/chat', data);
    return response;
  },

  // 获取会话列表
  getSessions: async (userId = 'default') => {
    const response = await apiClient.get('/ai/sessions', {
      params: { user_id: userId },
    });
    return response;
  },

  // 获取历史消息
  getMessages: async (sessionId) => {
    const response = await apiClient.get(`/ai/messages/${sessionId}`);
    return response;
  },
};

export default aiService;
