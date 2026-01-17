package parser

import (
	"fmt"
	"log"
	"path/filepath"
)

// getExampleFilePath is a helper for examples to locate test .gro files.
func getExampleFilePath(file string) string {
	return filepath.Join("..", "..", "test_data", file)
}

// ExampleLoadSchema demonstrates loading a schema with includes.
func ExampleLoadSchema() {
	path := getExampleFilePath("happy/go-project.gro")

	schema, err := LoadSchema(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(schema.Name)
	fmt.Println(len(schema.Require))
	fmt.Println(len(schema.Allow))
	fmt.Println(len(schema.Deny))

	// Output:
	// go-project
	// 6
	// 1
	// 3
}

// ExampleParseByteUnits demonstrates parsing human-readable byte units.
func ExampleParseByteUnits() {
	tests := []string{"1B", "1KB", "2MB", "3GB", "5TB", "10"}

	for _, s := range tests {
		bytes, err := ParseByteUnits(s)
		if err != nil {
			fmt.Printf("%s: error %v\n", s, err)
			continue
		}
		fmt.Printf("%s = %d bytes\n", s, bytes)
	}

	// Output:
	// 1B = 1 bytes
	// 1KB = 1024 bytes
	// 2MB = 2097152 bytes
	// 3GB = 3221225472 bytes
	// 5TB = 5497558138880 bytes
	// 10 = 10 bytes
}
