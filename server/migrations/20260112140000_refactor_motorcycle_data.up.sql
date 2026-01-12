-- Добавляем новое поле data типа JSONB для структурированных данных
ALTER TABLE "motorcycle" ADD COLUMN IF NOT EXISTS data JSONB DEFAULT '{}';

-- Создаем функцию для парсинга description в структурированные данные
CREATE OR REPLACE FUNCTION parse_motorcycle_description(description_text TEXT)
RETURNS JSONB AS $$
DECLARE
    result JSONB := '{}';
    mileage_match TEXT;
    volume_match TEXT;
    frame_match TEXT;
BEGIN
    IF description_text IS NULL OR description_text = '' THEN
        RETURN result;
    END IF;
    
    -- Парсим пробег
    mileage_match := (regexp_matches(description_text, 'Пробег:\s*(\d+)\s*км', 'i'))[1];
    IF mileage_match IS NOT NULL THEN
        result := jsonb_set(result, '{mileage}', to_jsonb(mileage_match::integer));
        result := jsonb_set(result, '{mileage_unit}', '"км"');
    END IF;
    
    -- Парсим объем
    volume_match := (regexp_matches(description_text, 'Объем:\s*(\d+)\s*сс', 'i'))[1];
    IF volume_match IS NOT NULL THEN
        result := jsonb_set(result, '{volume}', to_jsonb(volume_match::integer));
        result := jsonb_set(result, '{volume_unit}', '"сс"');
    END IF;
    
    -- Парсим номер рамы
    frame_match := (regexp_matches(description_text, 'Номер рамы:\s*([A-Z0-9-]+)', 'i'))[1];
    IF frame_match IS NOT NULL THEN
        result := jsonb_set(result, '{frame_number}', to_jsonb(frame_match));
    END IF;
    
    -- Добавляем пустое поле для даты прибытия
    result := jsonb_set(result, '{arrival_date}', '""');
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Переносим данные из description в data
UPDATE "motorcycle" 
SET data = parse_motorcycle_description(description)
WHERE description IS NOT NULL AND description != '';

-- Удаляем функцию после использования
DROP FUNCTION parse_motorcycle_description(TEXT);

-- Удаляем поле description
ALTER TABLE "motorcycle" DROP COLUMN IF EXISTS description;

-- Создаем индекс для поиска по JSONB данным
CREATE INDEX IF NOT EXISTS idx_motorcycle_data_gin ON "motorcycle" USING GIN (data);
