package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigValid(t *testing.T) {
	actual, err := parseConfig("testdata/configs/valid.yml")

	expected := configRoot{
		Binaries: []configBinary{
			{
				Name:    "foo",
				Source:  "https://foo.com/{{ .Version }}/bin",
				Version: "1.2.0",
			},
			{
				Name:    "bar",
				Source:  "https://bar.com/releases/v{{ .Version }}/assets/v{{ .Version }}-bin",
				Version: "0.16.5",
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
