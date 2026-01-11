import axios, { AxiosError } from 'axios';
import { API_URL } from '../shared/constants';

// Создаем клиент с базовым URL
const client = axios.create({
  baseURL: `${API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 секунд таймаут
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
  (error: AxiosError) => {
    if (error.response) {
      // Сервер ответил с кодом ошибки
      console.error('API Error:', {
        status: error.response.status,
        statusText: error.response.statusText,
        data: error.response.data,
        url: error.config?.url,
      });
    } else if (error.request) {
      // Запрос был отправлен, но ответа не получено
      console.error('API Request Error:', {
        message: error.message,
        url: error.config?.url,
      });
    } else {
      // Ошибка при настройке запроса
      console.error('API Setup Error:', error.message);
    }
    return Promise.reject(error);
  }
);

export default client;
