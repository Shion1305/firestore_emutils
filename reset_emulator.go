package emutils

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
)

var regexHost, regexProjID regexp.Regexp

func init() {
	regexHost = *regexp.MustCompile(`^[^:]+:\d{1,5}$`)
	regexProjID = *regexp.MustCompile(`^[a-z][a-z0-9-]{4,28}[a-z0-9]$`)
}

func ResetEmulator(projectID string) error {
	if !regexProjID.MatchString(projectID) {
		return fmt.Errorf("project ID must match %s, %s", regexProjID.String(), projectID)
	}
	emuHost, err := loadEmulatorHostEnv()
	if err != nil {
		return err
	}
	return execResetReq(emuHost, projectID)
}

func loadEmulatorHostEnv() (host string, err error) {
	emuHost := os.Getenv("FIRESTORE_EMULATOR_HOST")
	emuPort := os.Getenv("FIRESTORE_EMULATOR_PORT")
	if emuHost == "" {
		if emuPort == "" {
			// if all variables are not set, throw error
			return "", fmt.Errorf("missing FIRESTORE_EMULATOR_HOST")
		}
		return fmt.Sprintf("localhost:%s", emuPort), nil
	}
	// validate if Host is in the format of host:port
	if regexHost.MatchString(emuHost) {
		return emuHost, nil
	}
	if emuPort == "" {
		return "", fmt.Errorf("invalid FIRESTORE_EMULATOR_HOST, should be host:port but got: %s", emuHost)
	}
	return fmt.Sprintf("%s:%s", emuHost, emuPort), nil
}

func execResetReq(emuHost, projectID string) error {
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
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
