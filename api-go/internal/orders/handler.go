package orders

import (
    "errors"
    "net/http"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"
)
type Handler struct {
	repo OrderRepository
}

func NewHandler(repo OrderRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("pageindex", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "page must be an integer greater than or equal to 1",
		})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pagesize", "20"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "page_size must be an integer between 1 and 100",
		})
		return
	}

	queryPtr := func(key string) *string {
		value, exists := c.GetQuery(key)
		if !exists {
			return nil
		}
		return &value
	}

	filter := OrderFilter{
		CuisineCategory: queryPtr("cuisine_category"),
		OrderID:         queryPtr("order_id"),
		MetroCity:       queryPtr("metro_city"),
		OrderChannel:    queryPtr("order_channel"),
	}

	orders, err := h.repo.GetList(c.Request.Context(), page, pageSize, filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not load orders",
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetByID(c *gin.Context) {
	orderID := strings.TrimSpace(c.Param("id"))
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	order, err := h.repo.GetByID(c.Request.Context(), orderID)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load order"})
		return
	}

	c.JSON(http.StatusOK, order)
}