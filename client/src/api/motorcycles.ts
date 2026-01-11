import client from './client';

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

export const getMotorcycles = async (filters?: FilterMotorcycle): Promise<Motorcycle[]> => {
  // Строим query параметры только если они есть
  const queryParams: Record<string, string> = {};
  
  if (filters?.status) {
    queryParams.status = filters.status;
  }
  if (filters?.title && filters.title.trim()) {
    queryParams.title = filters.title.trim();
  }
  if (filters?.minPrice !== undefined && filters.minPrice > 0) {
    queryParams.minPrice = filters.minPrice.toString();
  }
  if (filters?.maxPrice !== undefined && filters.maxPrice > 0) {
    queryParams.maxPrice = filters.maxPrice.toString();
  }

  try {
    // Если есть параметры, передаем их, иначе делаем простой GET запрос
    const response = Object.keys(queryParams).length > 0
      ? await client.get<{ body: Motorcycle[] } | Motorcycle[]>('/motorcycles', { params: queryParams })
      : await client.get<{ body: Motorcycle[] } | Motorcycle[]>('/motorcycles');
    
    // Обрабатываем ответ - Huma может возвращать { body: [...] } или просто [...]
    let data: Motorcycle[] = [];
    
    if (response.data) {
      if (Array.isArray(response.data)) {
        data = response.data;
      } else if (typeof response.data === 'object' && 'body' in response.data) {
        const body = (response.data as { body: Motorcycle[] }).body;
        data = Array.isArray(body) ? body : [];
      }
    }
    
    return data;
  } catch (error) {
    console.error('getMotorcycles error:', error);
    throw error;
  }
};

export const getMotorcycle = async (id: string): Promise<Motorcycle> => {
  try {
    const response = await client.get<{ body: Motorcycle } | Motorcycle>(`/motorcycles/${id}`);
    
    // Обрабатываем ответ
    if (response.data && typeof response.data === 'object' && 'body' in response.data) {
      return (response.data as { body: Motorcycle }).body;
    }
    
    return response.data as Motorcycle;
  } catch (error) {
    console.error('getMotorcycle error:', error);
    throw error;
  }
};
