package integration_test

import (
	"os"
	"testing"
)

func TestIntegrationRealNetworkIsOptIn(t *testing.T) {
	requireIntegration(t, "SEEKMOON_INTEGRATION_NETWORK")
	t.Skip("real network integration implementation is intentionally outside default WP13 acceptance")
}

func TestIntegrationMoonCLIProbeMutationIsOptIn(t *testing.T) {
	requireIntegration(t, "SEEKMOON_INTEGRATION_MOONCLI")
	t.Skip("real Moon CLI probe mutation is intentionally outside default WP13 acceptance")
}

func TestIntegrationGitHubAPIIsOptIn(t *testing.T) {
	requireIntegration(t, "SEEKMOON_INTEGRATION_GITHUB")
	t.Skip("real GitHub API integration is intentionally outside default WP13 acceptance")
}

func requireIntegration(t *testing.T, env string) {
	t.Helper()
	if os.Getenv(env) == "" {
		t.Skipf("%s is not set; integration test skipped by default", env)
	}
}
