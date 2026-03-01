package validator

import (
	"path/filepath"
	"testing"

	"github.com/Cliper27/grove/internal/parser"
)

func TestValidate_FullValidationLogic(t *testing.T) {
	tmp := t.TempDir()

	// --- build file structure ---
	mkdir(t, tmp, "config")
	mkdir(t, tmp, "assets")

	touch(t, tmp, "main.go")
	touch(t, tmp, "LICENSE")
	touch(t, tmp, "README.md") // allowed via glob

	// nested required file
	touch(t, tmp, filepath.Join("config", "app.yaml"))

	// allowed asset
	touch(t, tmp, filepath.Join("assets", "logo.png"))

	// denied files
	touch(t, tmp, "file.tmp")
	touch(t, tmp, "debug_trace.log")

	// denied folder
	mkdir(t, tmp, "secrets")

	//
	// --- nested schemas ---
	//

	configSchema := &parser.Schema{
		Name: "config-schema",
		Require: []*parser.Node{
			fileNode("app.yaml"),
		},
	}

	assetsSchema := &parser.Schema{
		Name: "assets-schema",
		Allow: []*parser.Node{
			globFile("*.png"),
		},
	}

	//
	// --- root schema ---
	//

	rootSchema := &parser.Schema{
		Name: "app-layout",

		Require: []*parser.Node{
			dirNode("config", configSchema),
			fileNode("main.go"),
			fileNode("LICENSE"),
		},

		Allow: []*parser.Node{
			dirNode("assets", assetsSchema),
			globFile("*.md"),
		},

		Deny: []*parser.Node{
			dirNode("secrets", nil),
			globFile("*.tmp"),
			regexFile("^debug_.*\\.log$"),
		},
	}

	result := Validate(tmp, rootSchema)

	if result.Valid {
		t.Fatal("expected validation to fail due to denied entries")
	}

	// Ensure deny bubbling works
	if len(result.Reasons) == 0 {
		t.Fatal("expected root node to contain deny reasons")
	}

	// Ensure nested schema worked (config/app.yaml exists)
	configNode := findChild(result, "config")
	if configNode == nil || !configNode.Valid {
		t.Fatal("expected config folder to be valid")
	}

	// Ensure secrets folder is present and marked invalid
	secretsNode := findChild(result, "secrets")
	if secretsNode == nil || secretsNode.Valid {
		t.Fatal("expected secrets folder to be invalid")
	}
}
