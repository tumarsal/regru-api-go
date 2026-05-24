package regru

import (
	"context"
	"fmt"
	"net/url"
)

// UserService handles user-related API calls.
type UserService struct {
	client *Client
}

// User returns the user service.
func (c *Client) User() *UserService {
	return &UserService{client: c}
}

// Nop tests user API availability.
func (s *UserService) Nop(ctx context.Context) error {
	return s.client.call(ctx, "user", "nop", nil, nil)
}

// CreateUserRequest is a request to create a new user.
type CreateUserRequest struct {
	UserLogin         string   `json:"user_login"`
	UserPassword      string   `json:"user_password"`
	UserEmail         string   `json:"user_email"`
	UserIP            string   `json:"user_ip"`
	DefaultCountryCode string  `json:"default_country_code,omitempty"`

	// Simplified registration - only email
	EmailOnly bool `json:"email_only,omitempty"`

	// Optional profile fields
	UserFirstName string `json:"user_first_name,omitempty"`
	UserLastName  string `json:"user_last_name,omitempty"`
	UserCompany   string `json:"user_company,omitempty"`
	UserPhone     string `json:"user_phone,omitempty"`
	UserFax       string `json:"user_fax,omitempty"`
	UserAddr      string `json:"user_addr,omitempty"`
	UserCity      string `json:"user_city,omitempty"`
	UserState     string `json:"user_state,omitempty"`
	UserPostcode  string `json:"user_postcode,omitempty"`
	UserLanguage  string `json:"user_language,omitempty"` // ru or en

	// Other options
	UserSubscribe    int      `json:"user_subsribe,omitempty"` // note: typo matches API
	UserMailNotify   int      `json:"user_mailnotify,omitempty"`
	SetMeAsReferrer  int      `json:"set_me_as_referrer,omitempty"`
	CheckOnly        int      `json:"check_only,omitempty"`
	WhiteListIPs     []string `json:"white_list_ips,omitempty"`
}

// CreateUser registers a new user (reseller only).
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (string, error) {
	params := map[string]interface{}{
		"user_email": req.UserEmail,
		"user_ip":    req.UserIP,
	}

	if req.EmailOnly {
		params["email_only"] = "1"
	} else {
		params["user_login"] = req.UserLogin
		params["user_password"] = req.UserPassword
	}

	if req.DefaultCountryCode != "" {
		params["default_country_code"] = req.DefaultCountryCode
	}
	if req.UserFirstName != "" {
		params["user_first_name"] = req.UserFirstName
	}
	if req.UserLastName != "" {
		params["user_last_name"] = req.UserLastName
	}
	if req.UserCompany != "" {
		params["user_company"] = req.UserCompany
	}
	if req.UserPhone != "" {
		params["user_phone"] = req.UserPhone
	}
	if req.UserFax != "" {
		params["user_fax"] = req.UserFax
	}
	if req.UserAddr != "" {
		params["user_addr"] = req.UserAddr
	}
	if req.UserCity != "" {
		params["user_city"] = req.UserCity
	}
	if req.UserState != "" {
		params["user_state"] = req.UserState
	}
	if req.UserPostcode != "" {
		params["user_postcode"] = req.UserPostcode
	}
	if req.UserLanguage != "" {
		params["user_language"] = req.UserLanguage
	}
	if req.UserSubscribe != 0 {
		params["user_subsribe"] = req.UserSubscribe
	}
	if req.UserMailNotify != 0 {
		params["user_mailnotify"] = req.UserMailNotify
	}
	if req.SetMeAsReferrer != 0 {
		params["set_me_as_referrer"] = req.SetMeAsReferrer
	}
	if req.CheckOnly != 0 {
		params["check_only"] = req.CheckOnly
	}
	if len(req.WhiteListIPs) > 0 {
		for _, ip := range req.WhiteListIPs {
			params["white_list_ips"] = append(
				params["white_list_ips"].([]string), ip)
		}
	}

	var answer struct {
		UserID string `json:"user_id"`
	}
	if err := s.client.call(ctx, "user", "create", params, &answer); err != nil {
		return "", err
	}
	return answer.UserID, nil
}

// GetStatistics returns user account statistics.
func (s *UserService) GetStatistics(ctx context.Context) (*UserStatistics, error) {
	var stats UserStatistics
	if err := s.client.call(ctx, "user", "get_statistics", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetBalance returns the user's account balance.
func (s *UserService) GetBalance(ctx context.Context) (*UserBalance, error) {
	var balance UserBalance
	if err := s.client.call(ctx, "user", "get_balance", nil, &balance); err != nil {
		return nil, err
	}
	return &balance, nil
}

// SetResellerURL sets the redirect URL for external services (reseller only).
func (s *UserService) SetResellerURL(ctx context.Context, redirectURL string) error {
	params := map[string]interface{}{
		"url": redirectURL,
	}
	return s.client.call(ctx, "user", "set_reseller_url", params, nil)
}

// GetResellerURL returns the redirect URL for external services (reseller only).
func (s *UserService) GetResellerURL(ctx context.Context) (string, error) {
	var answer struct {
		URL string `json:"url"`
	}
	if err := s.client.call(ctx, "user", "get_reseller_url", nil, &answer); err != nil {
		return "", err
	}
	return answer.URL, nil
}

// GetPersons returns contact persons with their identification status.
func (s *UserService) GetPersons(ctx context.Context) ([]Person, error) {
	var answer struct {
		Persons []Person `json:"persons"`
	}
	if err := s.client.call(ctx, "user", "get_persons", nil, &answer); err != nil {
		return nil, err
	}
	return answer.Persons, nil
}
