-- Удаляем старую таблицу cat и связанные индексы
DROP INDEX IF EXISTS idx_cat_name;
DROP INDEX IF EXISTS idx_cat_owner_id;
DROP TABLE IF EXISTS "cat";

-- Добавляем поле is_admin в таблицу user
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS is_admin BOOLEAN DEFAULT FALSE NOT NULL;

-- Создаем таблицу motorcycle
CREATE TABLE IF NOT EXISTS "motorcycle" (
    id VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    price DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'RUB',
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'reserved', 'sold', 'draft')),
    source_url VARCHAR(500) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу motorcycle_photo
CREATE TABLE IF NOT EXISTS "motorcycle_photo" (
    id VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid(),
    motorcycle_id VARCHAR(255) NOT NULL,
    s3_url VARCHAR(500) NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (motorcycle_id) REFERENCES "motorcycle"(id) ON DELETE CASCADE
);

-- Создаем индексы
CREATE INDEX idx_motorcycle_status ON "motorcycle"(status);
CREATE INDEX idx_motorcycle_created_at ON "motorcycle"(created_at);
CREATE INDEX idx_motorcycle_photo_motorcycle_id ON "motorcycle_photo"(motorcycle_id);
CREATE INDEX idx_motorcycle_photo_order ON "motorcycle_photo"(motorcycle_id, "order");
CREATE INDEX idx_user_is_admin ON "user"(is_admin);

