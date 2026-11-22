package update

import "strings"

// ImageRepo is the upstream Docker image registry. This fork's docker.yml
// only smoke-tests an image and does not publish to this registry.
const ImageRepo = "ghcr.io/yetone/magpie"

// Image formats an upstream image reference using an exact version tag
// rather than a floating tag. It does not imply that this fork publishes
// images for its own release versions.
func Image(version string) string {
	return ImageRepo + ":" + strings.TrimPrefix(version, "v")
}
