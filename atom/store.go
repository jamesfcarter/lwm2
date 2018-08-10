package atom

import (
	"log"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Store struct {
	sync.RWMutex
	atoms map[string]xproto.Atom
}

type AtomLoader interface {
	AtomLoad(string) (xproto.Atom, error)
}

type XAtomLoader struct {
	conn *xgb.Conn
}

func (x *XAtomLoader) AtomLoad(name string) (xproto.Atom, error) {
	reply, err := xproto.InternAtom(x.conn, true, uint16(len(name)),
		name).Reply()
	return reply.Atom, err
}

func Loader(x *xgb.Conn) *XAtomLoader {
	return &XAtomLoader{x}
}

var required []string

func Require(names ...string) {
	required = append(required, names...)
}

// NewStore returns a store of interned Atoms. In normal use the
// AtomLoader argument should be atom.Loader(x) where x is a *xgb.Conn
func NewStore(x AtomLoader) *Store {
	store := &Store{
		atoms: make(map[string]xproto.Atom),
	}
	for _, name := range required {
		atom, err := x.AtomLoad(name)
		if err != nil {
			log.Fatalf("failed to intern atom %s: %v", name, err)
		}
		store.atoms[name] = atom
	}
	return store
}

func (s *Store) Atom(name string) xproto.Atom {
	s.RLock()
	defer s.RUnlock()
	return s.atoms[name]
}
