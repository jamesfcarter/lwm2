package wm

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
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
	Window *xwindow.Window
	Name   string
}
type Clients map[xproto.Window]*Client

func Init() (*Wm, error) {
	var err error
	wm := &Wm{}
	wm.X, err = xgbutil.NewConn()
	if err != nil {
		return nil, err
	}
	wm.loadScreens()
	err = wm.OnEachScreen(loadClients)
	if err != nil {
		return nil, err
	}
	err = wm.OnEachScreen(manageClients)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

func (wm *Wm) OnEachScreen(f func(*Screen) error) error {
	for _, s := range wm.Screens {
		err := f(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wm *Wm) Run() {
	xevent.Main(wm.X)
}

func (wm *Wm) loadScreens() {
	setup := wm.X.Setup()
	wm.Screens = make(Screens, 0, len(setup.Roots))
	for _, si := range setup.Roots {
		wm.Screens = append(wm.Screens, &Screen{
			Root:    xwindow.New(wm.X, si.Root),
			Wm:      wm,
			Clients: make(Clients),
		})
	}
}

func loadClients(s *Screen) error {
	windows, err := s.UnmanagedWindows()
	if err != nil {
		return err
	}
	for _, w := range windows {
		s.Clients[w] = &Client{
			Screen: s,
			Window: xwindow.New(s.X(), w),
		}
	}
	return nil
}

func manageClients(s *Screen) error {
	return s.OnEachClient(func(c *Client) error { return c.Manage() })
}
