// Простая система аналитики для отслеживания заходов
import { getTelegramInitData } from './telegram';

class SimpleAnalytics {
  private hasRecordedVisit = false;

  public async recordVisit(source?: string): Promise<void> {
    // Записываем только один раз за сессию
    if (this.hasRecordedVisit) {
      return;
    }

    try {
      const initData = getTelegramInitData();
      if (!initData) {
        console.warn('No Telegram initData available for analytics');
        return;
      }

      const apiUrl = `${window.api?.API_URL || 'http://localhost:8000'}/api/v1`;
      
      await fetch(`${apiUrl}/analytics/visit`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-API-Token': initData
        },
        body: JSON.stringify({
          source: source || this.getSourceFromUrl()
        })
      });

      this.hasRecordedVisit = true;
      console.log('Visit recorded successfully', { source });
    } catch (error) {
      console.warn('Failed to record visit:', error);
    }
  }

  private getSourceFromUrl(): string {
    // Пытаемся определить источник из URL или Telegram данных
    const urlParams = new URLSearchParams(window.location.search);
    const startParam = urlParams.get('startapp');
    
    // Проверяем Telegram WebApp данные
    if (window.Telegram?.WebApp?.initDataUnsafe?.start_param) {
      return window.Telegram.WebApp.initDataUnsafe.start_param;
    }

    if (startParam) {
      return startParam;
    }

    // Определяем по URL
    if (window.location.href.includes('t.me')) {
      return 'telegram_webapp';
    }

    return 'direct';
  }

  public async getUserStats(): Promise<any> {
    try {
      const initData = getTelegramInitData();
      if (!initData) {
        return null;
      }

      const apiUrl = `${window.api?.API_URL || 'http://localhost:8000'}/api/v1`;
      
      const response = await fetch(`${apiUrl}/analytics/my-stats`, {
        method: 'GET',
        headers: {
          'X-API-Token': initData
        }
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }

      const data = await response.json();
      return data.body;
    } catch (error) {
      console.warn('Failed to get user stats:', error);
      return null;
    }
  }
}

// Создаем глобальный экземпляр
export const analytics = new SimpleAnalytics();

// Добавляем в глобальную область для отладки
if (typeof window !== 'undefined') {
  (window as any).analytics = analytics;
}
