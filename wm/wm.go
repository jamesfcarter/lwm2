package wm

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Wm struct {
	X       *xgbutil.XUtil
	Screens Screens
}

type Screen struct {
	Wm      *Wm
	Root    *xwindow.Window
	Clients Clients
}
type Screens []*Screen

type Client struct {
	Screen *Screen
}
type Clients []*Client

func Init() (*Wm, error) {
	var err error
	wm := &Wm{}
	wm.X, err = xgbutil.NewConn()
	if err != nil {
		return nil, err
	}
	wm.loadScreens()
	return wm, nil
}

func (wm *Wm) loadScreens() {
	setup := wm.X.Setup()
	wm.Screens = make(Screens, 0, len(setup.Roots))
	for _, si := range setup.Roots {
		wm.Screens = append(wm.Screens, &Screen{
			Root: xwindow.New(wm.X, si.Root),
			Wm:   wm,
		})
	}
}
