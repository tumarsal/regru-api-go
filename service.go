package regru

import (
	"context"
	"fmt"
)

// ServiceService handles service-related API calls.
type ServiceService struct {
	client *Client
}

// Service returns the service API handler.
func (c *Client) Service() *ServiceService {
	return &ServiceService{client: c}
}

// Nop tests service API availability.
func (s *ServiceService) Nop(ctx context.Context, services ...ServiceID) error {
	params := make(map[string]interface{})
	if len(services) > 0 {
		params["services"] = services
	}
	return s.client.call(ctx, "service", "nop", params, nil)
}

// GetPrices returns prices for service activation/renewal.
func (s *ServiceService) GetPrices(ctx context.Context, serviceType string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if serviceType != "" {
		params["servtype"] = serviceType
	}
	var answer map[string]interface{}
	if err := s.client.call(ctx, "service", "get_prices", params, &answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// GetServTypeDetails returns price and general data for a service type.
func (s *ServiceService) GetServTypeDetails(ctx context.Context, serviceType string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"servtype": serviceType,
	}
	var answer map[string]interface{}
	if err := s.client.call(ctx, "service", "get_servtype_details", params, &answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// CreateServiceRequest is a request to order a new service.
type CreateServiceRequest struct {
	DomainName   string                 `json:"dname"`
	ServiceType  string                 `json:"servtype"`
	SubType      string                 `json:"subtype,omitempty"`
	Period       int                    `json:"period,omitempty"`
	PayType      string                 `json:"pay_type,omitempty"`
	OkIfNoMoney  int                    `json:"ok_if_no_money,omitempty"`
	ServiceID    string                 `json:"service_id,omitempty"`
	// Service-specific parameters
	ExtraParams  map[string]interface{} `json:"-"`
}

// Create orders a new service.
func (s *ServiceService) Create(ctx context.Context, req CreateServiceRequest) (*ServiceCreateResponse, error) {
	params := map[string]interface{}{
		"dname":     req.DomainName,
		"servtype":  req.ServiceType,
	}
	if req.SubType != "" {
		params["subtype"] = req.SubType
	}
	if req.Period > 0 {
		params["period"] = req.Period
	}
	if req.PayType != "" {
		params["pay_type"] = req.PayType
	}
	if req.OkIfNoMoney != 0 {
		params["ok_if_no_money"] = req.OkIfNoMoney
	}
	if req.ServiceID != "" {
		params["service_id"] = req.ServiceID
	}
	for k, v := range req.ExtraParams {
		params[k] = v
	}

	var resp ServiceCreateResponse
	if err := s.client.call(ctx, "service", "create", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ServiceCreateResponse is the response for service creation.
type ServiceCreateResponse struct {
	Descr     string `json:"descr"`
	Payment   string `json:"payment"`
	PayNotes  string `json:"pay_notes"`
	BillID    string `json:"bill_id"`
	ServiceID string `json:"service_id"`
}

// CheckCreate validates service order parameters without actually ordering.
func (s *ServiceService) CheckCreate(ctx context.Context, req CreateServiceRequest) error {
	params := map[string]interface{}{
		"dname":     req.DomainName,
		"servtype":  req.ServiceType,
	}
	if req.SubType != "" {
		params["subtype"] = req.SubType
	}
	if req.Period > 0 {
		params["period"] = req.Period
	}
	for k, v := range req.ExtraParams {
		params[k] = v
	}
	return s.client.call(ctx, "service", "check_create", params, nil)
}

// Delete deletes a service.
func (s *ServiceService) Delete(ctx context.Context, service ServiceID) error {
	params := serviceToMap(service)
	return s.client.call(ctx, "service", "delete", params, nil)
}

// GetInfoRequest is a request to get service information.
type GetInfoRequest struct {
	Services       []ServiceID `json:"services,omitempty"`
	Domains        []ServiceID `json:"domains,omitempty"`
	ShowHidden     int         `json:"show_hidden,omitempty"`
	ShowContactsOnly int       `json:"show_contacts_only,omitempty"`
	SeparateGroups int         `json:"separate_groups,omitempty"`
}

// GetInfo returns information about services.
func (s *ServiceService) GetInfo(ctx context.Context, req GetInfoRequest) ([]ServiceInfo, error) {
	params := make(map[string]interface{})
	if len(req.Services) > 0 {
		params["services"] = req.Services
	}
	if len(req.Domains) > 0 {
		params["domains"] = req.Domains
	}
	if req.ShowHidden != 0 {
		params["show_hidden"] = req.ShowHidden
	}
	if req.ShowContactsOnly != 0 {
		params["show_contacts_only"] = req.ShowContactsOnly
	}
	if req.SeparateGroups != 0 {
		params["separate_groups"] = req.SeparateGroups
	}

	var answer struct {
		Services []ServiceInfo `json:"services"`
	}
	if err := s.client.call(ctx, "service", "get_info", params, &answer); err != nil {
		return nil, err
	}
	return answer.Services, nil
}

// GetList returns a list of active services.
func (s *ServiceService) GetList(ctx context.Context, serviceType string) ([]ServiceListItem, error) {
	params := map[string]interface{}{}
	if serviceType != "" {
		params["servtype"] = serviceType
	}
	var answer struct {
		Services []ServiceListItem `json:"services"`
	}
	if err := s.client.call(ctx, "service", "get_list", params, &answer); err != nil {
		return nil, err
	}
	return answer.Services, nil
}

// GetFolders returns folders that contain the service.
func (s *ServiceService) GetFolders(ctx context.Context, service ServiceID) ([]Folder, error) {
	params := serviceToMap(service)
	var answer struct {
		Folders []Folder `json:"folders"`
	}
	if err := s.client.call(ctx, "service", "get_folders", params, &answer); err != nil {
		return nil, err
	}
	return answer.Folders, nil
}

// GetDetails returns general information about a service.
func (s *ServiceService) GetDetails(ctx context.Context, service ServiceID) (map[string]interface{}, error) {
	params := serviceToMap(service)
	var answer map[string]interface{}
	if err := s.client.call(ctx, "service", "get_details", params, &answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// Renew renews a domain or service.
func (s *ServiceService) Renew(ctx context.Context, service ServiceID, period int, payType string) error {
	params := serviceToMap(service)
	params["period"] = period
	if payType != "" {
		params["pay_type"] = payType
	}
	return s.client.call(ctx, "service", "renew", params, nil)
}

// GetBills returns invoices linked to services (reseller only).
func (s *ServiceService) GetBills(ctx context.Context, services ...ServiceID) ([]Bill, error) {
	params := map[string]interface{}{}
	if len(services) > 0 {
		params["services"] = services
	}
	var answer struct {
		Bills []Bill `json:"bills"`
	}
	if err := s.client.call(ctx, "service", "get_bills", params, &answer); err != nil {
		return nil, err
	}
	return answer.Bills, nil
}

// SetAutorenewFlag sets or removes the autorenew flag.
func (s *ServiceService) SetAutorenewFlag(ctx context.Context, service ServiceID, flagValue int) error {
	params := serviceToMap(service)
	params["flag_value"] = flagValue
	return s.client.call(ctx, "service", "set_autorenew_flag", params, nil)
}

// Suspend suspends a service (for domains - removes delegation).
func (s *ServiceService) Suspend(ctx context.Context, service ServiceID) error {
	params := serviceToMap(service)
	return s.client.call(ctx, "service", "suspend", params, nil)
}

// Resume resumes a service (for domains - delegates).
func (s *ServiceService) Resume(ctx context.Context, service ServiceID) error {
	params := serviceToMap(service)
	return s.client.call(ctx, "service", "resume", params, nil)
}

// Upgrade upgrades the service subtype (plan).
func (s *ServiceService) Upgrade(ctx context.Context, service ServiceID, newSubType string) error {
	params := serviceToMap(service)
	params["subtype"] = newSubType
	return s.client.call(ctx, "service", "upgrade", params, nil)
}

// PartControlGrant grants partial control to another user.
func (s *ServiceService) PartControlGrant(ctx context.Context, service ServiceID, grantToUser string) error {
	params := serviceToMap(service)
	params["grant_to"] = grantToUser
	return s.client.call(ctx, "service", "partcontrol_grant", params, nil)
}

// PartControlRevoke revokes partial control from another user.
func (s *ServiceService) PartControlRevoke(ctx context.Context, service ServiceID, revokeFromUser string) error {
	params := serviceToMap(service)
	params["revoke_from"] = revokeFromUser
	return s.client.call(ctx, "service", "partcontrol_revoke", params, nil)
}
