import { useState } from 'react';
import { FilterMotorcycle } from '../api/motorcycles';

interface MotorcycleFilterProps {
  onFilterChange: (filters: FilterMotorcycle) => void;
}

export const MotorcycleFilter = ({ onFilterChange }: MotorcycleFilterProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const [filters, setFilters] = useState<FilterMotorcycle>({});

  const handleStatusChange = (status: 'available' | 'reserved' | 'sold' | undefined) => {
    const newFilters = { ...filters, status };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const handleTitleChange = (title: string) => {
    const newFilters = { ...filters, title: title || undefined };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const handlePriceChange = (field: 'minPrice' | 'maxPrice', value: string) => {
    const numValue = value ? parseFloat(value) : undefined;
    const newFilters = { ...filters, [field]: numValue };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const resetFilters = () => {
    const emptyFilters: FilterMotorcycle = {};
    setFilters(emptyFilters);
    onFilterChange(emptyFilters);
  };

  const hasActiveFilters = filters.status || filters.title || filters.minPrice || filters.maxPrice;

  return (
    <div className="mb-6">
      <div className="flex items-center gap-3 mb-3">
        <input
          type="text"
          value={filters.title || ''}
          onChange={(e) => handleTitleChange(e.target.value)}
          placeholder="Поиск по названию"
          className="flex-1 bg-white text-gray-900 px-4 py-2.5 rounded-lg border-0 focus:outline-none focus:ring-2 focus:ring-red-500 text-sm"
        />
        <button
          onClick={() => setIsOpen(!isOpen)}
          className={`px-4 py-2.5 rounded-lg text-sm font-medium transition-colors ${
            hasActiveFilters
              ? 'bg-red-500 text-white'
              : 'bg-white text-gray-900'
          }`}
        >
          {isOpen ? 'Скрыть' : 'Фильтры'}
        </button>
        {hasActiveFilters && (
          <button
            onClick={resetFilters}
            className="px-3 py-2.5 rounded-lg bg-white text-gray-600 hover:text-gray-900 text-sm"
          >
            ✕
          </button>
        )}
      </div>

      {isOpen && (
        <div className="bg-white rounded-lg p-4 space-y-3">
          <div>
            <div className="flex gap-2 mb-2">
              {(['available', 'reserved', 'sold'] as const).map((status) => (
                <button
                  key={status}
                  onClick={() => handleStatusChange(filters.status === status ? undefined : status)}
                  className={`px-3 py-1.5 text-xs font-medium rounded transition-colors ${
                    filters.status === status
                      ? 'bg-red-500 text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  {status === 'available' ? 'В продаже' : status === 'reserved' ? 'Бронь' : 'Продано'}
                </button>
              ))}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <input
                type="number"
                value={filters.minPrice || ''}
                onChange={(e) => handlePriceChange('minPrice', e.target.value)}
                placeholder="Цена от"
                className="w-full px-3 py-2 text-sm bg-gray-50 border border-gray-200 rounded focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent text-gray-900"
              />
            </div>
            <div>
              <input
                type="number"
                value={filters.maxPrice || ''}
                onChange={(e) => handlePriceChange('maxPrice', e.target.value)}
                placeholder="Цена до"
                className="w-full px-3 py-2 text-sm bg-gray-50 border border-gray-200 rounded focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent text-gray-900"
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

