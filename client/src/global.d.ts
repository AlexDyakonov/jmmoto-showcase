interface Window {
  api: {
    BOT_USERNAME: string;
    API_URL: string;
  };
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
