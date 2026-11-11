package mcpserver

import (
	"context"
	"fmt"

	"example.com/go-mcp-enterprise/internal/inventory"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type InventoryService interface {
	GetItem(
		ctx context.Context,
		sku string,
	) (*inventory.Item, error)

	ListLowStock(
		ctx context.Context,
		threshold int,
	) ([]inventory.Item, error)
}

type GetInventoryInput struct {
	SKU string `json:"sku" jsonschema:"商品唯一 SKU 编号"`
}

type InventoryItemOutput struct {
	SKU               string `json:"sku"`
	Name              string `json:"name"`
	TotalQuantity     int    `json:"total_quantity"`
	AvailableQuantity int    `json:"available_quantity"`
	ReservedQuantity  int    `json:"reserved_quantity"`
	ReorderPoint      int    `json:"reorder_point"`
}

type ListLowStockInput struct {
	Threshold int `json:"threshold" jsonschema:"可用库存阈值，返回可用库存小于或等于该值的商品"`
}

type ListLowStockOutput struct {
	Threshold int                   `json:"threshold"`
	Count     int                   `json:"count"`
	Items     []InventoryItemOutput `json:"items"`
}

func RegisterInventoryTools(
	server *mcp.Server,
	service InventoryService,
) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "inventory_get_item",
			Description: "根据商品 SKU 查询库存详情，包括总库存、可用库存、已预留库存和补货预警数量。只读取数据，不修改库存。",
		},
		func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input GetInventoryInput,
		) (*mcp.CallToolResult, InventoryItemOutput, error) {
			item, err := service.GetItem(ctx, input.SKU)
			if err != nil {
				return nil, InventoryItemOutput{},
					fmt.Errorf("get inventory item: %w", err)
			}

			return nil, toInventoryOutput(item), nil
		},
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "inventory_list_low_stock",
			Description: "查询可用库存小于或等于指定阈值的商品，按 SKU 排序返回。只读取数据，不修改库存。",
		},
		func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input ListLowStockInput,
		) (*mcp.CallToolResult, ListLowStockOutput, error) {
			items, err := service.ListLowStock(
				ctx,
				input.Threshold,
			)
			if err != nil {
				return nil, ListLowStockOutput{},
					fmt.Errorf("list low stock inventory: %w", err)
			}

			output := ListLowStockOutput{
				Threshold: input.Threshold,
				Count:     len(items),
				Items:     make([]InventoryItemOutput, 0, len(items)),
			}

			for i := range items {
				output.Items = append(
					output.Items,
					toInventoryOutput(&items[i]),
				)
			}

			return nil, output, nil
		},
	)
}

func toInventoryOutput(item *inventory.Item) InventoryItemOutput {
	return InventoryItemOutput{
		SKU:               item.SKU,
		Name:              item.Name,
		TotalQuantity:     item.TotalQuantity,
		AvailableQuantity: item.AvailableQuantity,
		ReservedQuantity:  item.ReservedQuantity,
		ReorderPoint:      item.ReorderPoint,
	}
}
