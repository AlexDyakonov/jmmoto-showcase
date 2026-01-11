import { useEffect, useState } from 'react';
import { getMotorcycles, Motorcycle, FilterMotorcycle } from './api/motorcycles';
import { MotorcycleCard } from './components/MotorcycleCard';
import { MotorcycleFilter } from './components/MotorcycleFilter';

// Расширяем типы для window.Telegram
declare global {
  interface Window {
    Telegram?: {
      WebApp: {
        initData: string;
        initDataUnsafe: {
          user?: {
            id: number;
            first_name: string;
            last_name?: string;
            username?: string;
            photo_url?: string;
          };
        };
        ready: () => void;
        expand: () => void;
        close: () => void;
        setHeaderColor: (color: string) => void;
        setBackgroundColor: (color: string) => void;
      };
    };
  }
}

function App() {
  const [motorcycles, setMotorcycles] = useState<Motorcycle[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<FilterMotorcycle>({});

  // Инициализация Telegram WebApp
  useEffect(() => {
    console.log('Initializing Telegram WebApp...');
    
    if (window.Telegram?.WebApp) {
      const tg = window.Telegram.WebApp;
      
      // Настраиваем цвета обертки
      tg.setHeaderColor('#0E0E0E');
      tg.setBackgroundColor('#0E0E0E');
      
      tg.ready();
      tg.expand();
      
      console.log('Telegram WebApp initialized');
    } else {
      console.log('Telegram WebApp not available');
    }
  }, []);

  // Загрузка мотоциклов
  useEffect(() => {
    let isCancelled = false;
    
    const loadMotorcycles = async () => {
      console.log('Loading motorcycles with filters:', filters);
      
      try {
        setLoading(true);
        setError(null);
        
        const data = await getMotorcycles(filters);
        
        if (!isCancelled) {
          console.log('Loaded motorcycles:', data.length);
          setMotorcycles(data || []);
        }
      } catch (err) {
        if (!isCancelled) {
          console.error('Failed to load motorcycles:', err);
          setError(err instanceof Error ? err.message : 'Не удалось загрузить мотоциклы');
          setMotorcycles([]);
        }
      } finally {
        if (!isCancelled) {
          setLoading(false);
        }
      }
    };

    loadMotorcycles();
    
    // Cleanup function для отмены запроса
    return () => {
      isCancelled = true;
    };
  }, [filters]);

  const handleFilterChange = (newFilters: FilterMotorcycle) => {
    console.log('Filter changed:', newFilters);
    setFilters(newFilters);
  };

  return (
    <div className="min-h-screen bg-[#0E0E0E] text-white">
      <div className="max-w-6xl mx-auto px-4 py-6">
        {/* Фильтр */}
        <MotorcycleFilter onFilterChange={handleFilterChange} />

        {/* Контент */}
        {loading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
            <p className="mt-4 text-gray-400">Загрузка мотоциклов...</p>
          </div>
        ) : error ? (
          <div className="text-center py-12">
            <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-6 max-w-md mx-auto">
              <p className="text-red-400 mb-4">{error}</p>
              <button
                onClick={() => window.location.reload()}
                className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 transition-colors"
              >
                Обновить страницу
              </button>
            </div>
          </div>
        ) : motorcycles.length === 0 ? (
          <div className="text-center py-12">
            <div className="bg-gray-800/50 rounded-lg p-8 max-w-md mx-auto">
              <p className="text-gray-400 text-lg">Мотоциклы не найдены</p>
              <p className="text-gray-500 text-sm mt-2">
                Попробуйте изменить фильтры поиска
              </p>
            </div>
          </div>
        ) : (
          <>
            <div className="mb-6 text-sm text-gray-400">
              Найдено мотоциклов: <span className="text-white font-medium">{motorcycles.length}</span>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {motorcycles.map((motorcycle) => (
                <MotorcycleCard 
                  key={motorcycle.id} 
                  motorcycle={motorcycle}
                  onClick={() => {
                    console.log('Motorcycle clicked:', motorcycle.id);
                    // Здесь можно добавить навигацию к детальной странице
                  }}
                />
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default App;