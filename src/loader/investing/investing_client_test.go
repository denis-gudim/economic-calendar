package investing

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuffleRequestParams(t *testing.T) {
	// Arrange
	params := url.Values{}

	for i := 0; i < 26; i++ {
		params.Add(strconv.Itoa(i), string(rune('A'+i)))
	}

	expected := params.Encode()

	// Act
	actual := shuffleRequestParams(&params)

	// Assert
	assert.NotEqual(t, expected, actual)
}
