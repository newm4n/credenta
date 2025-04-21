package credenta

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCredentaDB_GetUser(t *testing.T) {
	cdb := NewCredentaDB()

	u, err := cdb.NewUser("DEFAULT", "USERID", "password", nil, IdTypeUserId, VerificationMethodPLAIN)

	assert.NoError(t, err)
	u.AddRole(0)
	u.AddRole(1)
	u.AddRole(64)
	assert.NoError(t, u.SetsAttribute("ATTRA", "AttributeA"))
	assert.True(t, u.HasAttribute("ATTRA"))

	t.Log(u)
	u.RemoveRole(1)
	u.RemoveRole(64)
	t.Log(u)

	_, err = cdb.NewUser("DEFAULT", "USERID", "password", nil, IdTypeUserId, VerificationMethodPLAIN)
	assert.Error(t, err)
}
