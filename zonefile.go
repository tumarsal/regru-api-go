package regru

import (
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

// ZoneFile is the YAML document used by zone sync (domain + records).
type ZoneFile struct {
	Domain  string           `yaml:"domain,omitempty"`
	Records []ZoneFileRecord `yaml:"records"`
}

// ZoneFileRecord is one DNS record in a zone YAML file.
type ZoneFileRecord struct {
	Type          string `yaml:"type"`
	Subdomain     string `yaml:"subdomain,omitempty"`
	Content       string `yaml:"content,omitempty"`
	Priority      string `yaml:"priority,omitempty"`
	CanonicalName string `yaml:"canonical_name,omitempty"`
	DNSServer     string `yaml:"dns_server,omitempty"`
	RecordNumber  string `yaml:"record_number,omitempty"`
	Service       string `yaml:"service,omitempty"`
	Target        string `yaml:"target,omitempty"`
	Port          string `yaml:"port,omitempty"`
	Weight        string `yaml:"weight,omitempty"`
	Text          string `yaml:"text,omitempty"`
	Flags         int    `yaml:"flags,omitempty"`
	Tag           string `yaml:"tag,omitempty"`
	Value         string `yaml:"value,omitempty"`
}

// LoadZoneFile reads and parses a zone YAML file.
func LoadZoneFile(path string) (*ZoneFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var spec ZoneFile
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("разбор YAML: %w", err)
	}
	return &spec, nil
}

// WriteZoneFile marshals and writes a zone YAML file.
func WriteZoneFile(path string, spec *ZoneFile) error {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("формирование YAML: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("запись файла: %w", err)
	}
	return nil
}

// MarshalZoneFile returns YAML bytes for a zone file.
func MarshalZoneFile(spec *ZoneFile) ([]byte, error) {
	return yaml.Marshal(spec)
}

// ZoneFileFromRecords builds a ZoneFile from API resource records.
func ZoneFileFromRecords(domain string, records []ResourceRecord) *ZoneFile {
	return &ZoneFile{
		Domain:  domain,
		Records: zoneFileRecordsFromAPI(records),
	}
}

func zoneFileRecordsFromAPI(records []ResourceRecord) []ZoneFileRecord {
	out := make([]ZoneFileRecord, 0, len(records))
	for _, r := range records {
		row := ZoneFileRecord{
			Type:      r.Rectype,
			Subdomain: r.Subname,
			Content:   r.Content,
		}
		if p := r.Priority.String(); p != "" {
			row.Priority = p
		}
		out = append(out, row)
	}
	return out
}

// DNSRecordsFromZoneFile converts YAML records to DNSRecord for AddRecord / PlanZoneSync.
func DNSRecordsFromZoneFile(rows []ZoneFileRecord) ([]DNSRecord, error) {
	out := make([]DNSRecord, 0, len(rows))
	for i, row := range rows {
		if strings.TrimSpace(row.Type) == "" {
			return nil, fmt.Errorf("запись %d: не указан type", i+1)
		}
		rec := DNSRecord{
			RecordType:    row.Type,
			SubDomain:     row.Subdomain,
			Content:       row.Content,
			Priority:      row.Priority,
			CanonicalName: row.CanonicalName,
			DNSServer:     row.DNSServer,
			RecordNumber:  row.RecordNumber,
			Service:       row.Service,
			Target:        row.Target,
			Port:          row.Port,
			Weight:        row.Weight,
			Text:          row.Text,
			Flags:         row.Flags,
			Tag:           row.Tag,
			Value:         row.Value,
		}
		if strings.ToUpper(strings.TrimSpace(rec.RecordType)) == "HTTPS" && rec.Priority == "" {
			rec.Priority = "0"
		}
		out = append(out, rec)
	}
	return out, nil
}
