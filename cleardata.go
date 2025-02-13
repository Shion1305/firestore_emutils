package emutils

import (
	"fmt"
	"net/http"
)

type EmulatorInfo struct {
	Host string
	Port int
}

func NewEmulator(
	Host string,
	Port int,
) EmulatorInfo {
	return EmulatorInfo{
		Host, Port,
	}
}

// ClearData sends DELETE request to emulator host
func (info EmulatorInfo) ClearData(projectID string) error {
	url := fmt.Sprintf("http://%s:%d/emulator/v1/projects/%s/databases/(default)/documents",
		info.Host, info.Port, projectID)
	req, err := http.NewRequest(
		"DELETE",
		url,
		nil,
	)
	if err != nil {
		return fmt.Errorf("faield to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("faield to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
