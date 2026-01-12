import { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { MotorcycleList } from './pages/MotorcycleList';
import { MotorcycleDetail } from './pages/MotorcycleDetail';
import { initializeTelegramWebApp } from './utils/telegram';
import { analytics } from './utils/analytics';

function AppContent() {
  // Инициализация Telegram WebApp и аналитики
  useEffect(() => {
    // Инициализируем Telegram WebApp
    initializeTelegramWebApp();
    
    // Записываем заход пользователя
    analytics.recordVisit();
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