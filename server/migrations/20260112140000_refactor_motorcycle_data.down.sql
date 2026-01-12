-- Удаляем индекс
DROP INDEX IF EXISTS idx_motorcycle_data_gin;

-- Добавляем обратно поле description
ALTER TABLE "motorcycle" ADD COLUMN IF NOT EXISTS description TEXT;

-- Создаем функцию для восстановления description из data
CREATE OR REPLACE FUNCTION restore_motorcycle_description(data_json JSONB)
RETURNS TEXT AS $$
DECLARE
    result TEXT := '';
    mileage INTEGER;
    volume INTEGER;
    frame_number TEXT;
BEGIN
    IF data_json IS NULL THEN
        RETURN NULL;
    END IF;
    
    -- Восстанавливаем пробег
    mileage := (data_json->>'mileage')::INTEGER;
    IF mileage IS NOT NULL THEN
        result := result || 'Пробег: ' || mileage || ' км. ';
    END IF;
    
    -- Восстанавливаем объем
    volume := (data_json->>'volume')::INTEGER;
    IF volume IS NOT NULL THEN
        result := result || 'Объем: ' || volume || ' сс. ';
    END IF;
    
    -- Восстанавливаем номер рамы
    frame_number := data_json->>'frame_number';
    IF frame_number IS NOT NULL AND frame_number != '' THEN
        result := result || 'Номер рамы: ' || frame_number || '. ';
    END IF;
    
    -- Убираем лишние пробелы в конце
    result := TRIM(result);
    
    RETURN CASE WHEN result = '' THEN NULL ELSE result END;
END;
$$ LANGUAGE plpgsql;

-- Восстанавливаем данные из data в description
UPDATE "motorcycle" 
SET description = restore_motorcycle_description(data)
WHERE data IS NOT NULL AND data != '{}';

-- Удаляем функцию после использования
DROP FUNCTION restore_motorcycle_description(JSONB);

-- Удаляем поле data
ALTER TABLE "motorcycle" DROP COLUMN IF EXISTS data;
