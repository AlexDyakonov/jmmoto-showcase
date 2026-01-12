-- Простая таблица для отслеживания заходов пользователей
CREATE TABLE user_visits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    source VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для быстрых запросов по пользователю и времени
CREATE INDEX idx_user_visits_user_id ON user_visits(user_id);
CREATE INDEX idx_user_visits_created_at ON user_visits(created_at);
CREATE INDEX idx_user_visits_user_time ON user_visits(user_id, created_at);

-- Простое представление для подсчета частоты заходов
CREATE VIEW user_visit_stats AS
SELECT 
    user_id,
    COUNT(*) as total_visits,
    COUNT(DISTINCT DATE(created_at)) as unique_days,
    MIN(created_at) as first_visit,
    MAX(created_at) as last_visit,
    EXTRACT(DAYS FROM (MAX(created_at) - MIN(created_at))) + 1 as days_span,
    ROUND(COUNT(*)::numeric / GREATEST(EXTRACT(DAYS FROM (MAX(created_at) - MIN(created_at))) + 1, 1), 2) as avg_visits_per_day
FROM user_visits 
GROUP BY user_id;
