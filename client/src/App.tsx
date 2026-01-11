import { useEffect, useState } from 'react';
import { setTelegramToken } from './api/client';
import { getMotorcycles, Motorcycle, FilterMotorcycle } from './api/motorcycles';
import { MotorcycleCard } from './components/MotorcycleCard';
import { MotorcycleFilter } from './components/MotorcycleFilter';

function MotorcycleShowcase() {
  const [motorcycles, setMotorcycles] = useState<Motorcycle[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<FilterMotorcycle>({});
  const [isInitialized, setIsInitialized] = useState(false);

  // Инициализация Telegram WebApp - выполняется один раз при монтировании
  useEffect(() => {
    const initTelegram = () => {
      if (window.Telegram?.WebApp) {
        const tg = window.Telegram.WebApp;
        
        // Настраиваем цвета обертки
        tg.setHeaderColor('#0E0E0E');
        tg.setBackgroundColor('#0E0E0E');
        
        tg.ready();
        tg.expand();
        
        // Устанавливаем токен Telegram для API запросов (если доступен)
        const initData = tg.initData;
        if (initData) {
          setTelegramToken(initData);
        }
      }
      
      // Проверяем, что API_URL доступен
      if (window.api?.API_URL) {
        setIsInitialized(true);
      } else {
        // Если API_URL еще не загружен, ждем немного и проверяем снова
        const checkInterval = setInterval(() => {
          if (window.api?.API_URL) {
            setIsInitialized(true);
            clearInterval(checkInterval);
          }
        }, 100);
        
        // Таймаут на случай, если config.js не загрузится
        setTimeout(() => {
          clearInterval(checkInterval);
          if (!window.api?.API_URL) {
            console.warn('API_URL not found, using default');
            setIsInitialized(true);
          }
        }, 2000);
      }
    };

    initTelegram();
  }, []);

  // Загрузка мотоциклов - выполняется после инициализации и при изменении фильтров
  useEffect(() => {
    if (!isInitialized) {
      return;
    }

    let cancelled = false;

    const loadMotorcycles = async () => {
      try {
        setLoading(true);
        setError(null);
        
        // Вызываем API без фильтров при первой загрузке
        const data = await getMotorcycles(filters);
        
        if (cancelled) {
          return;
        }
        
        if (Array.isArray(data)) {
          setMotorcycles(data);
        } else {
          console.error('Invalid data format:', data);
          setMotorcycles([]);
          setError('Неверный формат данных');
        }
      } catch (err) {
        if (cancelled) {
          return;
        }
        
        console.error('Failed to load motorcycles:', err);
        setError('Не удалось загрузить мотоциклы');
        setMotorcycles([]);
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    loadMotorcycles();

    return () => {
      cancelled = true;
    };
  }, [isInitialized, filters]);

  return (
    <div className="min-h-screen bg-[#0E0E0E] text-white">
      <div className="max-w-6xl mx-auto px-4 py-4">
        {/* Компактный фильтр */}
        <MotorcycleFilter onFilterChange={setFilters} />

        {/* Список мотоциклов */}
        {loading ? (
          <div className="text-center py-16">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-2 border-white border-t-transparent"></div>
            <p className="mt-4 text-sm text-gray-400">Загрузка...</p>
          </div>
        ) : error ? (
          <div className="text-center py-16">
            <p className="text-sm text-red-400 mb-4">{error}</p>
            <button
              onClick={() => {
                setError(null);
                setLoading(true);
                getMotorcycles(filters)
                  .then((data) => {
                    if (Array.isArray(data)) {
                      setMotorcycles(data);
                    } else {
                      setMotorcycles([]);
                    }
                  })
                  .catch((err) => {
                    console.error('Retry failed:', err);
                    setError('Не удалось загрузить мотоциклы');
                  })
                  .finally(() => {
                    setLoading(false);
                  });
              }}
              className="px-6 py-2.5 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-colors"
            >
              Повторить
            </button>
          </div>
        ) : !motorcycles || motorcycles.length === 0 ? (
          <div className="text-center py-16">
            <p className="text-sm text-gray-400">Мотоциклы не найдены</p>
          </div>
        ) : (
          <>
            <div className="mb-3 text-xs text-gray-400 uppercase tracking-wider">
              Найдено: {motorcycles.length}
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {motorcycles.map((motorcycle) => (
                <MotorcycleCard key={motorcycle.id} motorcycle={motorcycle} />
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
}

function App() {
  return <MotorcycleShowcase />;
}

export default App;
