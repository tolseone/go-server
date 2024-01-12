-- migrations/002_create_item_table.sql
CREATE TABLE IF NOT EXISTS public.item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    rarity VARCHAR(20) NOT NULL,
    description VARCHAR(100) NOT NULL
);
