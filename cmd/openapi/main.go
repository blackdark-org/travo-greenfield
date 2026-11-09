// Command openapi prints the daemon's OpenAPI 3.1 spec as JSON. It needs no device.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/BlackDark/travo-greenfield/internal/api"
)

// specVersion keeps the committed spec stable; the real build version is not part of the contract.
const specVersion = "0.0.0"

func render(w io.Writer) error {
	_, a := api.New(specVersion)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(a.OpenAPI())
}

func main() {
	if err := render(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "openapi:", err)
		os.Exit(1)
	}
}
