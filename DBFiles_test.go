package credenta

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestListUserDataFiles(t *testing.T) {
	mdf, err := ListUserDataFiles(".", "/data/user")
	assert.NoError(t, err)
	for realm, names := range mdf {
		namesString := strings.Join(names, ", ")
		t.Logf("%s realms contains %s", realm, namesString)
	}
}
