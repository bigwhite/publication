package inventory

type Item struct {
	SKU               string `json:"sku"`
	Name              string `json:"name"`
	TotalQuantity     int    `json:"total_quantity"`
	AvailableQuantity int    `json:"available_quantity"`
	ReservedQuantity  int    `json:"reserved_quantity"`
	ReorderPoint      int    `json:"reorder_point"`
}
