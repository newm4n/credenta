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

	assert.NoError(t, g.SetAttribute("G1N1", "int", "1"))
	assert.Error(t, g.SetAttribute("G1N1", "int", "2"))

	//fmt.Println(g.GetAttributeList())
}
