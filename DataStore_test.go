package credenta

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCredentaDB_GetUser(t *testing.T) {
	cDB, err := NewCredentaDB()
	assert.NoError(t, err)

	ctx := context.WithValue(context.Background(), ETX_USER, "TestUser")

	u, err := cDB.NewUser(ctx, "DEFAULT", "USERID", "password", nil, IdTypeUserId, VerificationMethodPLAIN)
	assert.NoError(t, err)

	err = u.StoreOrSaveToFile(ctx)
	assert.NoError(t, err)

	u.AddRole(0)
	u.AddRole(1)
	u.AddRole(64)

	assert.NoError(t, u.SetsAttribute("ATTRA", "AttributeA"))
	assert.True(t, u.HasAttribute("ATTRA"))
	assert.NoError(t, u.SetiAttribute("ATTRB", 123))
	assert.True(t, u.HasAttribute("ATTRB"))
	assert.NoError(t, u.SetiAttribute("ATTRC", 123))
	assert.True(t, u.HasAttribute("ATTRC"))

	t.Log("------------")

	assert.True(t, u.HasAttribute("ATTRA"))
	assert.True(t, u.HasAttribute("ATTRB"))

	t.Log("------------")

	t.Log(u)
	u.RemoveRole(1)
	u.RemoveRole(64)
	t.Log(u)

	_, err = cDB.NewUser(ctx, "DEFAULT", "USERID", "password", nil, IdTypeUserId, VerificationMethodPLAIN)
	assert.Error(t, err)

	usr, err := cDB.GetUser(ctx, "DEFAULT", "USERID")
	assert.NoError(t, err)
	assert.NoError(t, usr.DeleteFile(ctx))

}

func TestCredentaDB_GetRoleMasksOfGroups(t *testing.T) {
	cDB, err := NewCredentaDB()
	assert.NoError(t, err)

	ctx := context.WithValue(context.Background(), ETX_USER, "TestUser")

	elder, err := cDB.NewGroup(ctx, "RA", "GroupElder", nil)
	assert.NoError(t, err)
	elder.AddRole(0)
	assert.True(t, elder.HasRole(0))
	assert.NoError(t, elder.StoreOrSaveToFile(ctx))

	eGroup, err := cDB.GetGroup(ctx, "RA", "GroupElder")
	assert.NoError(t, err)
	assert.NotNil(t, eGroup)

	son, err := cDB.NewGroup(ctx, "RA", "GroupSon", []string{"GroupElder"})
	assert.NoError(t, err)
	son.AddRole(1)
	assert.True(t, son.HasRole(1))
	assert.NoError(t, son.StoreOrSaveToFile(ctx))

	grandson, err := cDB.NewGroup(ctx, "RA", "GroupGrand", []string{"GroupSon"})
	assert.NoError(t, err)
	grandson.AddRole(2)
	assert.True(t, grandson.HasRole(2))
	assert.NoError(t, grandson.StoreOrSaveToFile(ctx))

	roles := cDB.GetRoleMasksOfGroups(ctx, "RA", "GroupGrand")
	assert.True(t, isRoleFlagOn(roles, 0))
	assert.True(t, isRoleFlagOn(roles, 1))
	assert.True(t, isRoleFlagOn(roles, 2))
	assert.False(t, isRoleFlagOn(roles, 3))

	assert.NoError(t, elder.DeleteFile(ctx))
	assert.NoError(t, son.DeleteFile(ctx))
	assert.NoError(t, grandson.DeleteFile(ctx))

}

func isRoleFlagOn(roles []uint64, roleSquence int) bool {
	seq, bit := toUint64ByBit(roleSquence)
	return isBitFlagOn(roles[seq], bit)
}
