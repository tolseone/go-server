CREATE TABLE IF NOT EXISTS public.trade (
    trade_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    offered_items UUID[] NOT NULL,
    requested_items UUID[] NOT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp,
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE,
    CONSTRAINT valid_items CHECK (
        array_length(offered_items, 1) IS NOT NULL AND
        array_length(requested_items, 1) IS NOT NULL
    )
);
