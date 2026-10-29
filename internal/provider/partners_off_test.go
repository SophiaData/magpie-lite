package provider

import (
	"context"
	"testing"
	"time"
)

// This build lists no partners, and that is not a setting: naming a list in
// the environment must not make one appear, fetch, or count.
func TestPartnersStayOffWhateverTheEnvironmentSays(t *testing.T) {
	t.Setenv("MAGPIE_PARTNERS", "https://example.test/api/partners")
	resetPartners()
	t.Cleanup(resetPartners)

	if got := partnersSource(); got != "off" {
		t.Fatalf("partnersSource = %q, want off", got)
	}
	if got := Partners(); len(got) != 0 {
		t.Fatalf("Partners = %v", got)
	}
	if got := PartnersNow(0); len(got) != 0 {
		t.Fatalf("PartnersNow = %v", got)
	}
	if PartnerNoticed("acme") {
		t.Fatal("a partner was noticed")
	}
	if got := PartnerCounts(time.Now()); len(got) != 0 {
		t.Fatalf("PartnerCounts = %v", got)
	}
	if got := keptPartnerIcons(); len(got) != 0 {
		t.Fatalf("keptPartnerIcons = %v", got)
	}
	if p := partnerPreset("acme"); p != nil {
		t.Fatalf("partnerPreset = %+v", p)
	}

	// The kept-fresh loop returns at once rather than fetching.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	go func() { defer close(done); KeepPartnersFresh(ctx) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("KeepPartnersFresh kept running with no partners to fetch")
	}
}
