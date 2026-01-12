import { useState, useEffect } from 'react';
import { getCurrentUserOrCreate, User } from '../api/motorcycles';

interface UseCurrentUserResult {
  user: User | null;
  loading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
  isRegistering: boolean;
}

export const useCurrentUser = (): UseCurrentUserResult => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [isRegistering, setIsRegistering] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadUser = async () => {
    try {
      setLoading(true);
      setError(null);
      setIsRegistering(false);
      
      const userData = await getCurrentUserOrCreate();
      setUser(userData);
      setLoading(false);
    } catch (err) {
      console.error('Failed to load or create user:', err);
      
      // Не показываем ошибки авторизации пользователю
      if (err instanceof Error && err.message.includes('авторизация')) {
        // Показываем состояние регистрации
        setIsRegistering(true);
        setLoading(true);
        setUser(null);
      } else {
        // Для других ошибок показываем сообщение
        setError(err instanceof Error ? err.message : 'Не удалось загрузить пользователя');
        setUser(null);
        setLoading(false);
        setIsRegistering(false);
      }
    }
  };

  useEffect(() => {
    loadUser();
  }, []);

  const refetch = async () => {
    await loadUser();
  };

  return {
    user,
    loading,
    error,
    refetch,
    isRegistering,
  };
};
