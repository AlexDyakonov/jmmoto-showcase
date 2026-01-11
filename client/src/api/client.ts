import axios from 'axios';
import { API_URL } from '../shared/constants';

const client = axios.create({
  baseURL: `${API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Добавляем токен Telegram в заголовки
export const setTelegramToken = (token: string) => {
  if (token) {
    client.defaults.headers.common['X-API-Token'] = token;
  } else {
    delete client.defaults.headers.common['X-API-Token'];
  }
};

// Обработка ошибок
client.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export default client;
