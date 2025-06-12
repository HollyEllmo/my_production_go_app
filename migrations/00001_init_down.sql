BEGIN;

-- Удаляем данные
DELETE FROM public.currency WHERE name = 'рубли' OR name = 'dollar';

-- Удаляем таблицы
DROP TABLE IF EXISTS public.product;
DROP TABLE IF EXISTS public.currency;
DROP TABLE IF EXISTS public.category;

-- Удаляем расширения (осторожно, может влиять на другие части системы)
-- DROP EXTENSION IF EXISTS pgcrypto;

COMMIT;