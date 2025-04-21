package credenta

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroup(t *testing.T) {
	g := &CGroup{
		Name:         "G1",
		ParentGroups: nil,
		Attributes:   []*Attribute{},
	}

	assert.NoError(t, g.SetiAttribute("G1N1", 1))
	assert.Error(t, g.SetiAttribute("G1N1", 2))

	assert.NoError(t, g.SetsAttribute("G1S1", "One"))
	assert.Error(t, g.SetsAttribute("G1S1", "Two"))

	//fmt.Println(g.GetAttributeList())
}
