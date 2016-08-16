package base

import (
	"reflect"
	"testing"
)

func TestRemove(b *testing.T) {
	keys := []string{"a", "b", "c", "d", "e"}
	mp := NewAvaliableMap(keys)
	mp.Remove("a")
	if !reflect.DeepEqual(mp.keys, []string{"b", "c", "d", "e"}) {
		b.Error("a:值不相等:real:", mp.keys, "expect:", []string{"b", "c", "d", "e"})
	}
	mp.Remove("e")
	if !reflect.DeepEqual(mp.keys, []string{"b", "c", "d"}) {
		b.Error("e:值不相等:real:", mp.keys, []string{"b", "c", "d"})
	}
	mp.Remove("c")
	if !reflect.DeepEqual(mp.keys, []string{"b", "d"}) {
		b.Error("c:值不相等:real:", mp.keys, []string{"b", "d"})
	}
	mp.Remove("a")
	if !reflect.DeepEqual(mp.keys, []string{"b", "d"}) {
		b.Error("c:值不相等:real:", mp.keys, []string{"b", "d"})
	}
	mp.Remove("b")
	if !reflect.DeepEqual(mp.keys, []string{"d"}) {
		b.Error("b:值不相等:real:", mp.keys, []string{"d"})
	}
	mp.Remove("d")
	if !reflect.DeepEqual(mp.keys, []string{}) {
		b.Error("d:值不相等:real:", mp.keys, []string{})
	}
}
