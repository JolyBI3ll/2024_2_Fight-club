-- 1. Создаем сервисного пользователя
CREATE ROLE service_user WITH LOGIN PASSWORD 'service_user_password';

-- 2. Ограничиваем возможности пользователя
ALTER ROLE service_user WITH NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;

-- 3. Предоставляем права на все таблицы в схеме public
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO service_user;

-- 4. Предоставляем доступ к последовательностям (если они используются)
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO service_user;

-- 5. Автоматически предоставляем права на новые объекты (таблицы, последовательности)
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO service_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO service_user;
