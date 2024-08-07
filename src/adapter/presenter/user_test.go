package presenter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutputCreateResult(t *testing.T) {
	/* Arrange */
	expected := "{}\n"
	c, rec := newRouter()
	up := &UserPresenter{c: c}

	/* Act */
	actual := up.OutputCreateResult()

	/* Assert */
	// up.OutputCreateResultがJSONを返すこと
	if assert.NoError(t, actual) {
		assert.Equal(t, expected, rec.Body.String())
	}
}
