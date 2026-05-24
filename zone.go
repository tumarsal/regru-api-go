package regru

import (
	"context"
	"fmt"
)

// ZoneService handles DNS zone-related API calls.
type ZoneService struct {
	client *Client
}

// Zone returns the zone service.
func (c *Client) Zone() *ZoneService {
	return &ZoneService{client: c}
}

// Nop tests zone API availability.
func (s *ZoneService) Nop(ctx context.Context) error {
	return s.client.call(ctx, "zone", "nop", nil, nil)
}

// GetZoneRecords returns all DNS records for a domain.
func (s *ZoneService) GetRecords(ctx context.Context, domainName string) ([]ZoneRecord, error) {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	var answer struct {
		Records []ZoneRecord `json:"records"`
	}
	if err := s.client.call(ctx, "zone", "get_records", params, &answer); err != nil {
		return nil, err
	}
	return answer.Records, nil
}

// AddRecord adds a DNS record to a domain's zone.
func (s *ZoneService) AddRecord(ctx context.Context, domainName string, record DNSRecord) (*ZoneRecord, error) {
	params := map[string]interface{}{
		"domain_name": domainName,
		"subdomain":   record.SubDomain,
		"content":     record.Content,
		"rectype":     record.RecordType,
	}
	if record.Priority != "" {
		params["priority"] = record.Priority
	}
	if record.TTL > 0 {
		params["ttl"] = record.TTL
	}

	var answer struct {
		RecordID int `json:"record_id"`
	}
	if err := s.client.call(ctx, "zone", "add_record", params, &answer); err != nil {
		return nil, err
	}
	return &ZoneRecord{
		RecordID:  answer.RecordID,
		SubDomain: record.SubDomain,
		Content:   record.Content,
		RecordType: record.RecordType,
		Priority:  record.Priority,
	}, nil
}

// UpdateRecord updates an existing DNS record.
func (s *ZoneService) UpdateRecord(ctx context.Context, domainName string, recordID int, record DNSRecord) error {
	params := map[string]interface{}{
		"domain_name": domainName,
		"record_id":   recordID,
		"subdomain":   record.SubDomain,
		"content":     record.Content,
		"rectype":     record.RecordType,
	}
	if record.Priority != "" {
		params["priority"] = record.Priority
	}
	return s.client.call(ctx, "zone", "update_record", params, nil)
}

// DeleteRecord removes a DNS record from a domain's zone.
func (s *ZoneService) DeleteRecord(ctx context.Context, domainName string, recordID int) error {
	params := map[string]interface{}{
		"domain_name": domainName,
		"record_id":   recordID,
	}
	return s.client.call(ctx, "zone", "delete_record", params, nil)
}

// DeleteRecords removes multiple DNS records by type/content.
func (s *ZoneService) DeleteRecords(ctx context.Context, domainName, recordType, subdomain, content string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	if recordType != "" {
		params["rectype"] = recordType
	}
	if subdomain != "" {
		params["subdomain"] = subdomain
	}
	if content != "" {
		params["content"] = content
	}
	return s.client.call(ctx, "zone", "delete_records", params, nil)
}

// ClearZone removes all DNS records from a domain's zone.
func (s *ZoneService) ClearZone(ctx context.Context, domainName string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	return s.client.call(ctx, "zone", "clear", params, nil)
}

// TmplCreate creates a DNS template.
func (s *ZoneService) TmplCreate(ctx context.Context, name string, records []DNSRecord) error {
	params := map[string]interface{}{
		"name":    name,
		"records": records,
	}
	return s.client.call(ctx, "zone", "tmpl_create", params, nil)
}

// TmplGet returns a DNS template.
func (s *ZoneService) TmplGet(ctx context.Context, name string) ([]DNSRecord, error) {
	params := map[string]interface{}{
		"name": name,
	}
	var answer struct {
		Records []DNSRecord `json:"records"`
	}
	if err := s.client.call(ctx, "zone", "tmpl_get", params, &answer); err != nil {
		return nil, err
	}
	return answer.Records, nil
}

// TmplDelete deletes a DNS template.
func (s *ZoneService) TmplDelete(ctx context.Context, name string) error {
	params := map[string]interface{}{
		"name": name,
	}
	return s.client.call(ctx, "zone", "tmpl_delete", params, nil)
}

// TmplApply applies a DNS template to a domain.
func (s *ZoneService) TmplApply(ctx context.Context, domainName, tmplName string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
		"tmpl_name":   tmplName,
	}
	return s.client.call(ctx, "zone", "tmpl_apply", params, nil)
}

// Add_Alias is a helper to add an A/ALIAS record.
func (s *ZoneService) Add_Alias(ctx context.Context, domainName, subdomain, ip string) (*ZoneRecord, error) {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Content:    ip,
		RecordType: "A",
	})
}

// Add_CNAME is a helper to add a CNAME record.
func (s *ZoneService) Add_CNAME(ctx context.Context, domainName, subdomain, target string) (*ZoneRecord, error) {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Content:    target,
		RecordType: "CNAME",
	})
}

// Add_MX is a helper to add an MX record.
func (s *ZoneService) Add_MX(ctx context.Context, domainName, subdomain, mailserver, priority string) (*ZoneRecord, error) {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Content:    mailserver,
		RecordType: "MX",
		Priority:   priority,
	})
}

// Add_TXT is a helper to add a TXT record.
func (s *ZoneService) Add_TXT(ctx context.Context, domainName, subdomain, text string) (*ZoneRecord, error) {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Content:    text,
		RecordType: "TXT",
	})
}

// Add_NS is a helper to add an NS record.
func (s *ZoneService) Add_NS(ctx context.Context, domainName, subdomain, nameserver string) (*ZoneRecord, error) {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Content:    nameserver,
		RecordType: "NS",
	})
}

// Add_SRV is a helper to add an SRV record.
func (s *ZoneService) Add_SRV(ctx context.Context, domainName, subdomain, content, priority string) (*ZoneRecord, error) {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Content:    content,
		RecordType: "SRV",
		Priority:   priority,
	})
}

// Add_AAAA is a helper to add an AAAA (IPv6) record.
func (s *ZoneService) Add_AAAA(ctx context.Context, domainName, subdomain, ipv6 string) (*ZoneRecord, error) {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Content:    ipv6,
		RecordType: "AAAA",
	})
}
