import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getMotorcycle, updateMotorcycle, updateMotorcycleStatus, Motorcycle, PatchMotorcycle } from '../api/motorcycles';
import { isTelegramWebAppAvailable } from '../utils/telegram';
import { PhotoCarousel } from '../components/PhotoCarousel';
import { useCurrentUser } from '../hooks/useCurrentUser';

export const MotorcycleDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useCurrentUser();
  const [motorcycle, setMotorcycle] = useState<Motorcycle | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [editingField, setEditingField] = useState<string | null>(null);
  const [editValues, setEditValues] = useState<Partial<Motorcycle>>({});

  // Загружаем мотоцикл
  useEffect(() => {
    if (!id) {
      setError('ID мотоцикла не указан');
      setLoading(false);
      return;
    }

    const loadMotorcycle = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const data = await getMotorcycle(id);
        setMotorcycle(data);
        setEditValues(data);
      } catch (err) {
        console.error('Failed to load motorcycle:', err);
        setError(err instanceof Error ? err.message : 'Не удалось загрузить мотоцикл');
      } finally {
        setLoading(false);
      }
    };

    loadMotorcycle();
  }, [id]);

  const handleBack = () => {
    navigate('/');
  };

  const handleBuyClick = () => {
    if (isTelegramWebAppAvailable()) {
      // Открываем чат с администратором
      window.open('https://t.me/sharabarov', '_blank');
    }
  };

  const handleEditStart = (field: string) => {
    setEditingField(field);
  };

  const handleEditCancel = () => {
    setEditingField(null);
    if (motorcycle) {
      setEditValues(motorcycle);
    }
  };

  const handleEditSave = async (field: string) => {
    if (!motorcycle || !editValues || !user?.isAdmin || saving) return;

    setSaving(true);
    setError(null);
    
    try {
      // Подготавливаем данные для обновления
      const updates: PatchMotorcycle = {};
      
      if (field === 'title' && editValues.title !== undefined) {
        updates.title = editValues.title;
      } else if (field === 'price' && editValues.price !== undefined) {
        updates.price = editValues.price;
      } else if (field === 'arrival_date' && editValues.data?.arrival_date !== undefined) {
        updates.data = editValues.data;
      }

      // Отправляем запрос на сервер
      const updatedMotorcycle = await updateMotorcycle(motorcycle.id, updates);
      
      // Обновляем состояние
      setMotorcycle(updatedMotorcycle);
      setEditValues(updatedMotorcycle);
      setEditingField(null);
    } catch (err) {
      console.error('Failed to save changes:', err);
      setError(err instanceof Error ? err.message : 'Не удалось сохранить изменения');
      
      // Возвращаем исходные значения при ошибке
      if (motorcycle) {
        setEditValues(motorcycle);
      }
    } finally {
      setSaving(false);
    }
  };

  const handleStatusChange = async (newStatus: Motorcycle['status']) => {
    if (!motorcycle || !user?.isAdmin || saving) return;

    setSaving(true);
    setError(null);
    
    try {
      // Отправляем запрос на сервер
      const updatedMotorcycle = await updateMotorcycleStatus(motorcycle.id, newStatus);
      
      // Обновляем состояние
      setMotorcycle(updatedMotorcycle);
      setEditValues(updatedMotorcycle);
    } catch (err) {
      console.error('Failed to update status:', err);
      setError(err instanceof Error ? err.message : 'Не удалось обновить статус');
    } finally {
      setSaving(false);
    }
  };

  const getStatusColor = (status: Motorcycle['status']) => {
    switch (status) {
      case 'available': return 'bg-green-500';
      case 'reserved': return 'bg-yellow-500';
      case 'sold': return 'bg-red-500';
      case 'draft': return 'bg-gray-500';
      default: return 'bg-gray-500';
    }
  };

  const getStatusText = (status: Motorcycle['status']) => {
    switch (status) {
      case 'available': return 'В продаже';
      case 'reserved': return 'Забронирован';
      case 'sold': return 'Продан';
      case 'draft': return 'Черновик';
      default: return 'Неизвестно';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-[#0E0E0E] text-white flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-white mb-4"></div>
          <p className="text-gray-400">Загрузка мотоцикла...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-[#0E0E0E] text-white flex items-center justify-center">
        <div className="text-center">
          <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-6 max-w-md mx-auto">
            <p className="text-red-400 mb-4">{error}</p>
            <button
              onClick={handleBack}
              className="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
            >
              Назад к списку
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (!motorcycle) {
    return (
      <div className="min-h-screen bg-[#0E0E0E] text-white flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 mb-4">Мотоцикл не найден</p>
          <button
            onClick={handleBack}
            className="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
          >
            Назад к списку
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#0E0E0E] text-white">
      <div className="max-w-4xl mx-auto px-4 py-6">
        {/* Заголовок с кнопкой назад */}
        <div className="flex items-center justify-between mb-6">
          <button
            onClick={handleBack}
            className="flex items-center text-gray-400 hover:text-white transition-colors"
          >
            <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            Назад к списку
          </button>
          
          {/* Статус */}
          <div className="flex items-center">
            <span className={`px-3 py-1 rounded-full text-sm font-medium text-white ${getStatusColor(motorcycle.status)}`}>
              {getStatusText(motorcycle.status)}
            </span>
          </div>
        </div>

        {/* Карусель фотографий */}
        <div className="mb-8">
          <PhotoCarousel 
            photos={motorcycle.photos || []} 
            title={motorcycle.title}
          />
        </div>

        {/* Основная информация */}
        <div className="bg-gray-900/50 rounded-lg p-6 mb-6">
          {/* Название */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-400 mb-2">Название</label>
            {user?.isAdmin && editingField === 'title' ? (
              <div className="flex gap-2">
                <input
                  type="text"
                  value={editValues.title || ''}
                  onChange={(e) => setEditValues({ ...editValues, title: e.target.value })}
                  className="flex-1 bg-gray-800 border border-gray-700 rounded px-3 py-2 text-white"
                  autoFocus
                />
                <button
                  onClick={() => handleEditSave('title')}
                  disabled={saving}
                  className="px-3 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {saving ? (
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  ) : (
                    '✓'
                  )}
                </button>
                <button
                  onClick={handleEditCancel}
                  className="px-3 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
                >
                  ✕
                </button>
              </div>
            ) : (
              <div className="flex items-center justify-between">
                <h1 className="text-2xl font-bold text-white">{motorcycle.title}</h1>
                {user?.isAdmin && (
                  <button
                    onClick={() => handleEditStart('title')}
                    className="text-gray-400 hover:text-white transition-colors"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                    </svg>
                  </button>
                )}
              </div>
            )}
          </div>

          {/* Цена */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-400 mb-2">Цена</label>
            {user?.isAdmin && editingField === 'price' ? (
              <div className="flex gap-2">
                <input
                  type="number"
                  value={editValues.price || ''}
                  onChange={(e) => setEditValues({ ...editValues, price: Number(e.target.value) })}
                  className="flex-1 bg-gray-800 border border-gray-700 rounded px-3 py-2 text-white"
                  autoFocus
                />
                <button
                  onClick={() => handleEditSave('price')}
                  disabled={saving}
                  className="px-3 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {saving ? (
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  ) : (
                    '✓'
                  )}
                </button>
                <button
                  onClick={handleEditCancel}
                  className="px-3 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
                >
                  ✕
                </button>
              </div>
            ) : (
              <div className="flex items-center justify-between">
                <p className="text-xl font-semibold text-green-400">
                  {motorcycle.price.toLocaleString()} {motorcycle.currency}
                </p>
                {user?.isAdmin && (
                  <button
                    onClick={() => handleEditStart('price')}
                    className="text-gray-400 hover:text-white transition-colors"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                    </svg>
                  </button>
                )}
              </div>
            )}
          </div>

          {/* Технические характеристики */}
          {motorcycle.data && (
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-400 mb-3">Технические характеристики</label>
              <div className="space-y-2">
                {motorcycle.data.mileage && (
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">Пробег:</span>
                    <span className="text-white font-medium">
                      {motorcycle.data.mileage.toLocaleString()} {motorcycle.data.mileage_unit || 'км'}
                    </span>
                  </div>
                )}
                {motorcycle.data.volume && (
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">Объем двигателя:</span>
                    <span className="text-white font-medium">
                      {motorcycle.data.volume} {motorcycle.data.volume_unit || 'сс'}
                    </span>
                  </div>
                )}
                {motorcycle.data.frame_number && (
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">Номер рамы:</span>
                    <span className="text-white font-medium font-mono">{motorcycle.data.frame_number}</span>
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Дата прибытия */}
          {motorcycle.data?.arrival_date && (
            <div className="mb-4">
              {user?.isAdmin && editingField === 'arrival_date' ? (
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={editValues.data?.arrival_date || ''}
                    onChange={(e) => setEditValues({ 
                      ...editValues, 
                      data: { 
                        ...editValues.data, 
                        arrival_date: e.target.value 
                      } 
                    })}
                    className="flex-1 bg-gray-800 border border-gray-700 rounded px-3 py-2 text-white"
                    autoFocus
                    placeholder="Введите дату прибытия"
                  />
                  <button
                    onClick={() => handleEditSave('arrival_date')}
                    disabled={saving}
                    className="px-3 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {saving ? (
                      <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                    ) : (
                      '✓'
                    )}
                  </button>
                  <button
                    onClick={handleEditCancel}
                    className="px-3 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
                  >
                    ✕
                  </button>
                </div>
              ) : (
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <span className="text-gray-400 text-sm">Ориентировочная дата прибытия во Владивосток:</span>
                    <p className="text-white font-bold">{motorcycle.data.arrival_date}</p>
                  </div>
                  {user?.isAdmin && (
                    <button
                      onClick={() => handleEditStart('arrival_date')}
                      className="text-gray-400 hover:text-white transition-colors ml-2"
                    >
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                      </svg>
                    </button>
                  )}
                </div>
              )}
            </div>
          )}

          {/* Ссылка на источник */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-400 mb-2">Источник</label>
            <a
              href={motorcycle.sourceUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-400 hover:text-blue-300 transition-colors break-all"
            >
              {motorcycle.sourceUrl}
            </a>
          </div>
        </div>

        {/* Админские кнопки статуса */}
        {user?.isAdmin && (
          <div className="bg-gray-900/50 rounded-lg p-6 mb-6">
            <h3 className="text-lg font-semibold text-white mb-4">Управление статусом</h3>
            <div className="flex flex-wrap gap-2">
              {(['available', 'reserved', 'sold', 'draft'] as const).map((status) => (
                <button
                  key={status}
                  onClick={() => handleStatusChange(status)}
                  disabled={saving}
                  className={`px-4 py-2 rounded font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 ${
                    motorcycle.status === status
                      ? `${getStatusColor(status)} text-white`
                      : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                  }`}
                >
                  {saving && motorcycle.status !== status ? (
                    <div className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin"></div>
                  ) : null}
                  {getStatusText(status)}
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Кнопка покупки для обычных пользователей */}
        {!user?.isAdmin && motorcycle.status === 'available' && (
          <div className="bg-gray-900/50 rounded-lg p-6">
            <button
              onClick={handleBuyClick}
              className="w-full bg-green-600 hover:bg-green-700 text-white font-semibold py-3 px-6 rounded-lg transition-colors"
            >
              Купить мотоцикл
            </button>
            <p className="text-gray-400 text-sm text-center mt-2">
              Нажмите, чтобы связаться с продавцом
            </p>
          </div>
        )}

        {/* Информация о датах */}
        <div className="mt-6 text-sm text-gray-500 text-center">
          <p>Создано: {new Date(motorcycle.createdAt).toLocaleDateString('ru-RU')}</p>
          <p>Обновлено: {new Date(motorcycle.updatedAt).toLocaleDateString('ru-RU')}</p>
        </div>
      </div>
    </div>
  );
};
