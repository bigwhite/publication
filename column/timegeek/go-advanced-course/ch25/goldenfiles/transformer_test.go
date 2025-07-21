package goldenfiles

import (
	"flag" // To define the -update flag
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp" // For better diffs
)

var update = flag.Bool("update", false, "Update golden files with actual output.")

// TestMain can be used to parse flags.
func TestMain(m *testing.M) {
	flag.Parse()
	code := m.Run()
	os.Exit(code)
}

func TestTransformText_Golden(t *testing.T) {
	testCases := []struct {
		name      string
		inputFile string // Relative to testdata/
		// Golden file will be inputFile with .golden suffix
	}{
		{name: "Case1_HelloWorld", inputFile: "case1_input.txt"},
		{name: "Case2_GoTesting", inputFile: "case2_input.txt"},
		{name: "Case3_EmptyInput", inputFile: "case3_empty_input.txt"}, // Expects an empty input file
		{name: "Case4_SpecialChars", inputFile: "case4_special_input.txt"},
	}

	// Create dummy input files for Case3 and Case4 if they don't exist
	// This is just for the example to be self-contained for generation.
	// In a real scenario, these files would already exist with meaningful content.
	os.MkdirAll(filepath.Join("testdata"), 0755)
	os.WriteFile(filepath.Join("testdata", "case3_empty_input.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join("testdata", "case4_special_input.txt"), []byte("你好, 世界 & < > \" '"), 0644)

	for _, tc := range testCases {
		currentTC := tc // Capture range variable
		t.Run(currentTC.name, func(t *testing.T) {
			// t.Parallel() // Golden file tests often modify files, so parallelism needs care

			inputPath := filepath.Join("testdata", currentTC.inputFile)
			// Golden file path is derived from input file name
			goldenPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "_output.golden"

			inputBytes, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatalf("Failed to read input file %s: %v", inputPath, err)
			}
			inputContent := string(inputBytes)

			actualOutput := TransformText(inputContent)

			if *update { // Check if the -update flag is set
				// Write the actual output to the golden file
				err := os.WriteFile(goldenPath, []byte(actualOutput), 0644)
				if err != nil {
					t.Fatalf("Failed to update golden file %s: %v", goldenPath, err)
				}
				t.Logf("Golden file %s updated.", goldenPath)
				// After updating, we might want to skip the comparison for this run,
				// or we can let it compare to ensure what we wrote is what we get if read back.
				// For this example, we'll just log and continue (which means it might PASS if written correctly).
			}

			// Read the golden file for comparison
			expectedOutputBytes, err := os.ReadFile(goldenPath)
			if err != nil {
				// If golden file doesn't exist and we are not in -update mode, it's an error.
				// Or, it could mean this is the first run for a new test case,
				// and you might want to automatically create it (similar to -update).
				// For strictness, we'll consider it an error here if not in -update mode.
				t.Fatalf("Failed to read golden file %s: %v. Run with -args -update to create it.", goldenPath, err)
			}
			expectedOutput := string(expectedOutputBytes)

			// Compare actual output with golden file content
			// Using github.com/google/go-cmp/cmp for better diffs
			if diff := cmp.Diff(expectedOutput, actualOutput); diff != "" {
				t.Errorf("TransformText() output does not match golden file %s. Diff (-golden +actual):\n%s",
					goldenPath, diff)
			}
		})
	}
}
