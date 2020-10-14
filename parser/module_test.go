package parser

import (
	"testing"
)

func TestModule(t *testing.T) {
	t.Run("Basic module parsing", func(t *testing.T) {
		m := NewRootModule("testdata/module_consumer")

		err := m.Resolve()
		if err != nil {
			t.Errorf("Failed to evaluate module %s: %s", m.ModuleBlock.Name, err.Error())
			t.FailNow()
		}

		// assert.Equal(t, b.BasicBlocks[0].FirstLabel, "first_label")
		// assert.Equal(t, b.BasicBlocks[0].SecondLabel, "second_label")
		// assert.Equal(t, b.BasicBlocks[0].Number, 42)
		// assert.Equal(t, b.BasicBlocks[0].String, "spiffing")
		// assert.Equal(t, b.BasicBlocks[0].OptionalString, "")
		// processHCL(t, stringHCL, &m)

		// for _, attr := range m.Modules[0].Params.Attributes {
		// 	println(attr.Name)
		// }
	})
}
