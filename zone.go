package regru

import (
	"context"
	"fmt"
	"strings"
)

// ZoneService handles DNS zone-related API calls.
type ZoneService struct {
	client *Client
}

// Zone returns the zone service.
func (c *Client) Zone() *ZoneService {
	return &ZoneService{client: c}
}

type zoneDomainsAnswer struct {
	Domains []zoneDomainResult `json:"domains"`
}

type zoneDomainResult struct {
	Dname     string             `json:"dname"`
	Result    string             `json:"result"`
	ServiceID string             `json:"service_id,omitempty"`
	ErrorCode string             `json:"error_code,omitempty"`
	ErrorText string             `json:"error_text,omitempty"`
	RRs       []ResourceRecord   `json:"rrs,omitempty"`
}

func (s *ZoneService) zoneCall(ctx context.Context, function string, params map[string]interface{}) error {
	var answer zoneDomainsAnswer
	if err := s.client.call(ctx, "zone", function, params, &answer); err != nil {
		return err
	}
	domainName, _ := params["domain_name"].(string)
	return checkZoneDomainResults(answer.Domains, domainName)
}

func checkZoneDomainResults(domains []zoneDomainResult, expectDname string) error {
	if len(domains) == 0 {
		return nil
	}
	for _, d := range domains {
		if expectDname != "" && d.Dname != expectDname {
			continue
		}
		if d.Result != ResultSuccess {
			msg := d.ErrorCode
			if d.ErrorText != "" {
				msg = d.ErrorCode + ": " + d.ErrorText
			}
			if msg == "" {
				msg = d.Result
			}
			return fmt.Errorf("zone %s: %s", d.Dname, msg)
		}
		return nil
	}
	// single domain without dname match — check first
	if len(domains) == 1 && domains[0].Result != ResultSuccess {
		d := domains[0]
		msg := d.ErrorCode
		if d.ErrorText != "" {
			msg = d.ErrorCode + ": " + d.ErrorText
		}
		return fmt.Errorf("zone: %s", msg)
	}
	return nil
}

// Nop tests zone API availability.
func (s *ZoneService) Nop(ctx context.Context) error {
	return s.client.call(ctx, "zone", "nop", nil, nil)
}

// GetResourceRecords returns DNS resource records for a domain (zone/get_resource_records).
func (s *ZoneService) GetResourceRecords(ctx context.Context, domainName string) ([]ResourceRecord, error) {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	var answer zoneDomainsAnswer
	if err := s.client.call(ctx, "zone", "get_resource_records", params, &answer); err != nil {
		return nil, err
	}
	if err := checkZoneDomainResults(answer.Domains, domainName); err != nil {
		return nil, err
	}
	for _, d := range answer.Domains {
		if d.Dname == domainName || len(answer.Domains) == 1 {
			return d.RRs, nil
		}
	}
	return nil, nil
}

// GetRecords is an alias for GetResourceRecords.
func (s *ZoneService) GetRecords(ctx context.Context, domainName string) ([]ResourceRecord, error) {
	return s.GetResourceRecords(ctx, domainName)
}

// FilterRecords returns records matching all non-empty filter fields.
func FilterRecords(records []ResourceRecord, f RecordFilter) []ResourceRecord {
	if f.Type == "" && f.Subdomain == "" && f.Content == "" {
		return records
	}
	out := make([]ResourceRecord, 0, len(records))
	for _, r := range records {
		if f.Type != "" && !strings.EqualFold(r.Rectype, f.Type) {
			continue
		}
		if f.Subdomain != "" && r.Subname != f.Subdomain {
			continue
		}
		if f.Content != "" && r.Content != f.Content {
			continue
		}
		out = append(out, r)
	}
	return out
}

// AddRecord adds a DNS record using the appropriate zone/add_* method.
func (s *ZoneService) AddRecord(ctx context.Context, domainName string, record DNSRecord) error {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	if record.SubDomain == "" {
		record.SubDomain = "@"
	}

	typ := strings.ToUpper(strings.TrimSpace(record.RecordType))
	var fn string

	switch typ {
	case "A", "ALIAS":
		if record.Content == "" {
			return fmt.Errorf("%s: укажите content (IPv4)", typ)
		}
		fn = "add_alias"
		params["subdomain"] = record.SubDomain
		params["ipaddr"] = record.Content
	case "AAAA":
		if record.Content == "" {
			return fmt.Errorf("AAAA: укажите content (IPv6)")
		}
		fn = "add_aaaa"
		params["subdomain"] = record.SubDomain
		params["ipaddr"] = record.Content
	case "CNAME":
		fn = "add_cname"
		params["subdomain"] = record.SubDomain
		cname := record.CanonicalName
		if cname == "" {
			cname = record.Content
		}
		if cname == "" {
			return fmt.Errorf("CNAME: укажите canonical_name или content")
		}
		params["canonical_name"] = cname
	case "MX":
		fn = "add_mx"
		params["subdomain"] = record.SubDomain
		mail := record.Content
		if mail == "" {
			return fmt.Errorf("MX: укажите content (mail_server)")
		}
		params["mail_server"] = mail
		if record.Priority != "" {
			params["priority"] = record.Priority
		}
	case "NS":
		fn = "add_ns"
		params["subdomain"] = record.SubDomain
		ns := record.DNSServer
		if ns == "" {
			ns = record.Content
		}
		if ns == "" {
			return fmt.Errorf("NS: укажите dns_server или content")
		}
		params["dns_server"] = ns
		if record.RecordNumber != "" {
			params["record_number"] = record.RecordNumber
		} else if record.Priority != "" {
			params["record_number"] = record.Priority
		}
	case "TXT":
		fn = "add_txt"
		params["subdomain"] = record.SubDomain
		text := record.Text
		if text == "" {
			text = record.Content
		}
		if text == "" {
			return fmt.Errorf("TXT: укажите text или content")
		}
		params["text"] = text
	case "SRV":
		fn = "add_srv"
		if record.Service == "" {
			return fmt.Errorf("SRV: укажите service (например _sip._udp)")
		}
		params["service"] = record.Service
		if record.Target != "" {
			params["target"] = record.Target
		} else if record.Content != "" {
			params["target"] = record.Content
		} else {
			return fmt.Errorf("SRV: укажите target или content")
		}
		if record.Port != "" {
			params["port"] = record.Port
		}
		if record.Priority != "" {
			params["priority"] = record.Priority
		}
		if record.Weight != "" {
			params["weight"] = record.Weight
		}
	case "CAA":
		fn = "add_caa"
		params["subdomain"] = record.SubDomain
		params["flags"] = record.Flags
		if record.Tag == "" {
			return fmt.Errorf("CAA: укажите tag (issue, issuewild, iodef)")
		}
		params["tag"] = record.Tag
		val := record.Value
		if val == "" {
			val = record.Content
		}
		if val == "" {
			return fmt.Errorf("CAA: укажите value или content")
		}
		params["value"] = val
	case "HTTPS":
		fn = "add_https"
		params["subdomain"] = record.SubDomain
		if record.Target == "" {
			return fmt.Errorf("HTTPS: укажите target")
		}
		params["target"] = record.Target
		if record.Priority != "" {
			params["priority"] = record.Priority
		}
		if record.Value != "" {
			params["value"] = record.Value
		}
	default:
		return fmt.Errorf("неподдерживаемый тип записи %q (A, AAAA, CNAME, MX, NS, TXT, SRV, CAA, HTTPS)", typ)
	}

	return s.zoneCall(ctx, fn, params)
}

// RemoveRecord deletes a DNS record (zone/remove_record).
func (s *ZoneService) RemoveRecord(ctx context.Context, domainName, subdomain, recordType, content, priority string) error {
	if subdomain == "" {
		return fmt.Errorf("subdomain обязателен")
	}
	if recordType == "" {
		return fmt.Errorf("record_type обязателен")
	}
	params := map[string]interface{}{
		"domain_name": domainName,
		"subdomain":   subdomain,
		"record_type": strings.ToUpper(recordType),
	}
	if content != "" {
		params["content"] = content
	}
	if priority != "" {
		params["priority"] = priority
	}
	return s.zoneCall(ctx, "remove_record", params)
}

// DeleteRecord is an alias for RemoveRecord (без record_id — API REG.RU его не возвращает).
func (s *ZoneService) DeleteRecord(ctx context.Context, domainName, subdomain, recordType, content, priority string) error {
	return s.RemoveRecord(ctx, domainName, subdomain, recordType, content, priority)
}

// ClearZone removes all DNS records from a domain's zone.
func (s *ZoneService) ClearZone(ctx context.Context, domainName string) error {
	params := map[string]interface{}{
		"domain_name": domainName,
	}
	return s.zoneCall(ctx, "clear", params)
}

// Add_Alias adds an A record (zone/add_alias).
func (s *ZoneService) Add_Alias(ctx context.Context, domainName, subdomain, ip string) error {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Content:    ip,
		RecordType: "A",
	})
}

// Add_AAAA adds an AAAA record (zone/add_aaaa).
func (s *ZoneService) Add_AAAA(ctx context.Context, domainName, subdomain, ipv6 string) error {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Content:    ipv6,
		RecordType: "AAAA",
	})
}

// Add_CNAME adds a CNAME record.
func (s *ZoneService) Add_CNAME(ctx context.Context, domainName, subdomain, target string) error {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:     subdomain,
		CanonicalName: target,
		RecordType:    "CNAME",
	})
}

// Add_MX adds an MX record.
func (s *ZoneService) Add_MX(ctx context.Context, domainName, subdomain, mailserver, priority string) error {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Content:    mailserver,
		RecordType: "MX",
		Priority:   priority,
	})
}

// Add_TXT adds a TXT record.
func (s *ZoneService) Add_TXT(ctx context.Context, domainName, subdomain, text string) error {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:  subdomain,
		Text:       text,
		RecordType: "TXT",
	})
}

// Add_NS adds an NS record.
func (s *ZoneService) Add_NS(ctx context.Context, domainName, subdomain, nameserver, recordNumber string) error {
	return s.AddRecord(ctx, domainName, DNSRecord{
		SubDomain:    subdomain,
		DNSServer:    nameserver,
		RecordNumber: recordNumber,
		RecordType:   "NS",
	})
}

// Add_SRV adds an SRV record.
func (s *ZoneService) Add_SRV(ctx context.Context, domainName, service, target, port, priority, weight string) error {
	return s.AddRecord(ctx, domainName, DNSRecord{
		Service:    service,
		Target:     target,
		Port:       port,
		Priority:   priority,
		Weight:     weight,
		RecordType: "SRV",
	})
}
