package gui

import (
	"testing"
	"time"

	"github.com/yetone/magpie/internal/settings"
)

// This build never asks the feed by itself (autoUpdatesDisabled): however
// the settings are, nothing is due. A check is the user's own (Settings'
// Check, `magpie update`).
func TestUpdateDue(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	now := time.Now()
	s := settings.Load()
	for _, c := range []struct {
		off   bool
		every int
		ago   time.Duration
	}{
		{false, 0, 0},
		{false, 360, 5*time.Hour + 59*time.Minute},
		{false, 360, 6 * time.Hour},
		{false, 30, 30 * time.Minute},
		{false, 1440, 24 * time.Hour},
		{true, 30, 48 * time.Hour},
	} {
		s.NoAutoUpdate, s.UpdateEvery = c.off, c.every
		if got := updateDue(s, now.Add(-c.ago), now); got {
			t.Errorf("off %v, every %d, %v ago: due %v, want never due", c.off, c.every, c.ago, got)
		}
	}
	if autoUpdatesDisabled != true {
		t.Fatal("automatic updates are not disabled in this build")
	}
	// never asked: not due either
	s.NoAutoUpdate, s.UpdateEvery = false, 1440
	if updateDue(s, time.Time{}, now) {
		t.Error("never asked: due")
	}
	// only the presets are kept
	s.NoAutoUpdate, s.UpdateEvery = false, 120
	if err := settings.Save(s); err == nil {
		t.Error("every 120 minutes saved")
	}
}

// A check, asked for or by the clock, is when the next one counts from.
func TestUpdateAskedAt(t *testing.T) {
	u := &updater{}
	before := time.Now()
	if !u.begin() {
		t.Fatal("not begun")
	}
	if u.asked.Before(before) {
		t.Fatalf("asked at %v, before %v", u.asked, before)
	}
}
