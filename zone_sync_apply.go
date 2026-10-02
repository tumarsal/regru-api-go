package regru

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

// SyncOptions controls SyncZone behaviour.
type SyncOptions struct {
	DryRun bool
	Prune  bool // when true (default for callers), remove records absent from the file
}

// SyncChangeKind is one planned or applied change.
type SyncChangeKind string

const (
	SyncChangeAdd    SyncChangeKind = "add"
	SyncChangeRemove SyncChangeKind = "remove"
)

// SyncChange is a human-readable line for a single DNS change.
type SyncChange struct {
	Kind SyncChangeKind
	Line string
}

// SyncResult summarizes a SyncZone run.
type SyncResult struct {
	// Exported is true when the YAML file did not exist and was created (or would be, in dry-run).
	Exported bool
	// ExportYAML is set on dry-run export when the file was not written.
	ExportYAML string
	// Message is a short summary for the user.
	Message string
	// Changes lists +/- lines (empty when Exported or no diff).
	Changes []SyncChange
	Added   int
	Removed int
}

// SyncZone synchronizes a domain zone with a YAML file.
// If path does not exist, current records are exported to that file (no sync).
func SyncZone(ctx context.Context, client *Client, domain, path string, opts SyncOptions) (*SyncResult, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	domain = strings.TrimSpace(domain)
	path = strings.TrimSpace(path)
	if domain == "" {
		return nil, fmt.Errorf("domain is required")
	}
	if path == "" {
		return nil, fmt.Errorf("filepath is required")
	}

	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return exportZone(ctx, client, domain, path, opts.DryRun)
		}
		return nil, fmt.Errorf("чтение файла: %w", err)
	}

	spec, err := LoadZoneFile(path)
	if err != nil {
		return nil, err
	}
	if len(spec.Records) == 0 {
		return nil, fmt.Errorf("в файле нет записей (ключ records)")
	}
	if spec.Domain != "" && spec.Domain != domain {
		return nil, fmt.Errorf("домен в аргументе (%s) не совпадает с domain в файле (%s)", domain, spec.Domain)
	}

	desired, err := DNSRecordsFromZoneFile(spec.Records)
	if err != nil {
		return nil, err
	}

	current, err := client.Zone().GetResourceRecords(ctx, domain)
	if err != nil {
		return nil, err
	}

	plan := PlanZoneSync(desired, current, opts.Prune)
	result := &SyncResult{
		Added:   len(plan.ToAdd),
		Removed: len(plan.ToRemove),
	}

	if len(plan.ToAdd) == 0 && len(plan.ToRemove) == 0 {
		result.Message = "изменений нет"
		return result, nil
	}

	for _, r := range plan.ToRemove {
		line := fmt.Sprintf("- %s %s %s", r.Rectype, r.Subname, r.Content)
		if p := r.Priority.String(); p != "" {
			line += " (prio " + p + ")"
		}
		result.Changes = append(result.Changes, SyncChange{Kind: SyncChangeRemove, Line: line})
	}
	for _, d := range plan.ToAdd {
		key := DesiredRecordKey(d)
		parts := strings.Split(key, "\x00")
		line := fmt.Sprintf("+ %s %s %s", parts[0], parts[1], parts[2])
		if len(parts) > 3 && parts[3] != "" {
			line += " (prio " + parts[3] + ")"
		}
		result.Changes = append(result.Changes, SyncChange{Kind: SyncChangeAdd, Line: line})
	}

	if opts.DryRun {
		result.Message = "(dry-run: изменения не применены)"
		return result, nil
	}

	for _, r := range plan.ToRemove {
		if err := client.Zone().RemoveRecord(ctx, domain, r.Subname, r.Rectype, r.Content, r.Priority.String()); err != nil {
			return result, fmt.Errorf("удаление %s %s: %w", r.Rectype, r.Subname, err)
		}
	}
	for _, d := range plan.ToAdd {
		if err := client.Zone().AddRecord(ctx, domain, d); err != nil {
			return result, fmt.Errorf("добавление %s %s: %w", d.RecordType, d.SubDomain, err)
		}
	}

	result.Message = fmt.Sprintf("синхронизация %s завершена: удалено %d, добавлено %d", domain, len(plan.ToRemove), len(plan.ToAdd))
	return result, nil
}

func exportZone(ctx context.Context, client *Client, domain, path string, dryRun bool) (*SyncResult, error) {
	records, err := client.Zone().GetResourceRecords(ctx, domain)
	if err != nil {
		return nil, err
	}
	spec := ZoneFileFromRecords(domain, records)
	data, err := MarshalZoneFile(spec)
	if err != nil {
		return nil, fmt.Errorf("формирование YAML: %w", err)
	}

	result := &SyncResult{Exported: true}
	if dryRun {
		result.ExportYAML = string(data)
		result.Message = fmt.Sprintf("(dry-run: файл %s не записан)", path)
		return result, nil
	}
	if err := WriteZoneFile(path, spec); err != nil {
		return nil, err
	}
	result.Message = fmt.Sprintf("файл %s создан: %d записей из API", path, len(spec.Records))
	return result, nil
}
