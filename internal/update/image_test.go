package update

import (
	"os"
	"regexp"
	"testing"
)

// The updater still names the upstream's versioned image rather than a
// floating tag. This fork does not publish Docker images: its workflow only
// builds and smoke-tests an image, so it must not be required to push tags.
func TestImagePointsToUpstreamAndForkDockerIsSmokeOnly(t *testing.T) {
	for in, want := range map[string]string{
		"0.1.4479":  "ghcr.io/yetone/magpie:0.1.4479",
		"v0.1.4479": "ghcr.io/yetone/magpie:0.1.4479",
	} {
		if got := Image(in); got != want {
			t.Errorf("Image(%q) = %q, want %q", in, got, want)
		}
	}
	if ImageRepo != "ghcr.io/yetone/magpie" {
		t.Errorf("ImageRepo = %q, want upstream image", ImageRepo)
	}

	y, err := os.ReadFile("../../.github/workflows/docker.yml")
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`docker build -t magpie-smoke`).Match(y) ||
		!regexp.MustCompile(`docker run -d --name smoke`).Match(y) {
		t.Error("fork docker.yml must build and run a smoke-test image")
	}
	if regexp.MustCompile(`docker (login|push)`).Match(y) ||
		regexp.MustCompile(`IMAGE: ghcr\.io/`).Match(y) {
		t.Error("fork docker.yml must not log in to or publish an image")
	}
}
