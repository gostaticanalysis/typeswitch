package typeswitch_test

import (
	"testing"

	"github.com/gostaticanalysis/typeswitch"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, typeswitch.Analyzer, "a")
}