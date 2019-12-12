package contrib

import "testing"

func TestJSONMarshalMust(t *testing.T) {
	var (
		s = `{name:10,"age":999}`
	)
	t.Log(string(JSONMarshalMust(s)))
}
