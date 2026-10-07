package utils

import (
	"fmt"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
)

type OrderData struct {
	OrderID      uint
	CustomerName string
	TotalAmount  float64
}

func GenerateOrderInvoice(order OrderData) ([]byte, error) {
	cfg := config.NewBuilder().Build()
	m := maroto.New(cfg)

	m.AddRows(
		text.NewRow(12, fmt.Sprintf("INVOICE FOR ORDER #%s", order.OrderID)),
		text.NewRow(8, fmt.Sprintf("Customer: %s", order.CustomerName)),
		text.NewRow(8, fmt.Sprintf("Total Amount: $%.2f", order.TotalAmount)),
	)

	doc, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate pdf document: %w", err)
	}

	return doc.GetBytes(), nil
}
