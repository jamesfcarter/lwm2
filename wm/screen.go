package wm

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
)

func (s *Screen) X() *xgbutil.XUtil {
	return s.Wm.X
}

func (s *Screen) Conn() *xgb.Conn {
	return s.X().Conn()
}

func (s *Screen) UnmanagedWindows() ([]xproto.Window, error) {
	windows := make([]xproto.Window, 0, 100)
	qtr, err := xproto.QueryTree(s.Conn(), s.Root.Id).Reply()
	if err != nil {
		return nil, err
	}
	for _, win := range qtr.Children {
		if s.Clients[win] != nil {
			continue
		}
		attr, err := xproto.GetWindowAttributes(s.Conn(), win).Reply()
		if err != nil {
			return nil, err
		}
		if attr.OverrideRedirect {
			continue
		}
		windows = append(windows, win)
	}
	return windows, nil
}

func (s *Screen) OnEachClient(f func(*Client) error) error {
	for _, c := range s.Clients {
		err := f(c)
		if err != nil {
			return err
		}
	}
	return nil
}
