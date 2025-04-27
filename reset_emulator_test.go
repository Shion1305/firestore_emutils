package emutils_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	emutils "github.com/Shion1305/firestore_emutils"
)

func TestResetEmulatorIntegration(t *testing.T) {
	defaultValue := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if defaultValue == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST not set; skipping integration test.")
	}
	defer func() {
		_ = os.Setenv("FIRESTORE_EMULATOR_HOST", defaultValue)
		_ = os.Unsetenv("FIRESTORE_EMULATOR_PORT")
	}()
	tests := []struct {
		name      string
		projectID string
		setup     func(t *testing.T)
		reset     func() error
		wantErr   bool
	}{
		{
			name:      "Missing FIRESTORE_EMULATOR_HOST",
			projectID: "dummy-project",
			setup: func(t *testing.T) {
				_ = os.Unsetenv("FIRESTORE_EMULATOR_HOST")
			},
			reset: func() error {
				return nil
			},
			wantErr: true,
		},
		{
			name:      "Valid Reset with FIRESTORE_EMULATOR_HOST",
			projectID: "test-project",
			setup: func(t *testing.T) {
				_ = os.Setenv("FIRESTORE_EMULATOR_HOST", defaultValue)
			},
			reset: func() error {
				return os.Unsetenv("FIRESTORE_EMULATOR_HOST")
			},
			wantErr: false,
		},
		{
			name:      "Valid Reset with FIRESTORE_EMULATOR_PORT",
			projectID: "test-project",
			setup: func(t *testing.T) {
				splits := strings.Split(defaultValue, ":")
				require.Len(t, splits, 2)
				p := splits[1]
				_ = os.Unsetenv("FIRESTORE_EMULATOR_HOST")
				_ = os.Setenv("FIRESTORE_EMULATOR_PORT", p)
			},
			reset: func() error {
				return os.Unsetenv("FIRESTORE_EMULATOR_PORT")
			},
			wantErr: false,
		},
		{
			name:      "Valid Reset with FIRESTORE_EMULATOR_HOST,FIRESTORE_EMULATOR_PORT",
			projectID: "test-project",
			setup: func(t *testing.T) {
				splits := strings.Split(defaultValue, ":")
				require.Len(t, splits, 2)
				h := splits[0]
				p := splits[1]
				_ = os.Setenv("FIRESTORE_EMULATOR_HOST", h)
				_ = os.Setenv("FIRESTORE_EMULATOR_PORT", p)
			},
			reset: func() error {
				if err := os.Unsetenv("FIRESTORE_EMULATOR_HOST"); err != nil {
					return err
				}
				return os.Unsetenv("FIRESTORE_EMULATOR_PORT")
			},
			wantErr: false,
		},
		{
			name:      "Empty project ID",
			projectID: "",
			setup: func(t *testing.T) {
				_ = os.Setenv("FIRESTORE_EMULATOR_HOST", defaultValue)
			},
			reset: func() error {
				return os.Unsetenv("FIRESTORE_EMULATOR_HOST")
			},
			wantErr: true,
		},
		{
			name:      "Invalid project ID",
			projectID: "invalidID",
			setup: func(t *testing.T) {
				_ = os.Setenv("FIRESTORE_EMULATOR_HOST", defaultValue)
			},
			reset: func() error {
				return os.Unsetenv("FIRESTORE_EMULATOR_HOST")
			},
			wantErr: true,
		},
		{
			name:      "Invalid HOST",
			projectID: "test-project",
			setup: func(t *testing.T) {
				_ = os.Setenv("FIRESTORE_EMULATOR_HOST", "invalidhost")
			},
			reset: func() error {
				return os.Unsetenv("FIRESTORE_EMULATOR_HOST")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(t)
			defer func() {
				if err := tc.reset(); err != nil {
					t.Fatalf("failed to exec reset step on testcase %s: %v", tc.name, err)
				}
			}()

			err := emutils.ResetEmulator(tc.projectID)
			if tc.wantErr && err == nil {
				t.Errorf("expected an error, but got none")
			} else if !tc.wantErr && err != nil {
				t.Errorf("did not expect an error, got: %v", err)
			}
		})
	}
}

func TestResetEmulator_AddDocsThenReset(t *testing.T) {
	// Skip if the emulator host is not configured
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set, skipping integration test.")
	}

	ctx := context.Background()
	const (
		projectID  = "test-project"
		collection = "test-collection"
		document   = "testDoc"
	)

	// Create a Firestore client without authentication.
	// This client will connect to the local emulator.
	client, err := firestore.NewClient(ctx, projectID, option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// Insert a document
	_, err = client.Collection(collection).Doc(document).Set(ctx, map[string]interface{}{
		"hello": "world",
	})
	if err != nil {
		t.Fatalf("failed to add document: %v", err)
	}

	// Verify the document exists
	snapshot, err := client.Collection(collection).Doc(document).Get(ctx)
	if err != nil {
		t.Fatalf("failed to retrieve document: %v", err)
	}
	if val, ok := snapshot.Data()["hello"]; !ok || val != "world" {
		t.Fatalf("expected 'hello=world' in document, got: %v", snapshot.Data())
	}

	// Call ResetEmulator to clear data
	if err := emutils.ResetEmulator(projectID); err != nil {
		t.Fatalf("ResetEmulator returned an error: %v", err)
	}

	// Verify the document is no longer found
	_, err = client.Collection(collection).Doc(document).Get(ctx)
	if err == nil {
		t.Fatalf("expected a 'not found' error, but got none")
	}
	// Check for the Firestore NotFound status
	if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound {
		t.Fatalf("expected a NotFound error, got: %v", err)
	}
}
