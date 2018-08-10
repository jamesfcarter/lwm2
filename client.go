package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Client struct {
	Window       xproto.Window
	Parent       xproto.Window
	TransientFor xproto.Window
	Screen       *Screen
}
type Clients map[xproto.Window]*Client

func (c *Client) Manage(x *xgb.Conn) {
}
