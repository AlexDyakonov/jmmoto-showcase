import { Motorcycle } from '../api/motorcycles';

interface MotorcycleCardProps {
  motorcycle: Motorcycle;
  onClick?: () => void;
}

export const MotorcycleCard = ({ motorcycle, onClick }: MotorcycleCardProps) => {
  const mainPhoto = motorcycle.photos?.[0]?.s3Url;
  const statusLabels = {
    available: 'В продаже',
    reserved: 'Бронь',
    sold: 'Продано',
    draft: 'Черновик',
  };

  const statusColors = {
    available: 'bg-green-500 text-white',
    reserved: 'bg-yellow-500 text-white',
    sold: 'bg-gray-500 text-white',
    draft: 'bg-blue-500 text-white',
  };

  // Форматируем цену правильно
  const formatPrice = (price: number, currency: string) => {
    if (price === 0 || !price) {
      return 'Цена не указана';
    }
    return `${price.toLocaleString('ru-RU')} ${currency}`;
  };

  return (
    <div
      onClick={onClick}
      className="bg-white rounded-lg overflow-hidden cursor-pointer transition-all hover:shadow-lg"
    >
      <div className="aspect-[4/3] relative bg-gray-100">
        {mainPhoto ? (
          <img
            src={mainPhoto}
            alt={motorcycle.title}
            className="w-full h-full object-cover"
            onError={(e) => {
              (e.target as HTMLImageElement).src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="400" height="300"%3E%3Crect fill="%23e5e7eb" width="400" height="300"/%3E%3Ctext x="50%25" y="50%25" text-anchor="middle" dy=".3em" fill="%239ca3af" font-family="sans-serif" font-size="16"%3EНет фото%3C/text%3E%3C/svg%3E';
            }}
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gray-100">
            Нет фото
          </div>
        )}
        <div className={`absolute top-2 left-2 px-2 py-1 rounded text-xs font-semibold ${statusColors[motorcycle.status]}`}>
          {statusLabels[motorcycle.status]}
        </div>
      </div>
      <div className="p-4">
        <h3 className="text-base font-semibold text-gray-900 mb-2 line-clamp-2 leading-tight">
          {motorcycle.title}
        </h3>
        <div className="flex items-center justify-between">
          <span className="text-lg font-bold text-gray-900">
            {formatPrice(motorcycle.price, motorcycle.currency)}
          </span>
          {motorcycle.photos && motorcycle.photos.length > 1 && (
            <span className="text-xs text-gray-500 font-medium">
              +{motorcycle.photos.length - 1}
            </span>
          )}
        </div>
      </div>
    </div>
  );
};

