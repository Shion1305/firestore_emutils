# EMUTILS for Firestore

EMUTILS for Firestore is a small Go library designed to help you manage data in the official Firestore emulator. It’s particularly useful for integration tests or local development environments where resetting the emulator state quickly is important.

[![codecov](https://codecov.io/gh/Shion1305/firestore_emutils/graph/badge.svg?token=fi23EADawz)](https://codecov.io/gh/Shion1305/firestore_emutils)

## Features
- Clear all data in your local Firestore emulator with a single call.


## Usage

### Import this library

```bash
go get "github.com/Shion1305/firestore_emutils"
```

### Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/Shion1305/firestore_emutils"
)
func main() {
    // Configure these values to match your local setup
    host := "localhost"
    port := 8080
    projectID := "my_project"

    // Create a new Emulator instance
    emu := emutils.NewEmulator(host, port, projectID)

    // Clear all data from the emulator
    if err := emu.ClearData(); err != nil {
        log.Fatalf("Failed to clear emulator data: %v", err)
    }

    fmt.Println("Firestore emulator data cleared successfully!")
}
```

## Contributing

Contributions are always welcome! Please open a GitHub issue for bug reports or feature requests.
