// Package regru provides a Go client for REG.RU API 2.0.
package regru

import "encoding/json"

const (
	// DefaultBaseURL is the production API endpoint.
	DefaultBaseURL = "https://api.reg.ru"
	// TestBaseURL can be used for testing (same host, credentials "test"/"test").
	DefaultBaseURLPath = "/api/regru2"

	// Output formats.
	OutputJSON = "json"
	OutputYAML = "yaml"
	OutputXML  = "xml"
	OutputPlain = "plain"

	// Input formats.
	InputJSON = "json"
	InputXML  = "xml"
	// InputPlain means plain HTTP form parameters (default).
	InputPlain = "plain"

	// Pay types.
	PayTypePrepay = "prepay"
	PayTypeBank   = "bank"
	PayTypePBank  = "pbank"
	PayTypeYaCard = "yacard"

	// Currencies.
	CurrencyRUR = "RUR"
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	CurrencyUAH = "UAH"

	// Service types.
	ServTypeDomain          = "domain"
	ServTypeHostingISP      = "srv_hosting_ispmgr"
	ServTypeWebFwd          = "srv_webfwd"
	ServTypeParking         = "srv_parking"
	ServTypeSSLCertificate  = "srv_ssl_certificate"
	ServTypeVPS             = "srv_vps"
	ServTypeDedicatedServer = "srv_dedicated"

	// Result values.
	ResultSuccess = "success"
	ResultError   = "error"
)

// Response is the top-level response wrapper for all API calls.
type Response struct {
	Result       string          `json:"result"`
	Answer       json.RawMessage `json:"answer,omitempty"`
	ErrorCode    string          `json:"error_code,omitempty"`
	ErrorText    string          `json:"error_text,omitempty"`
	ErrorDetails string          `json:"error_details,omitempty"`
	Charset      string          `json:"charset,omitempty"`
	MessageStore interface{}     `json:"messagestore,omitempty"`
}

// IsSuccess returns true if the API call was successful.
func (r *Response) IsSuccess() bool {
	return r.Result == ResultSuccess
}

// Error returns an error description if the call failed.
func (r *Response) Error() string {
	if r.IsSuccess() {
		return ""
	}
	if r.ErrorText != "" {
		return r.ErrorCode + ": " + r.ErrorText
	}
	return r.ErrorCode
}

// AuthParams holds authentication fields sent with every request.
type AuthParams struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	// Sig is used for signature-based authentication (alternative to Password).
	Sig string `json:"sig,omitempty"`
}

// InputParams holds API control parameters sent with every request.
type InputParams struct {
	OutputContentType string `json:"output_content_type,omitempty"`
	OutputFormat      string `json:"output_format,omitempty"`
	InputFormat       string `json:"input_format,omitempty"`
	InputData         string `json:"input_data,omitempty"`
	IoEncoding        string `json:"io_encoding,omitempty"`
	ShowInputParams   int    `json:"show_input_params,omitempty"`
	Lang              string `json:"lang,omitempty"`
}

// ServiceID identifies a single service or domain.
type ServiceID struct {
	DomainName      string `json:"domain_name,omitempty"`
	ServiceID       string `json:"service_id,omitempty"`
	ServiceType     string `json:"servtype,omitempty"`
	SubType         string `json:"subtype,omitempty"`
	UpLinkServiceID string `json:"uplink_service_id,omitempty"`
}

// ServiceList is used for batch operations on multiple services.
type ServiceList struct {
	Domains  []ServiceID `json:"domains,omitempty"`
	Services []ServiceID `json:"services,omitempty"`
}

// DomainCheckResult represents the result of a domain availability check.
type DomainCheckResult struct {
	DomainName  string  `json:"dname"`
	Result      string  `json:"result"`
	ErrorCode   string  `json:"error_code,omitempty"`
	ErrorText   string  `json:"error_text,omitempty"`
	IsPremium   bool    `json:"is_premium,omitempty"`
	Price       float64 `json:"price,omitempty"`
}

// UserBalance represents the user account balance.
type UserBalance struct {
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Prepay   float64 `json:"prepay,omitempty"`
}

// Bill represents an invoice.
type Bill struct {
	BillID        string      `json:"bill_id"`
	BillDate      string      `json:"bill_date"`
	Currency      string      `json:"currency"`
	Payment       float64     `json:"payment"`
	TotalPayment  float64     `json:"total_payment"`
	PayType       string      `json:"pay_type"`
	PayStatus     string      `json:"pay_status"`
	Items         []BillItem  `json:"items"`
}

// BillItem represents a line item in an invoice.
type BillItem struct {
	ItemType  string `json:"itemtype"`
	DomainName string `json:"dname,omitempty"`
	ServType  string `json:"servtype,omitempty"`
	ServiceID string `json:"service_id,omitempty"`
	Action    string `json:"action,omitempty"`
}

// ServiceInfo represents detailed information about a service.
type ServiceInfo struct {
	DomainName string                 `json:"dname"`
	ServType   string                 `json:"servtype"`
	ServiceID  string                 `json:"service_id"`
	Result     string                 `json:"result"`
	ErrorCode  string                 `json:"error_code,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Contacts   interface{}            `json:"contacts,omitempty"`
}

// ServiceListItem represents an item in the service list.
type ServiceListItem struct {
	ServiceID   string `json:"service_id"`
	DomainName  string `json:"dname"`
	ServiceType string `json:"servtype"`
	SubType     string `json:"subtype,omitempty"`
}

// Folder represents a folder for organizing services/domains.
type Folder struct {
	FolderID   int    `json:"folder_id"`
	FolderName string `json:"folder_name"`
}

// DomainContact represents contact data for domain registration.
type DomainContact struct {
	// Russian legal/contact person
	RLastName      string `json:"r_last_name,omitempty"`
	RFirstName     string `json:"r_first_name,omitempty"`
	RMiddleName    string `json:"r_middle_name,omitempty"`
	RPhone         string `json:"r_phone,omitempty"`
	RFax           string `json:"r_fax,omitempty"`
	REmail         string `json:"r_email,omitempty"`
	RCountry       string `json:"r_country,omitempty"`
	RState         string `json:"r_state,omitempty"`
	RCity          string `json:"r_city,omitempty"`
	RAddr          string `json:"r_addr,omitempty"`
	RPostcode      string `json:"r_postcode,omitempty"`
	ROrgName       string `json:"r_org_name,omitempty"`

	// International contact
	ILastName      string `json:"i_last_name,omitempty"`
	IFirstName     string `json:"i_first_name,omitempty"`
	IMiddleName    string `json:"i_middle_name,omitempty"`
	IPhone         string `json:"i_phone,omitempty"`
	IFax           string `json:"i_fax,omitempty"`
	IEmail         string `json:"i_email,omitempty"`
	ICountry       string `json:"i_country,omitempty"`
	IState         string `json:"i_state,omitempty"`
	ICity          string `json:"i_city,omitempty"`
	IAddr          string `json:"i_addr,omitempty"`
	IPostcode      string `json:"i_postcode,omitempty"`
	IOrgName       string `json:"i_org_name,omitempty"`

	// Additional contacts for .RU/.SU/.РФ
	ALastName      string `json:"a_last_name,omitempty"`
	AFirstName     string `json:"a_first_name,omitempty"`
	AMiddleName    string `json:"a_middle_name,omitempty"`
	APhone         string `json:"a_phone,omitempty"`
	AFax           string `json:"a_fax,omitempty"`
	AEmail         string `json:"a_email,omitempty"`
	ACountry       string `json:"a_country,omitempty"`
	AState         string `json:"a_state,omitempty"`
	ACity          string `json:"a_city,omitempty"`
	AAddr          string `json:"a_addr,omitempty"`
	APostcode      string `json:"a_postcode,omitempty"`
	AOrgName       string `json:"a_org_name,omitempty"`

	// Birth date and passport for individuals (RU/SU/РФ)
	BirthDate      string `json:"birth_date,omitempty"`
	Passport       string `json:"passport,omitempty"`
	Person         string `json:"person,omitempty"`
}

// NameServer represents a DNS nameserver.
type NameServer struct {
	Name    string `json:"name"`
	IP      string `json:"ip,omitempty"`
}

// DNSRecord represents a DNS resource record.
type DNSRecord struct {
	SubDomain   string `json:"subdomain,omitempty"`
	Content     string `json:"content"`
	RecordType  string `json:"rectype,omitempty"`
	Priority    string `json:"priority,omitempty"`
	TTL         int    `json:"ttl,omitempty"`
}

// ZoneRecord represents a DNS zone record.
type ZoneRecord struct {
	RecordID    int    `json:"record_id"`
	SubDomain   string `json:"subdomain"`
	Content     string `json:"content"`
	RecordType  string `json:"rectype"`
	Priority    string `json:"priority,omitempty"`
}

// DomainPrice represents pricing information for a domain zone.
type DomainPrice struct {
	Zone          string  `json:"zone"`
	PriceRegister float64 `json:"price_register,omitempty"`
	PriceRenew    float64 `json:"price_renew,omitempty"`
	PriceTransfer float64 `json:"price_transfer,omitempty"`
	IsPremium     bool    `json:"is_premium,omitempty"`
}

// DeletedDomain represents a deleted/expired domain.
type DeletedDomain struct {
	DomainName      string `json:"domain_name"`
	DateDelete      string `json:"date_delete"`
	Registered      string `json:"registered"`
	FirstCreateDate string `json:"first_create_date"`
	YandexTIC       int    `json:"yandex_tic"`
	GooglePR        string `json:"google_pr"`
}

// ReregBid represents a bid on a dropping domain.
type ReregBid struct {
	DomainName string  `json:"domain_name"`
	Bid        float64 `json:"bid"`
}

// UserStatistics represents user account statistics.
type UserStatistics struct {
	TotalDomains    int `json:"total_domains,omitempty"`
	TotalServices   int `json:"total_services,omitempty"`
	ActiveDomains   int `json:"active_domains,omitempty"`
}

// Person represents a contact person with identification status.
type Person struct {
	PersonID   string `json:"person_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	IsVerified bool   `json:"is_verified"`
}

// SSLContact represents contact information for SSL certificate orders.
type SSLContact struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Title      string `json:"title,omitempty"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Fax        string `json:"fax,omitempty"`
	OrgName    string `json:"org_name,omitempty"`
}
