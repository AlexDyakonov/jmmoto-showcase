/**
 * Утилиты для работы с Telegram WebApp
 */

export interface TelegramUser {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  photo_url?: string;
}

export interface TelegramWebApp {
  initData: string;
  initDataUnsafe: {
    user?: TelegramUser;
  };
  ready: () => void;
  expand: () => void;
  close: () => void;
  setHeaderColor: (color: string) => void;
  setBackgroundColor: (color: string) => void;
}

declare global {
  interface Window {
    Telegram?: {
      WebApp: TelegramWebApp;
    };
  }
}

/**
 * Получает initData из Telegram WebApp
 * @returns строка initData или null если недоступно
 */
export const getTelegramInitData = (): string | null => {
  try {
    if (window.Telegram?.WebApp?.initData) {
      console.log('Telegram initData found:', window.Telegram.WebApp.initData);
      return window.Telegram.WebApp.initData;
    }
    
    console.log('Telegram initData not available');
    return null;
  } catch (error) {
    console.error('Error getting Telegram initData:', error);
    return null;
  }
};

/**
 * Получает информацию о пользователе из Telegram WebApp
 * @returns объект пользователя или null если недоступно
 */
export const getTelegramUser = (): TelegramUser | null => {
  try {
    if (window.Telegram?.WebApp?.initDataUnsafe?.user) {
      return window.Telegram.WebApp.initDataUnsafe.user;
    }
    return null;
  } catch (error) {
    console.error('Error getting Telegram user:', error);
    return null;
  }
};

/**
 * Проверяет, доступен ли Telegram WebApp
 * @returns true если Telegram WebApp доступен
 */
export const isTelegramWebAppAvailable = (): boolean => {
  return !!(window.Telegram?.WebApp);
};

/**
 * Инициализирует Telegram WebApp с настройками темы
 */
export const initializeTelegramWebApp = (): void => {
  if (window.Telegram?.WebApp) {
    const tg = window.Telegram.WebApp;
    
    // Настраиваем цвета
    tg.setHeaderColor('#0E0E0E');
    tg.setBackgroundColor('#0E0E0E');
    
    tg.ready();
    tg.expand();
    
    console.log('Telegram WebApp initialized');
  } else {
    console.log('Telegram WebApp not available');
  }
};
