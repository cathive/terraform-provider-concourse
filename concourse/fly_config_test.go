package concourse

import (
	"os"
	"testing"
)

func TestFlyRC_ImportConfigFromEnv(t *testing.T) {
	expected_test_file := "./testdata/flyrc.yml"
	expected_target_count := 2

	os.Setenv("FLYRC", expected_test_file)

	var flyrc FlyRc

	err := flyrc.ImportConfig()
	if err != nil {
		t.Fatalf("failed to import config from %s: %s", expected_test_file, err)
	}

	if flyrc.Filename != expected_test_file {
		t.Fatalf("filename stored in FlyRc (%) is not the same as the test file (%s)", flyrc.Filename, expected_test_file)
	}

	if len(flyrc.Targets) != expected_target_count {
		t.Fatalf("expected %d targets, but counted %d", expected_target_count, len(flyrc.Targets))
	}

}
