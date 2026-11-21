BEGIN;
CREATE TABLE IF NOT EXISTS users (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON users (lower(email));

CREATE TABLE IF NOT EXISTS revoked_tokens (
    jti        TEXT PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS revoked_tokens_expires_at_idx
    ON revoked_tokens (expires_at);

CREATE TABLE IF NOT EXISTS orders (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id text UNIQUE,
    metro_city text NOT NULL,
    cuisine_category text NOT NULL,
    order_channel text NOT NULL,
    order_value_usd numeric NOT NULL,
    discount_applied_usd numeric,
    delivery_distance_km numeric,
    estimated_prep_time_min integer,
    actual_delivery_time_min numeric,
    delivery_fee_usd numeric,
    customer_rating numeric,
    restaurant_price_tier text,
    is_peak_hour boolean,
    weather_condition text,
    net_effective_order_price numeric,
    discount_percentage numeric,
    delivery_speed_km_per_hr numeric,
    total_fulfillment_time_min numeric,
    delivery_fee_ratio numeric,
    customer_spending_tier text,
    logistics_friction_index numeric,
    estimated_restaurant_net_margin_usd numeric,
    customer_satisfaction_tier text,
    primary_satisfaction_driver text,
    cot_market_dynamics_rationale text,
    confidence_score numeric,
    created_by bigint REFERENCES users(id),
    status text NOT NULL DEFAULT 'pending',
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS orders_created_by_idx ON orders (created_by);
COMMIT;
