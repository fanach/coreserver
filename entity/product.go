package entity

const (
	// PriceUnitRMB yuan
	PriceUnitRMB = "￥"
	// PriceUnitUS dollor
	PriceUnitUS = "$"

	// DataFlowUnitMB MB
	DataFlowUnitMB = "MB"
	// DataFlowUnitGB GB
	DataFlowUnitGB = "GB"

	// ExpireUnitMonth Month
	ExpireUnitMonth = "Month"
	// ExpireUnitYear Year
	ExpireUnitYear = "Year"
)

// Product is what's in sale
// Price 5
// PriceUnit ￥
// DataFlow 1024
// DataFlowUnit MB/GB
// Expire 1
// ExpireUnit Month/Day/Year
type Product struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float32 `json:"price"`
	PriceUnit    string  `json:"price_unit"`
	DataFlow     float32 `json:"dataflow"`
	DataFlowUnit string  `json:"dataflow_unit"`
	Expire       float32 `json:"expire"`
	ExpireUnit   string  `json:"expire_unit"`
}
