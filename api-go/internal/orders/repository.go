package orders

import (
    "context"
	"encoding/json"
	"strings"
	"fmt"
    "github.com/jackc/pgx/v5/pgxpool"
)
type OrderRepository interface {
	GetByID(ctx context.Context, orderID string) (*OrderDetail, error)
	GetList(ctx context.Context, pageIndex, pageSize int, filter OrderFilter) ([]Order, error)
}

type PostgresOrderRepo struct {
	db *pgxpool.Pool
}

type OrderFilter struct {
	CuisineCategory *string
	OrderID         *string
	MetroCity       *string
	OrderChannel    *string
}

type Order struct {
	ID             string    `json:"id"`
	OrderID        *string  `json:"order_id"`
	MetroCity      string   `json:"metro_city"`
	CuisineCategory string   `json:"cuisine_category"`
	OrderChannel   string   `json:"order_channel"`
	OrderValueUSD  *json.Number   `json:"order_value_usd"`
	CustomerRating *json.Number `json:"customer_rating"`
}


type OrderDetail struct {
	ID                              int64        `json:"id"`
	OrderID                         *string      `json:"order_id"`
	MetroCity                       string       `json:"metro_city"`
	CuisineCategory                 string       `json:"cuisine_category"`
	OrderChannel                    string       `json:"order_channel"`
	OrderValueUSD                   json.Number  `json:"order_value_usd"`
	DiscountAppliedUSD              *json.Number `json:"discount_applied_usd"`
	DeliveryDistanceKM              *json.Number `json:"delivery_distance_km"`
	EstimatedPrepTimeMin            *int32       `json:"estimated_prep_time_min"`
	ActualDeliveryTimeMin           *json.Number `json:"actual_delivery_time_min"`
	DeliveryFeeUSD                  *json.Number `json:"delivery_fee_usd"`
	CustomerRating                  *json.Number `json:"customer_rating"`
	RestaurantPriceTier             *string      `json:"restaurant_price_tier"`
	IsPeakHour                      *bool        `json:"is_peak_hour"`
	WeatherCondition                *string      `json:"weather_condition"`
	NetEffectiveOrderPrice          *json.Number `json:"net_effective_order_price"`
	DiscountPercentage              *json.Number `json:"discount_percentage"`
	DeliverySpeedKMPerHr            *json.Number `json:"delivery_speed_km_per_hr"`
	TotalFulfillmentTimeMin         *json.Number `json:"total_fulfillment_time_min"`
	DeliveryFeeRatio                *json.Number `json:"delivery_fee_ratio"`
	CustomerSpendingTier            *string      `json:"customer_spending_tier"`
	LogisticsFrictionIndex          *json.Number `json:"logistics_friction_index"`
	EstimatedRestaurantNetMarginUSD *json.Number `json:"estimated_restaurant_net_margin_usd"`
	CustomerSatisfactionTier        *string      `json:"customer_satisfaction_tier"`
	PrimarySatisfactionDriver       *string      `json:"primary_satisfaction_driver"`
	COTMarketDynamicsRationale      *string      `json:"cot_market_dynamics_rationale"`
	ConfidenceScore                 *json.Number `json:"confidence_score"`
	Status                          string       `json:"status"`
}

func NewPostgresOrderRepo(db *pgxpool.Pool) *PostgresOrderRepo {
	return &PostgresOrderRepo{db: db}
}

func (r *PostgresOrderRepo) GetByID(ctx context.Context, orderID string) (*OrderDetail, error) {
	const query = `
		SELECT
			id,
			order_id,
			metro_city,
			cuisine_category,
			order_channel,
			order_value_usd::text,
			discount_applied_usd::text,
			delivery_distance_km::text,
			estimated_prep_time_min,
			actual_delivery_time_min::text,
			delivery_fee_usd::text,
			customer_rating::text,
			restaurant_price_tier,
			is_peak_hour,
			weather_condition,
			net_effective_order_price::text,
			discount_percentage::text,
			delivery_speed_km_per_hr::text,
			total_fulfillment_time_min::text,
			delivery_fee_ratio::text,
			customer_spending_tier,
			logistics_friction_index::text,
			estimated_restaurant_net_margin_usd::text,
			customer_satisfaction_tier,
			primary_satisfaction_driver,
			cot_market_dynamics_rationale,
			confidence_score::text,
			status
		FROM orders
		WHERE order_id = $1
	`

	var order OrderDetail
	var orderValueUSD string
	var (
		discountAppliedUSD              *string
		deliveryDistanceKM              *string
		actualDeliveryTimeMin           *string
		deliveryFeeUSD                  *string
		customerRating                  *string
		netEffectiveOrderPrice          *string
		discountPercentage              *string
		deliverySpeedKMPerHr            *string
		totalFulfillmentTimeMin         *string
		deliveryFeeRatio                *string
		logisticsFrictionIndex          *string
		estimatedRestaurantNetMarginUSD *string
		confidenceScore                 *string
	)

	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.OrderID,
		&order.MetroCity,
		&order.CuisineCategory,
		&order.OrderChannel,
		&orderValueUSD,
		&discountAppliedUSD,
		&deliveryDistanceKM,
		&order.EstimatedPrepTimeMin,
		&actualDeliveryTimeMin,
		&deliveryFeeUSD,
		&customerRating,
		&order.RestaurantPriceTier,
		&order.IsPeakHour,
		&order.WeatherCondition,
		&netEffectiveOrderPrice,
		&discountPercentage,
		&deliverySpeedKMPerHr,
		&totalFulfillmentTimeMin,
		&deliveryFeeRatio,
		&order.CustomerSpendingTier,
		&logisticsFrictionIndex,
		&estimatedRestaurantNetMarginUSD,
		&order.CustomerSatisfactionTier,
		&order.PrimarySatisfactionDriver,
		&order.COTMarketDynamicsRationale,
		&confidenceScore,
		&order.Status,
	)
	if err != nil {
		return nil, err
	}

	order.OrderValueUSD = json.Number(orderValueUSD)
	order.DiscountAppliedUSD = nullableJSONNumber(discountAppliedUSD)
	order.DeliveryDistanceKM = nullableJSONNumber(deliveryDistanceKM)
	order.ActualDeliveryTimeMin = nullableJSONNumber(actualDeliveryTimeMin)
	order.DeliveryFeeUSD = nullableJSONNumber(deliveryFeeUSD)
	order.CustomerRating = nullableJSONNumber(customerRating)
	order.NetEffectiveOrderPrice = nullableJSONNumber(netEffectiveOrderPrice)
	order.DiscountPercentage = nullableJSONNumber(discountPercentage)
	order.DeliverySpeedKMPerHr = nullableJSONNumber(deliverySpeedKMPerHr)
	order.TotalFulfillmentTimeMin = nullableJSONNumber(totalFulfillmentTimeMin)
	order.DeliveryFeeRatio = nullableJSONNumber(deliveryFeeRatio)
	order.LogisticsFrictionIndex = nullableJSONNumber(logisticsFrictionIndex)
	order.EstimatedRestaurantNetMarginUSD = nullableJSONNumber(estimatedRestaurantNetMarginUSD)
	order.ConfidenceScore = nullableJSONNumber(confidenceScore)

	return &order, nil
}

func nullableJSONNumber(value *string) *json.Number {
	if value == nil {
		return nil
	}

	number := json.Number(*value)
	return &number
}

func (r *PostgresOrderRepo) GetList(
	ctx context.Context,
	pageIndex, pageSize int,
	filter OrderFilter,
) ([]Order, error) {
	const baseQuery = `
		SELECT
			id,
			order_id,
			metro_city,
			cuisine_category,
			order_channel,
			order_value_usd::text,
			customer_rating::double precision
		FROM orders
	`

	conditions := make([]string, 0, 4)
	args := make([]any, 0, 6)

	filters := []struct {
		column string
		value  *string
	}{
		{"cuisine_category", filter.CuisineCategory},
		{"order_id", filter.OrderID},
		{"metro_city", filter.MetroCity},
		{"order_channel", filter.OrderChannel},
	}

	for _, field := range filters {
		if field.value == nil {
			continue
		}

		value := strings.TrimSpace(*field.value)
		if value == "" {
			continue
		}

		args = append(args, value)
		conditions = append(conditions, fmt.Sprintf("%s = $%d", field.column, len(args)))
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += "\nWHERE " + strings.Join(conditions, " AND ")
	}

	offset := (pageIndex - 1) * pageSize
	query += fmt.Sprintf("\nORDER BY id ASC\nLIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]Order, 0)
	for rows.Next() {
		var order Order
		if err := rows.Scan(
			&order.ID,
			&order.OrderID,
			&order.MetroCity,
			&order.CuisineCategory,
			&order.OrderChannel,
			&order.OrderValueUSD,
			&order.CustomerRating,
		); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}