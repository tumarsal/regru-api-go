package regru

import "strings"

// RecordSyncKey возвращает ключ записи в формате, совпадающем с колонками zone list
// (type, subdomain, content, priority).
func RecordSyncKey(recordType, subdomain, content, priority string) string {
	typ := strings.ToUpper(strings.TrimSpace(recordType))
	sub := strings.TrimSpace(subdomain)
	if sub == "" {
		sub = "@"
	}
	return typ + "\x00" + sub + "\x00" + content + "\x00" + priority
}

// ResourceRecordKey — ключ существующей записи из API.
func ResourceRecordKey(r ResourceRecord) string {
	return RecordSyncKey(r.Rectype, r.Subname, r.Content, r.Priority.String())
}

// DesiredRecordKey — ключ желаемой записи для добавления.
func DesiredRecordKey(d DNSRecord) string {
	sub := d.SubDomain
	content := d.Content
	priority := d.Priority

	typ := strings.ToUpper(strings.TrimSpace(d.RecordType))
	switch typ {
	case "CNAME":
		if d.CanonicalName != "" {
			content = d.CanonicalName
		}
	case "NS":
		if d.DNSServer != "" {
			content = d.DNSServer
		}
		if d.RecordNumber != "" {
			priority = d.RecordNumber
		}
	case "TXT":
		if d.Text != "" {
			content = d.Text
		}
	case "SRV":
		sub = d.Service
		if sub == "" {
			sub = d.SubDomain
		}
		// content в API для SRV может отличаться; для сравнения используем target/port/priority.
		if d.Target != "" {
			content = d.Target
		}
	case "CAA":
		if d.Tag != "" && d.Value != "" {
			content = d.Tag + " " + d.Value
		} else if d.Value != "" {
			content = d.Value
		}
	case "HTTPS":
		if d.Target != "" {
			content = d.Target
		}
	}

	return RecordSyncKey(d.RecordType, sub, content, priority)
}

// ZoneSyncPlan — записи для добавления и удаления при синхронизации зоны.
type ZoneSyncPlan struct {
	ToAdd    []DNSRecord
	ToRemove []ResourceRecord
}

// PlanZoneSync сравнивает желаемое состояние с текущим (полная замена при prune).
func PlanZoneSync(desired []DNSRecord, current []ResourceRecord, prune bool) ZoneSyncPlan {
	desiredKeys := make(map[string]DNSRecord, len(desired))
	for _, d := range desired {
		desiredKeys[DesiredRecordKey(d)] = d
	}

	currentKeys := make(map[string]ResourceRecord, len(current))
	for _, r := range current {
		currentKeys[ResourceRecordKey(r)] = r
	}

	var plan ZoneSyncPlan
	for k, d := range desiredKeys {
		if _, ok := currentKeys[k]; !ok {
			plan.ToAdd = append(plan.ToAdd, d)
		}
	}
	if prune {
		for k, r := range currentKeys {
			if _, ok := desiredKeys[k]; !ok {
				plan.ToRemove = append(plan.ToRemove, r)
			}
		}
	}
	return plan
}
