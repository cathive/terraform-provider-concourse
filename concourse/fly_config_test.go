package concourse

import (
	"testing"
	"os"
)

func TestFlyRC_ImportConfigFromEnv(t *testing.T) {
	test_file := "./test_fixtures/flyrc.yml"
	os.Setenv("FLYRC", test_file)

	var flyrc FlyRc

	err := flyrc.ImportConfig()
	if err != nil {
		t.Fatalf("failed to import config from %s: %s", test_file, err)
	}

	if flyrc.Filename != test_file {
		t.Fatalf("filename stored in FlyRc (%) is not the same as the test file (%s)", flyrc.Filename, test_file)
	}

}
