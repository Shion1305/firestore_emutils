package emutils

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"cloud.google.com/go/firestore"

	"github.com/stretchr/testify/require"
)

func TestEmulatorInfo_ClearData(t *testing.T) {
	host, port, err := setupEnv()
	require.NoError(t, err)

	tests := []struct {
		prerequires func() error
		name        string
		exec        func() error
		requires    func(t *testing.T) error
	}{
		{
			prerequires: func() error {
				ctx := context.Background()
				client, err := firestore.NewClient(ctx, "test")
				if err != nil {
					return err
				}
				colsRef := client.Collection("emulator")
				_, _, err = colsRef.Add(ctx, map[string]interface{}{
					"tests": "test",
				})
				return err
			},
			exec: func() error { return nil },
			requires: func(t *testing.T) error {
				ctx := context.Background()
				client, err := firestore.NewClient(ctx, "test")
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
				client, err := firestore.NewClient(ctx, "test")
				if err != nil {
					return err
				}
				colsRef := client.Collection("emulator")
				_, _, err = colsRef.Add(ctx, map[string]interface{}{
					"tests": "test",
				})
				return err
			},
			exec: func() error {
				emu := NewEmulator(host, port, "test")
				return emu.ClearData()
			},
			requires: func(t *testing.T) error {
				ctx := context.Background()
				client, err := firestore.NewClient(ctx, "test")
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

func setupEnv() (host string, port int, err error) {
	host = os.Getenv("EMULATOR_HOST")
	portStr := os.Getenv("EMULATOR_PORT")
	port64, err := strconv.ParseInt(portStr, 10, 32)
	hostParam := fmt.Sprintf("%s:%d", host, port)
	if err := os.Setenv("EMULATOR_HOST", hostParam); err != nil {
		return host, port, err
	}
	return host, int(port64), err
}
