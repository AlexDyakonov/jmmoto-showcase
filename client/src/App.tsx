import { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { MotorcycleList } from './pages/MotorcycleList';
import { MotorcycleDetail } from './pages/MotorcycleDetail';
import { initializeTelegramWebApp, getTelegramUser, isTelegramWebAppAvailable } from './utils/telegram';

function AppContent() {
  // Инициализация Telegram WebApp
  useEffect(() => {
    // Инициализируем Telegram WebApp
    initializeTelegramWebApp();
  }, []);

  return (
    <Routes>
      <Route path="/" element={<MotorcycleList />} />
      <Route path="/motorcycle/:id" element={<MotorcycleDetail />} />
    </Routes>
  );
}

function App() {
  return (
    <Router>
      <AppContent />
    </Router>
  );
}

export default App;