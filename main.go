package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"

	"log"
)

type Strut struct {
	Left, Right, Top, Bottom int
}

type Screen struct {
	xproto.ScreenInfo
}
type Screens []*Screen

func initXinerama(x *xgb.Conn) []xinerama.ScreenInfo {
	err := xinerama.Init(x)
	if err != nil {
		log.Printf("failed to init xinerama: %v", err)
	}

	qsr, err := xinerama.QueryScreens(x).Reply()
	if err != nil {
		log.Printf("failed to query xinerama screens: %v", err)
	}
	return qsr.ScreenInfo
}

func initScreens(x *xgb.Conn) Screens {
	setup := xproto.Setup(x)
	screens := make(Screens, 0, len(setup.Roots))
	for _, si := range setup.Roots {
		screens = append(screens, &Screen{
			ScreenInfo: si,
		})
	}
	return screens
}

func initClients(x *xgb.Conn, screens Screens) Clients {
	clients := make(Clients)
	for i, screen := range screens {
		qtr, err := xproto.QueryTree(x, screen.Root).Reply()
		if err != nil {
			log.Printf("failed to query screen %d's tree: %v", i, err)
			continue
		}
		for _, win := range qtr.Children {
			if clients[win] != nil {
				continue
			}
			attr, err := xproto.GetWindowAttributes(x, win).Reply()
			if err != nil {
				log.Printf("failed to read window attributes: %v", err)
				continue
			}
			if attr.OverrideRedirect {
				continue
			}
			clients[win] = &Client{
				Screen: screen,
			}
		}
	}
	return clients
}

func main() {
	x, err := xgb.NewConn()
	if err != nil {
		log.Fatalf("failed to init xgb: %v", err)
	}
	_ = initXinerama(x)
	screens := initScreens(x)
	log.Printf("%#v\n%d\n", screens, len(screens))
	clients := initClients(x, screens)
	log.Printf("%#v\n%d\n", clients, len(clients))
}
