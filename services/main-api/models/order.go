package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPendingPayment OrderStatus = "PendingPayment"
	StatusInProduction   OrderStatus = "InProduction"
	StatusQualityControl OrderStatus = "QualityControl"
	StatusShipped        OrderStatus = "Shipped"
	StatusDelivered      OrderStatus = "Delivered"
	StatusCancelled      OrderStatus = "Cancelled"
)

// Order represents a customer's confirmed B2B order
type Order struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CustomerID        uint           `json:"customer_id"`
	TotalAmountUSD    float64        `json:"total_amount_usd"`
	FinalSellingPrice float64        `json:"final_selling_price"`
	Currency          string         `json:"currency"`
	Status            OrderStatus    `json:"status"`
	PaymentID         string         `json:"payment_id"` // Reference to Stripe/Iyzico
	Items             []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`
	ShippingAddress   string         `json:"shipping_address"`
	EstimatedDelivery time.Time      `json:"estimated_delivery"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// OrderItem represents a single line item (e.g., a specific DXF part) in an order
type OrderItem struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrderID         uint           `json:"order_id"`
	PartName        string         `json:"part_name"` // E.g., file name uploaded
	Material        string         `json:"material"`
	ThicknessMm     float64        `json:"thickness_mm"`
	Quantity        int            `json:"quantity"`
	UnitCostUSD     float64        `json:"unit_cost_usd"`
	TotalCostUSD    float64        `json:"total_cost_usd"`
	LaserCutting    bool           `json:"laser_cutting"`
	BendingCount    int            `json:"bending_count"`
	ComplexBending  bool           `json:"complex_bending"`
	Painting        bool           `json:"painting"`
	TappingCount    int            `json:"tapping_count"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
