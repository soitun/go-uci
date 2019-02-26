package uci

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadExpected(t *testing.T, name string) *config {
	t.Helper()

	f, err := os.Open(filepath.Join("testdata", name+".json"))
	if err != nil {
		t.Fatalf("cannot open %s.json: %v", name, err)
	}
	defer f.Close()

	expected := &config{}
	err = json.NewDecoder(f).Decode(&expected)
	if err != nil {
		t.Fatalf("error decoding json: %v", err)
	}

	// The JSON dump does not contain empty slices (they're marked with
	// "omitempty"), but the decoder creates them anyway. To get the tests
	// to pass, we need to eliminate nil slices (sections of config and
	// options of section) manually.
	if expected.Sections == nil {
		expected.Sections = []*section{}
	}
	for _, sec := range expected.Sections {
		if sec.Options == nil {
			sec.Options = []*option{}
		}
	}
	return expected
}

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)

	for _, name := range []string{"system", "emptyfile", "emptysection", "luci", "ucitrack"} {
		t.Run(name, func(t *testing.T) {
			r := NewTree("testdata")
			err := r.LoadConfig(name, false)
			assert.NoError(err)

			actual := r.(*tree).configs[name]

			if dump["json"] {
				json.NewEncoder(os.Stderr).Encode(actual)
			}

			expected := loadExpected(t, name)
			assert.EqualValues(expected, actual)
		})
	}
}
