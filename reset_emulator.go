package emutils

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
)

const regexHost = `[^:]+:\d{1,5}`

func ResetEmulator(projectID string) error {
	emuHost := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if emuHost == "" {
		return fmt.Errorf("missing FIRESTORE_EMULATOR_HOST")
	}
	re := regexp.MustCompile(regexHost)
	if !re.MatchString(emuHost) {
		return fmt.Errorf("invalid FIRESTORE_EMULATOR_HOST, should be host:port but got: %s", emuHost)
	}

	url := fmt.Sprintf(
		"http://%s/emulator/v1/projects/%s/databases/(default)/documents",
		emuHost, projectID,
	)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
