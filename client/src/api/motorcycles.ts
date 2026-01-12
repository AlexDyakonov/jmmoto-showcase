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

import { getTelegramInitData, getTelegramUser } from '../utils/telegram';

export interface User {
  id: string;
  isAdmin: boolean;
  telegramId: number;
  telegramUsername: string;
  firstName: string;
  lastName: string;
  avatar: string;
}

export interface PatchMotorcycle {
  title?: string;
  price?: number;
  currency?: string;
  description?: string;
  status?: Motorcycle['status'];
}

export interface CreateUser {
  telegramId: number;
  telegramUsername: string;
  firstName: string;
  lastName: string;
  avatar: string;
}

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
    
    // Обрабатываем ответ - Huma возвращает { body: [...] }
    if (data && typeof data === 'object' && 'body' in data && Array.isArray(data.body)) {
      return data.body;
    }
    
    // Fallback: если данные приходят напрямую как массив
    if (Array.isArray(data)) {
      return data;
    }
    return [];
    
  } catch (error) {
    console.error('Error fetching motorcycles:', error);
    throw error;
  }
};

export const getMotorcycle = async (id: string): Promise<Motorcycle> => {
  try {
    const url = `${getApiBaseUrl()}/motorcycles/${id}`;
    
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

export const getCurrentUser = async (): Promise<User> => {
  try {
    const url = `${getApiBaseUrl()}/users/me`;
    
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
    console.error('Error fetching current user:', error);
    throw error;
  }
};

export const createUser = async (userData: CreateUser): Promise<User> => {
  try {
    const url = `${getApiBaseUrl()}/users/me`;
    
    const response = await fetch(url, {
      method: 'POST',
      headers: createApiHeaders(),
      body: JSON.stringify(userData),
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
    console.error('Error creating user:', error);
    throw error;
  }
};

export const getCurrentUserOrCreate = async (): Promise<User> => {
  try {
    // Сначала пытаемся получить существующего пользователя
    return await getCurrentUser();
  } catch (error) {
    // Если пользователь не найден (401), пытаемся создать нового
    if (error instanceof Error && error.message.includes('авторизация')) {
      // Получаем данные пользователя из Telegram
      const telegramUser = getTelegramUser();
      if (!telegramUser) {
        throw new Error('Данные пользователя Telegram недоступны');
      }

      // Создаем нового пользователя
      const newUserData: CreateUser = {
        telegramId: telegramUser.id,
        telegramUsername: telegramUser.username || '',
        firstName: telegramUser.first_name,
        lastName: telegramUser.last_name || '',
        avatar: telegramUser.photo_url || '',
      };

      // Создаем пользователя и возвращаем результат
      return await createUser(newUserData);
    }
    
    // Если это другая ошибка, пробрасываем её дальше
    throw error;
  }
};

export const updateMotorcycle = async (id: string, updates: PatchMotorcycle): Promise<Motorcycle> => {
  try {
    const url = `${getApiBaseUrl()}/admin/motorcycle/${id}`;
    
    const response = await fetch(url, {
      method: 'PATCH',
      headers: createApiHeaders(),
      body: JSON.stringify(updates),
    });
    
    if (!response.ok) {
      if (response.status === 401) {
        throw new Error('Требуется авторизация через Telegram');
      }
      if (response.status === 403) {
        throw new Error('Недостаточно прав для выполнения операции');
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
    console.error('Error updating motorcycle:', error);
    throw error;
  }
};

export const updateMotorcycleStatus = async (id: string, status: Motorcycle['status']): Promise<Motorcycle> => {
  try {
    const url = `${getApiBaseUrl()}/admin/motorcycle/${id}/status`;
    
    const response = await fetch(url, {
      method: 'PATCH',
      headers: createApiHeaders(),
      body: JSON.stringify({ status }),
    });
    
    if (!response.ok) {
      if (response.status === 401) {
        throw new Error('Требуется авторизация через Telegram');
      }
      if (response.status === 403) {
        throw new Error('Недостаточно прав для выполнения операции');
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
    console.error('Error updating motorcycle status:', error);
    throw error;
  }
};