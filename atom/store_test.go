package atom_test

import (
	"testing"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/jamesfcarter/lwm2/atom"
)

type FakeAtom struct {
	count xproto.Atom
	err   error
}

var fa *FakeAtom = &FakeAtom{}

func (fa *FakeAtom) AtomLoad(name string) (xproto.Atom, error) {
	fa.count += 1
	return fa.count, fa.err
}

func init() {
	atom.Require("test1", "test2", "test3")
}

func TestAtomsLoaded(t *testing.T) {
	store := atom.NewStore(fa)
	for _, test := range []struct {
		name string
		atom xproto.Atom
	}{
		{"test1", 1},
		{"test2", 2},
		{"test3", 3},
	} {
		a := store.Atom(test.name)
		if a != test.atom {
			t.Fatalf("wrong atom %s %d", test.name, a)
		}
	}
}
