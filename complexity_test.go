package complexity_test

import (
	"testing"

	"github.com/shoooooman/go-complexity-analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, complexity.Analyzer, []string{"halstead"}...)
}
