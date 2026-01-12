export interface Motorcycle {
  id: string;
  title: string;
  price: number;
  currency: string;
  description?: string;
  status: 'available' | 'reserved' | 'sold' | 'draft';
  sourceUrl: string;
  photos?: MotorcyclePhoto[];
  createdAt: string;
  updatedAt: string;
}

export interface MotorcyclePhoto {
  id: string;
  motorcycleId: string;
  s3Url: string;
  order: number;
  createdAt: string;
}

export interface FilterMotorcycle {
  status?: 'available' | 'reserved' | 'sold' | 'draft';
  title?: string;
  minPrice?: number;
  maxPrice?: number;
}

import { getTelegramInitData } from '../utils/telegram';

// Получаем API URL из глобальной конфигурации
const getApiBaseUrl = () => {
  return `${window.api.API_URL}/api/v1`;
};

// Создаем заголовки для API запросов с Telegram initData
const createApiHeaders = (): HeadersInit => {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };
  
  // Добавляем Telegram initData если доступно
  const initData = getTelegramInitData();
  if (initData) {
    headers['X-API-Token'] = initData;
    console.log('Added X-API-Token header with initData');
  } else {
    console.log('No Telegram initData available, sending anonymous request');
  }
  
  return headers;
};

export const getMotorcycles = async (filters?: FilterMotorcycle): Promise<Motorcycle[]> => {
  try {
    // Строим URL с параметрами
    let url = `${getApiBaseUrl()}/motorcycles`;
    const searchParams = new URLSearchParams();
    
    if (filters?.status) {
      searchParams.append('status', filters.status);
    }
    if (filters?.title && filters.title.trim()) {
      searchParams.append('title', filters.title.trim());
    }
    if (filters?.minPrice && filters.minPrice > 0) {
      searchParams.append('minPrice', filters.minPrice.toString());
    }
    if (filters?.maxPrice && filters.maxPrice > 0) {
      searchParams.append('maxPrice', filters.maxPrice.toString());
    }
    
    const queryString = searchParams.toString();
    if (queryString) {
      url += `?${queryString}`;
    }
    
    console.log('Fetching motorcycles from:', url);
    
    // Fetch запрос с Telegram initData
    const response = await fetch(url, {
      method: 'GET',
      headers: createApiHeaders(),
    });
    
    if (!response.ok) {
      if (response.status === 401) {
        throw new Error('Требуется авторизация через Telegram');
      }
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    
    const data = await response.json();
    console.log('Raw API response:', data);
    
    // Обрабатываем ответ - Huma возвращает { body: [...] }
    if (data && typeof data === 'object' && 'body' in data && Array.isArray(data.body)) {
      console.log('Found motorcycles in body:', data.body.length);
      return data.body;
    }
    
    // Fallback: если данные приходят напрямую как массив
    if (Array.isArray(data)) {
      console.log('Found motorcycles as direct array:', data.length);
      return data;
    }
    
    console.warn('Unexpected response format:', data);
    return [];
    
  } catch (error) {
    console.error('Error fetching motorcycles:', error);
    throw error;
  }
};

export const getMotorcycle = async (id: string): Promise<Motorcycle> => {
  try {
    const url = `${getApiBaseUrl()}/motorcycles/${id}`;
    console.log('Fetching motorcycle from:', url);
    
    const response = await fetch(url, {
      method: 'GET',
      headers: createApiHeaders(),
    });
    
    if (!response.ok) {
      if (response.status === 401) {
        throw new Error('Требуется авторизация через Telegram');
      }
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    
    const data = await response.json();
    
    // Обрабатываем ответ - Huma возвращает { body: {...} }
    if (data && typeof data === 'object' && 'body' in data) {
      return data.body;
    }
    
    // Fallback: если данные приходят напрямую
    return data;
    
  } catch (error) {
    console.error('Error fetching motorcycle:', error);
    throw error;
  }
};