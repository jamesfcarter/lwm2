package wm

import (
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
)

func (c *Client) X() *xgbutil.XUtil {
	return c.Window.X
}

func (c *Client) Conn() *xgb.Conn {
	return c.X().Conn()
}

func (c *Client) Id() xproto.Window {
	return c.Window.Id
}

func (c *Client) getName() {
	name, _ := ewmh.WmNameGet(c.X(), c.Id())
	if name != "" {
		c.Name = name
		return
	}
	name, _ = icccm.WmNameGet(c.X(), c.Id())
	c.Name = name
	return
}

func (c *Client) windowType() string {
	types, err := ewmh.WmWindowTypeGet(c.X(), c.Id())
	if err != nil || len(types) == 0 {
		return "_NET_WM_WINDOW_TYPE_NORMAL"
	}
	return types[0]
}

func (c *Client) Manage() error {
	c.getName()
	fmt.Printf("managing '%s' (%0x) %s\n", c.Name, c.Id(), c.windowType())
	return nil
}
