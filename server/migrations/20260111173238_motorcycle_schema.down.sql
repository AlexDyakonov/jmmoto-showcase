-- Удаляем индексы
DROP INDEX IF EXISTS idx_user_is_admin;
DROP INDEX IF EXISTS idx_motorcycle_photo_order;
DROP INDEX IF EXISTS idx_motorcycle_photo_motorcycle_id;
DROP INDEX IF EXISTS idx_motorcycle_created_at;
DROP INDEX IF EXISTS idx_motorcycle_status;

-- Удаляем таблицы
DROP TABLE IF EXISTS "motorcycle_photo";
DROP TABLE IF EXISTS "motorcycle";

-- Удаляем поле is_admin из таблицы user
ALTER TABLE "user" DROP COLUMN IF EXISTS is_admin;

-- Восстанавливаем таблицу cat (если нужно)
CREATE TABLE IF NOT EXISTS "cat" (
    id VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    owner_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES "user"(id) ON DELETE CASCADE
);

CREATE INDEX idx_cat_owner_id ON cat(owner_id);
CREATE INDEX idx_cat_name ON cat(name);

