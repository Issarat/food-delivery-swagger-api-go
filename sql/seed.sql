BEGIN;
CREATE TEMP TABLE stage_orders ON COMMIT DROP AS
SELECT order_id, metro_city, cuisine_category, order_channel, order_value_usd, discount_applied_usd, delivery_distance_km, estimated_prep_time_min, actual_delivery_time_min, delivery_fee_usd, customer_rating, restaurant_price_tier, is_peak_hour, weather_condition, net_effective_order_price, discount_percentage, delivery_speed_km_per_hr, total_fulfillment_time_min, delivery_fee_ratio, customer_spending_tier, logistics_friction_index, estimated_restaurant_net_margin_usd, customer_satisfaction_tier, primary_satisfaction_driver, cot_market_dynamics_rationale, confidence_score
FROM orders WITH NO DATA;

\copy stage_orders FROM '/data/global_food_delivery_enriched_market_2026.csv' WITH (FORMAT csv, HEADER MATCH, ENCODING 'UTF8')

INSERT INTO orders (order_id, metro_city, cuisine_category, order_channel, order_value_usd, discount_applied_usd, delivery_distance_km, estimated_prep_time_min, actual_delivery_time_min, delivery_fee_usd, customer_rating, restaurant_price_tier, is_peak_hour, weather_condition, net_effective_order_price, discount_percentage, delivery_speed_km_per_hr, total_fulfillment_time_min, delivery_fee_ratio, customer_spending_tier, logistics_friction_index, estimated_restaurant_net_margin_usd, customer_satisfaction_tier, primary_satisfaction_driver, cot_market_dynamics_rationale, confidence_score, status)
SELECT order_id, metro_city, cuisine_category, order_channel, order_value_usd, discount_applied_usd, delivery_distance_km, estimated_prep_time_min, actual_delivery_time_min, delivery_fee_usd, customer_rating, restaurant_price_tier, is_peak_hour, weather_condition, net_effective_order_price, discount_percentage, delivery_speed_km_per_hr, total_fulfillment_time_min, delivery_fee_ratio, customer_spending_tier, logistics_friction_index, estimated_restaurant_net_margin_usd, customer_satisfaction_tier, primary_satisfaction_driver, cot_market_dynamics_rationale, confidence_score, 'delivered'
FROM stage_orders
ON CONFLICT (order_id) DO NOTHING;
COMMIT;
