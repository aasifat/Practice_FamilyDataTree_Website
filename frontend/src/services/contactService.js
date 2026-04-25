import api from './api';

export const sendContactMessage = async (messageData) => {
  const response = await api.post('/contact', messageData);
  return response.data;
};

export const getAllMessages = async () => {
  const response = await api.get('/admin/messages');
  return response.data;
};
