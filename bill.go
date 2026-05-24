package regru

import (
	"context"
	"fmt"
)

// BillService handles invoice-related API calls.
type BillService struct {
	client *Client
}

// Bill returns the bill service.
func (c *Client) Bill() *BillService {
	return &BillService{client: c}
}

// Nop tests bill API availability.
func (s *BillService) Nop(ctx context.Context) error {
	return s.client.call(ctx, "bill", "nop", nil, nil)
}

// GetNotPayed returns a list of unpaid invoices.
func (s *BillService) GetNotPayed(ctx context.Context) ([]Bill, error) {
	var answer struct {
		Bills []Bill `json:"bills"`
	}
	if err := s.client.call(ctx, "bill", "get_not_payed", nil, &answer); err != nil {
		return nil, err
	}
	return answer.Bills, nil
}

// GetForPeriodRequest is a request to get bills for a specific period.
type GetForPeriodRequest struct {
	StartDate string `json:"start_date"` // ISO format: YYYY-MM-DD
	EndDate   string `json:"end_date"`   // ISO format: YYYY-MM-DD
	PayType   string `json:"pay_type,omitempty"`
	Limit     int    `json:"limit,omitempty"`  // default 100, max 1024
	Offset    int    `json:"offset,omitempty"`
	All       int    `json:"all,omitempty"`    // include inactive bills
}

// GetForPeriod returns invoices for a specified period.
func (s *BillService) GetForPeriod(ctx context.Context, req GetForPeriodRequest) ([]Bill, error) {
	params := map[string]interface{}{
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
	}
	if req.PayType != "" {
		params["pay_type"] = req.PayType
	}
	if req.Limit > 0 {
		params["limit"] = req.Limit
	}
	if req.Offset > 0 {
		params["offset"] = req.Offset
	}
	if req.All != 0 {
		params["all"] = req.All
	}

	var answer struct {
		Bills []Bill `json:"bills"`
	}
	if err := s.client.call(ctx, "bill", "get_for_period", params, &answer); err != nil {
		return nil, err
	}
	return answer.Bills, nil
}

// ChangePayTypeRequest is a request to change the payment method.
type ChangePayTypeRequest struct {
	BillID   string   `json:"bill_id,omitempty"`
	Bills    []string `json:"bills,omitempty"`
	PayType  string   `json:"pay_type"`  // prepay, yamoney, bank
	Currency string   `json:"currency"`  // RUR, USD
}

// ChangePayType changes the payment method for an invoice.
func (s *BillService) ChangePayType(ctx context.Context, req ChangePayTypeRequest) error {
	params := map[string]interface{}{
		"pay_type":  req.PayType,
		"currency":  req.Currency,
	}
	if req.BillID != "" {
		params["bill_id"] = req.BillID
	}
	if len(req.Bills) > 0 {
		params["bills"] = req.Bills
	}
	return s.client.call(ctx, "bill", "change_pay_type", params, nil)
}

// DeleteBill deletes an invoice.
func (s *BillService) DeleteBill(ctx context.Context, billID string) error {
	params := map[string]interface{}{
		"bill_id": billID,
	}
	return s.client.call(ctx, "bill", "delete", params, nil)
}
