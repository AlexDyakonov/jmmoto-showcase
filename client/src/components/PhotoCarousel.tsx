import { useState } from 'react';

interface Photo {
  id: string;
  s3Url: string;
  order: number;
}

interface PhotoCarouselProps {
  photos: Photo[];
  title: string;
}

export const PhotoCarousel: React.FC<PhotoCarouselProps> = ({ photos, title }) => {
  const [currentIndex, setCurrentIndex] = useState(0);

  if (!photos || photos.length === 0) {
    return (
      <div className="w-full aspect-video bg-gray-800 rounded-lg flex items-center justify-center">
        <div className="text-center text-gray-400">
          <svg className="w-16 h-16 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <p>Фотографии недоступны</p>
        </div>
      </div>
    );
  }

  // Сортируем фотографии по порядку
  const sortedPhotos = [...photos].sort((a, b) => a.order - b.order);

  const goToPrevious = () => {
    setCurrentIndex((prevIndex) => 
      prevIndex === 0 ? sortedPhotos.length - 1 : prevIndex - 1
    );
  };

  const goToNext = () => {
    setCurrentIndex((prevIndex) => 
      prevIndex === sortedPhotos.length - 1 ? 0 : prevIndex + 1
    );
  };

  const goToSlide = (index: number) => {
    setCurrentIndex(index);
  };

  return (
    <div className="relative w-full">
      {/* Основное изображение */}
      <div className="relative aspect-video bg-gray-800 rounded-lg overflow-hidden group">
        <img
          src={sortedPhotos[currentIndex].s3Url}
          alt={`${title} - фото ${currentIndex + 1}`}
          className="w-full h-full object-cover transition-opacity duration-300"
          loading="lazy"
        />
        
        {/* Градиент для лучшей видимости кнопок */}
        <div className="absolute inset-0 bg-gradient-to-r from-black/20 via-transparent to-black/20 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
        
        {/* Кнопки навигации */}
        {sortedPhotos.length > 1 && (
          <>
            {/* Предыдущая */}
            <button
              onClick={goToPrevious}
              className="absolute left-4 top-1/2 -translate-y-1/2 bg-black/50 hover:bg-black/70 text-white p-2 rounded-full opacity-0 group-hover:opacity-100 transition-all duration-300 backdrop-blur-sm"
              aria-label="Предыдущее фото"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            
            {/* Следующая */}
            <button
              onClick={goToNext}
              className="absolute right-4 top-1/2 -translate-y-1/2 bg-black/50 hover:bg-black/70 text-white p-2 rounded-full opacity-0 group-hover:opacity-100 transition-all duration-300 backdrop-blur-sm"
              aria-label="Следующее фото"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </>
        )}
        
        {/* Индикатор текущего фото */}
        {sortedPhotos.length > 1 && (
          <div className="absolute bottom-4 left-1/2 -translate-x-1/2 bg-black/50 backdrop-blur-sm rounded-full px-3 py-1 text-white text-sm">
            {currentIndex + 1} / {sortedPhotos.length}
          </div>
        )}
      </div>
      
      {/* Миниатюры */}
      {sortedPhotos.length > 1 && (
        <div className="mt-4 flex gap-2 overflow-x-auto pb-2">
          {sortedPhotos.map((photo, index) => (
            <button
              key={photo.id}
              onClick={() => goToSlide(index)}
              className={`flex-shrink-0 w-20 h-16 rounded-lg overflow-hidden border-2 transition-all duration-200 ${
                index === currentIndex
                  ? 'border-blue-500 ring-2 ring-blue-500/30'
                  : 'border-gray-600 hover:border-gray-400'
              }`}
            >
              <img
                src={photo.s3Url}
                alt={`${title} - миниатюра ${index + 1}`}
                className="w-full h-full object-cover"
                loading="lazy"
              />
            </button>
          ))}
        </div>
      )}
      
      {/* Поддержка свайпов на мобильных устройствах */}
      {sortedPhotos.length > 1 && (
        <div className="mt-2 text-center text-xs text-gray-500">
          Используйте стрелки или нажмите на миниатюры для навигации
        </div>
      )}
    </div>
  );
};
