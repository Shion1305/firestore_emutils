package emutils

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"

	"github.com/stretchr/testify/require"
)

func TestEmulatorInfo_ClearData(t *testing.T) {
	host, port, err := setupEnv(t)
	require.NoError(t, err)

	projectID := "test"

	tests := []struct {
		prerequires func() error
		name        string
		exec        func() error
		requires    func(t *testing.T) error
	}{
		{
			prerequires: func() error {
				ctx := context.Background()
				client, err := firestore.NewClient(ctx, projectID)
				if err != nil {
					return err
				}
				colsRef := client.Collection("emulator")
				_, _, err = colsRef.Add(ctx, map[string]interface{}{
					"Shion": "Ichikawa",
				})
				return err
			},
			exec: func() error { return nil },
			requires: func(t *testing.T) error {
				ctx := context.Background()
				client, err := firestore.NewClient(ctx, projectID)
				if err != nil {
					return fmt.Errorf("could not create Firestore client: %v", err)
				}
				colsRef := client.Collections(ctx)
				cols, err := colsRef.GetAll()
				if err != nil {
					return fmt.Errorf("could not get collections: %v", err)
				}
				require.Len(t, cols, 1)
				return nil
			},
		},
		{
			prerequires: func() error {
				ctx := context.Background()
				client, err := firestore.NewClient(ctx, projectID)
				if err != nil {
					return err
				}
				colsRef := client.Collection("emulator")
				_, _, err = colsRef.Add(ctx, map[string]interface{}{
					"Shion": "Ichikawa",
				})
				return err
			},
			exec: func() error {
				emu := NewEmulator(host, port)
				return emu.ClearData(projectID)
			},
			requires: func(t *testing.T) error {
				ctx := context.Background()
				client, err := firestore.NewClient(ctx, projectID)
				if err != nil {
					return fmt.Errorf("could not create Firestore client: %v", err)
				}
				colsRef := client.Collections(ctx)
				cols, err := colsRef.GetAll()
				if err != nil {
					return fmt.Errorf("could not get collections: %v", err)
				}
				require.Len(t, cols, 0)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, tt.prerequires(), "failed to satisfy prerequisites")
			require.NoError(t, tt.exec(), "error in execution")
			require.NoError(t, tt.requires(t), "failed to require emulator")
		})
	}
}

func setupEnv(t *testing.T) (host string, port int, err error) {
	hostRaw := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if hostRaw == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST not set, skipping test")
	}
	parts := strings.Split(hostRaw, ":")
	require.Len(t, parts, 2, "FIRESTORE_EMULATOR_HOST should be in the format host:port")
	host = parts[0]
	port, err = strconv.Atoi(parts[1])
	require.NoError(t, err, "FIRESTORE_EMULATOR_HOST port should be an integer")
	return host, port, err
}
