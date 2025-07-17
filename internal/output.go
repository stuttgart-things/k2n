package internal

import (
	"fmt"
	"os"
)

// SaveOutput writes the content to the given destination or stdout if destination is empty.
func SaveOutput(destination, content string) error {
	if destination == "" {
		fmt.Println(content)
		return nil
	}

	err := os.WriteFile(destination, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %w", destination, err)
	}

	fmt.Printf("Result written to %s\n", destination)
	return nil
}
