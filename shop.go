package regru

import "context"

// ShopService handles domain shop-related API calls.
type ShopService struct {
	client *Client
}

// Shop returns the shop service.
func (c *Client) Shop() *ShopService {
	return &ShopService{client: c}
}

// Nop tests shop API availability.
func (s *ShopService) Nop(ctx context.Context) error {
	return s.client.call(ctx, "shop", "nop", nil, nil)
}

// GetInfo returns shop information for a domain (reseller only).
func (s *ShopService) GetInfo(ctx context.Context, domainName string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	var answer map[string]interface{}
	if err := s.client.call(ctx, "shop", "get_info", params, &answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// Enable enables a domain for sale in the shop.
func (s *ShopService) Enable(ctx context.Context, domainName string, price float64, categoryID int) error {
	params := map[string]interface{}{
		"domain_name": domainName,
		"price":       price,
		"category_id": categoryID,
	}
	return s.client.call(ctx, "shop", "enable", params, nil)
}

// Disable removes a domain from the shop.
func (s *ShopService) Disable(ctx context.Context, domainName string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	return s.client.call(ctx, "shop", "disable", params, nil)
}

// GetCategories returns domain shop categories.
func (s *ShopService) GetCategories(ctx context.Context) ([]map[string]interface{}, error) {
	var answer struct {
		Categories []map[string]interface{} `json:"categories"`
	}
	if err := s.client.call(ctx, "shop", "get_categories", nil, &answer); err != nil {
		return nil, err
	}
	return answer.Categories, nil
}

// GetSuggestedTags returns popular tags for shop lots.
func (s *ShopService) GetSuggestedTags(ctx context.Context, limit int) ([]string, error) {
	params := map[string]interface{}{}
	if limit > 0 {
		params["limit"] = limit
	}
	var answer struct {
		Tags []string `json:"tags"`
	}
	if err := s.client.call(ctx, "shop", "get_suggested_tags", params, &answer); err != nil {
		return nil, err
	}
	return answer.Tags, nil
}
