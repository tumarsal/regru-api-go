package regru

import (
	"context"
	"fmt"
)

// DomainService handles domain-related API calls.
type DomainService struct {
	client *Client
}

// Domain returns the domain service.
func (c *Client) Domain() *DomainService {
	return &DomainService{client: c}
}

// Nop tests domain API availability.
func (s *DomainService) Nop(ctx context.Context) error {
	return s.client.call(ctx, "domain", "nop", nil, nil)
}

// GetPrices returns prices for all domain zones.
func (s *DomainService) GetPrices(ctx context.Context, currency string) (map[string]DomainPrice, error) {
	params := map[string]interface{}{}
	if currency != "" {
		params["currency"] = currency
	}
	var answer map[string]DomainPrice
	if err := s.client.call(ctx, "domain", "get_prices", params, &answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// GetSuggest suggests domain names based on keywords (reseller only).
func (s *DomainService) GetSuggest(ctx context.Context, word string, tlds []string) ([]string, error) {
	params := map[string]interface{}{
		"word": word,
	}
	if len(tlds) > 0 {
		params["tlds"] = tlds
	}
	var answer struct {
		Domains []string `json:"domains"`
	}
	if err := s.client.call(ctx, "domain", "get_suggest", params, &answer); err != nil {
		return nil, err
	}
	return answer.Domains, nil
}

// GetPremiumPrices returns premium domain prices (reseller only).
func (s *DomainService) GetPremiumPrices(ctx context.Context) (map[string]interface{}, error) {
	var answer map[string]interface{}
	if err := s.client.call(ctx, "domain", "get_premium_prices", nil, &answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// GetDeleted returns a list of deleted/expired domains (reseller only).
func (s *DomainService) GetDeleted(ctx context.Context, req GetDeletedRequest) ([]DeletedDomain, error) {
	params := map[string]interface{}{}
	if len(req.TLDs) > 0 {
		params["tlds"] = req.TLDs
	}
	if req.DeletedFrom != "" {
		params["deleted_from"] = req.DeletedFrom
	}
	if req.DeletedTo != "" {
		params["deleted_to"] = req.DeletedTo
	}
	if req.CreatedFrom != "" {
		params["created_from"] = req.CreatedFrom
	}
	if req.CreatedTo != "" {
		params["created_to"] = req.CreatedTo
	}
	if req.HideReg != 0 {
		params["hidereg"] = req.HideReg
	}
	if req.MinPR != 0 {
		params["min_pr"] = req.MinPR
	}
	if req.MinCY != 0 {
		params["min_cy"] = req.MinCY
	}

	var answer struct {
		Domains []DeletedDomain `json:"domains"`
	}
	if err := s.client.call(ctx, "domain", "get_deleted", params, &answer); err != nil {
		return nil, err
	}
	return answer.Domains, nil
}

// GetDeletedRequest is a request for deleted domains.
type GetDeletedRequest struct {
	TLDs        []string `json:"tlds,omitempty"`         // ru, рф, su
	DeletedFrom string   `json:"deleted_from,omitempty"` // YYYY-MM-DD
	DeletedTo   string   `json:"deleted_to,omitempty"`
	CreatedFrom string   `json:"created_from,omitempty"`
	CreatedTo   string   `json:"created_to,omitempty"`
	HideReg     int      `json:"hidereg,omitempty"`      // 1 = only available
	MinPR       int      `json:"min_pr,omitempty"`       // Google PR
	MinCY       int      `json:"min_cy,omitempty"`       // Yandex TIC
}

// Check checks domain availability for registration.
func (s *DomainService) Check(ctx context.Context, domainName string, isTransfer bool) (*DomainCheckResult, error) {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	if isTransfer {
		params["is_transfer"] = 1
	}
	var answer struct {
		Domains []DomainCheckResult `json:"domains"`
	}
	if err := s.client.call(ctx, "domain", "check", params, &answer); err != nil {
		return nil, err
	}
	if len(answer.Domains) > 0 {
		return &answer.Domains[0], nil
	}
	return nil, fmt.Errorf("no result in response")
}

// CheckBatch checks multiple domains at once.
func (s *DomainService) CheckBatch(ctx context.Context, domainNames []string, isTransfer bool) ([]DomainCheckResult, error) {
	domains := make([]map[string]string, len(domainNames))
	for i, d := range domainNames {
		domains[i] = map[string]string{"dname": d}
	}
	params := map[string]interface{}{
		"domains": domains,
	}
	if isTransfer {
		params["is_transfer"] = 1
	}
	var answer struct {
		Domains []DomainCheckResult `json:"domains"`
	}
	if err := s.client.call(ctx, "domain", "check", params, &answer); err != nil {
		return nil, err
	}
	return answer.Domains, nil
}

// CreateDomainRequest is a request to register a domain.
type CreateDomainRequest struct {
	DomainName        string         `json:"domain_name"`
	Period            int            `json:"period,omitempty"`      // default 1 year
	Nss               []NameServer   `json:"nss,omitempty"`         // custom nameservers
	UseDefaultNss     int            `json:"use_default_nss,omitempty"` // 1 = use default NS
	Contacts          DomainContact  `json:"contacts"`
	PrivatePersonFlag int            `json:"private_person_flag,omitempty"`
	PayType           string         `json:"pay_type,omitempty"`
	OkIfNoMoney       int            `json:"ok_if_no_money,omitempty"`
	AutoRenew         int            `json:"auto_renew,omitempty"`

	// For .RU/.SU/.РФ domains
	OrgRF     int    `json:"org_ru,omitempty"`

	// For some TLDs
	Description string `json:"description,omitempty"`
	Keywords    string `json:"keywords,omitempty"`

	// For premium domains
	PremiumPrice float64 `json:"premium_price,omitempty"`
}

// Create registers a new domain.
func (s *DomainService) Create(ctx context.Context, req CreateDomainRequest) (*DomainCreateResponse, error) {
	params := map[string]interface{}{
		"domain_name": req.DomainName,
	}
	if req.Period > 0 {
		params["period"] = req.Period
	}
	if len(req.Nss) > 0 {
		params["nss"] = req.Nss
	}
	if req.UseDefaultNss != 0 {
		params["use_default_nss"] = req.UseDefaultNss
	}
	if req.PrivatePersonFlag != 0 {
		params["private_person_flag"] = req.PrivatePersonFlag
	}
	if req.PayType != "" {
		params["pay_type"] = req.PayType
	}
	if req.OkIfNoMoney != 0 {
		params["ok_if_no_money"] = req.OkIfNoMoney
	}
	if req.AutoRenew != 0 {
		params["auto_renew"] = req.AutoRenew
	}
	if req.OrgRF != 0 {
		params["org_ru"] = req.OrgRF
	}
	if req.Description != "" {
		params["description"] = req.Description
	}
	if req.Keywords != "" {
		params["keywords"] = req.Keywords
	}
	if req.PremiumPrice > 0 {
		params["premium_price"] = req.PremiumPrice
	}

	// Add contacts
	addContactParams(params, req.Contacts)

	var resp DomainCreateResponse
	if err := s.client.call(ctx, "domain", "create", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DomainCreateResponse is the response for domain creation.
type DomainCreateResponse struct {
	ServiceID string `json:"service_id,omitempty"`
	Descr     string `json:"descr,omitempty"`
	Payment   string `json:"payment,omitempty"`
	BillID    string `json:"bill_id,omitempty"`
	Result    string `json:"result,omitempty"`
}

// Transfer initiates domain transfer from another registrar.
func (s *DomainService) Transfer(ctx context.Context, domainName, authInfo string, contacts *DomainContact, period int) (*DomainCreateResponse, error) {
	params := map[string]interface{}{
		"domain_name": domainName,
		"authinfo":    authInfo,
	}
	if period > 0 {
		params["period"] = period
	}
	if contacts != nil {
		addContactParams(params, *contacts)
	}
	var resp DomainCreateResponse
	if err := s.client.call(ctx, "domain", "transfer", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTransferStatus returns the status of a domain transfer.
func (s *DomainService) GetTransferStatus(ctx context.Context, domainName string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	var answer map[string]interface{}
	if err := s.client.call(ctx, "domain", "get_transfer_status", params, &answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// SetNewAuthInfo sets a new authorization code for domain transfers.
func (s *DomainService) SetNewAuthInfo(ctx context.Context, domainName, authInfo string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
		"authinfo":    authInfo,
	}
	return s.client.call(ctx, "domain", "set_new_authinfo", params, nil)
}

// CancelTransfer cancels a domain transfer.
func (s *DomainService) CancelTransfer(ctx context.Context, domainName string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	return s.client.call(ctx, "domain", "cancel_transfer", params, nil)
}

// GetReregData returns a list of dropping domains with characteristics (reseller only).
func (s *DomainService) GetReregData(ctx context.Context) (map[string]interface{}, error) {
	var answer map[string]interface{}
	if err := s.client.call(ctx, "domain", "get_rereg_data", nil, &answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// SetReregBids places bids on dropping domains.
func (s *DomainService) SetReregBids(ctx context.Context, bids []ReregBid) error {
	params := map[string]interface{}{
		"bids": bids,
	}
	return s.client.call(ctx, "domain", "set_rereg_bids", params, nil)
}

// GetUserReregBids returns the user's bids on dropping domains.
func (s *DomainService) GetUserReregBids(ctx context.Context) ([]ReregBid, error) {
	var answer struct {
		Bids []ReregBid `json:"bids"`
	}
	if err := s.client.call(ctx, "domain", "get_user_rereg_bids", nil, &answer); err != nil {
		return nil, err
	}
	return answer.Bids, nil
}

// GetDocsUploadURI returns a URL for uploading documents for .RU/.SU/.РФ domains.
func (s *DomainService) GetDocsUploadURI(ctx context.Context, domainName string) (string, error) {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	var answer struct {
		UploadURI string `json:"upload_uri"`
	}
	if err := s.client.call(ctx, "domain", "get_docs_upload_uri", params, &answer); err != nil {
		return "", err
	}
	return answer.UploadURI, nil
}

// UpdateContacts updates domain contact data.
func (s *DomainService) UpdateContacts(ctx context.Context, domainName string, contacts DomainContact) error {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	addContactParams(params, contacts)
	return s.client.call(ctx, "domain", "update_contacts", params, nil)
}

// UpdatePrivatePersonFlag changes the private person flag in WHOIS.
func (s *DomainService) UpdatePrivatePersonFlag(ctx context.Context, domainName string, flag int) error {
	params := map[string]interface{}{
		"domain_name":           domainName,
		"private_person_flag":   flag,
	}
	return s.client.call(ctx, "domain", "update_private_person_flag", params, nil)
}

// RegisterNS registers a nameserver in NSI-registry.
func (s *DomainService) RegisterNS(ctx context.Context, domainName, nsName, nsIP string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
		"nsname":      nsName,
		"nsip":        nsIP,
	}
	return s.client.call(ctx, "domain", "register_ns", params, nil)
}

// DeleteNS deletes a nameserver from NSI-registry.
func (s *DomainService) DeleteNS(ctx context.Context, domainName, nsName string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
		"nsname":      nsName,
	}
	return s.client.call(ctx, "domain", "delete_ns", params, nil)
}

// GetNSs returns the list of nameservers for a domain.
func (s *DomainService) GetNSs(ctx context.Context, domainName string) ([]NameServer, error) {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	var answer struct {
		Nss []NameServer `json:"nss"`
	}
	if err := s.client.call(ctx, "domain", "get_nss", params, &answer); err != nil {
		return nil, err
	}
	return answer.Nss, nil
}

// UpdateNSs updates the nameservers for a domain.
func (s *DomainService) UpdateNSs(ctx context.Context, domainName string, nss []NameServer) error {
	params := map[string]interface{}{
		"domain_name": domainName,
		"nss":         nss,
	}
	return s.client.call(ctx, "domain", "update_nss", params, nil)
}

// Delegate delegates a domain.
func (s *DomainService) Delegate(ctx context.Context, domainName string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	return s.client.call(ctx, "domain", "delegate", params, nil)
}

// Undelegate removes delegation from a domain.
func (s *DomainService) Undelegate(ctx context.Context, domainName string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	return s.client.call(ctx, "domain", "undelegate", params, nil)
}

// TransferToAnotherAccount transfers a domain to another account.
func (s *DomainService) TransferToAnotherAccount(ctx context.Context, domainName, targetUsername string) error {
	params := map[string]interface{}{
		"domain_name":     domainName,
		"target_username": targetUsername,
	}
	return s.client.call(ctx, "domain", "transfer_to_another_account", params, nil)
}

// LookAtEnteringList views domains being transferred to the account.
func (s *DomainService) LookAtEnteringList(ctx context.Context) ([]map[string]interface{}, error) {
	var answer struct {
		Domains []map[string]interface{} `json:"domains"`
	}
	if err := s.client.call(ctx, "domain", "look_at_entering_list", nil, &answer); err != nil {
		return nil, err
	}
	return answer.Domains, nil
}

// addContactParams adds contact fields to params map.
func addContactParams(params map[string]interface{}, c DomainContact) {
	fields := map[string]string{
		"r_last_name": c.RLastName,
		"r_first_name": c.RFirstName,
		"r_middle_name": c.RMiddleName,
		"r_phone": c.RPhone,
		"r_fax": c.RFax,
		"r_email": c.REmail,
		"r_country": c.RCountry,
		"r_state": c.RState,
		"r_city": c.RCity,
		"r_addr": c.RAddr,
		"r_postcode": c.RPostcode,
		"r_org_name": c.ROrgName,
		"i_last_name": c.ILastName,
		"i_first_name": c.IFirstName,
		"i_middle_name": c.IMiddleName,
		"i_phone": c.IPhone,
		"i_fax": c.IFax,
		"i_email": c.IEmail,
		"i_country": c.ICountry,
		"i_state": c.IState,
		"i_city": c.ICity,
		"i_addr": c.IAddr,
		"i_postcode": c.IPostcode,
		"i_org_name": c.IOrgName,
		"a_last_name": c.ALastName,
		"a_first_name": c.AFirstName,
		"a_middle_name": c.AMiddleName,
		"a_phone": c.APhone,
		"a_fax": c.AFax,
		"a_email": c.AEmail,
		"a_country": c.ACountry,
		"a_state": c.AState,
		"a_city": c.ACity,
		"a_addr": c.AAddr,
		"a_postcode": c.APostcode,
		"a_org_name": c.AOrgName,
		"birth_date": c.BirthDate,
		"passport": c.Passport,
		"person": c.Person,
	}
	for k, v := range fields {
		if v != "" {
			params[k] = v
		}
	}
}
