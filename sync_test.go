package regru

import "testing"

func TestPlanZoneSync_addAndPrune(t *testing.T) {
	desired := []DNSRecord{
		{RecordType: "A", SubDomain: "www", Content: "1.2.3.4"},
		{RecordType: "MX", SubDomain: "@", Content: "mail.example.ru", Priority: "10"},
	}
	current := []ResourceRecord{
		{Rectype: "A", Subname: "www", Content: "1.2.3.4"},
		{Rectype: "A", Subname: "old", Content: "9.9.9.9"},
	}

	plan := PlanZoneSync(desired, current, true)
	if len(plan.ToAdd) != 1 {
		t.Fatalf("ToAdd: want 1, got %d (%+v)", len(plan.ToAdd), plan.ToAdd)
	}
	if plan.ToAdd[0].RecordType != "MX" {
		t.Fatalf("ToAdd type: want MX, got %s", plan.ToAdd[0].RecordType)
	}
	if len(plan.ToRemove) != 1 {
		t.Fatalf("ToRemove: want 1, got %d (%+v)", len(plan.ToRemove), plan.ToRemove)
	}
	if plan.ToRemove[0].Subname != "old" {
		t.Fatalf("ToRemove subdomain: want old, got %s", plan.ToRemove[0].Subname)
	}
}

func TestPlanZoneSync_noPrune(t *testing.T) {
	desired := []DNSRecord{
		{RecordType: "A", SubDomain: "www", Content: "1.2.3.4"},
	}
	current := []ResourceRecord{
		{Rectype: "A", Subname: "old", Content: "9.9.9.9"},
	}

	plan := PlanZoneSync(desired, current, false)
	if len(plan.ToAdd) != 1 {
		t.Fatalf("ToAdd: want 1, got %d", len(plan.ToAdd))
	}
	if len(plan.ToRemove) != 0 {
		t.Fatalf("ToRemove: want 0 without prune, got %d", len(plan.ToRemove))
	}
}

func TestPlanZoneSync_unchanged(t *testing.T) {
	desired := []DNSRecord{
		{RecordType: "TXT", SubDomain: "@", Text: "v=spf1 ~all"},
	}
	current := []ResourceRecord{
		{Rectype: "TXT", Subname: "@", Content: "v=spf1 ~all"},
	}

	plan := PlanZoneSync(desired, current, true)
	if len(plan.ToAdd) != 0 || len(plan.ToRemove) != 0 {
		t.Fatalf("want empty plan, got add=%d remove=%d", len(plan.ToAdd), len(plan.ToRemove))
	}
}

func TestDesiredRecordKey_aliases(t *testing.T) {
	cases := []struct {
		name string
		d    DNSRecord
		want string
	}{
		{
			name: "cname via canonical_name",
			d:    DNSRecord{RecordType: "CNAME", SubDomain: "mail", CanonicalName: "mx.example.ru"},
			want: RecordSyncKey("CNAME", "mail", "mx.example.ru", ""),
		},
		{
			name: "empty subdomain becomes @",
			d:    DNSRecord{RecordType: "A", Content: "1.2.3.4"},
			want: RecordSyncKey("A", "@", "1.2.3.4", ""),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := DesiredRecordKey(tc.d)
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
