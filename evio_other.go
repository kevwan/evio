// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// +build !linux

package evio

import (
	"errors"
	"net"
	"os"
)

func (ln *listener) close() {
	if ln.ln != nil {
		ln.ln.Close()
	}
	if ln.network == "unix" {
		os.RemoveAll(ln.addr)
	}
}

func (ln *listener) system() error {
	return nil
}

func serve(events Events, listener *listener) error {
	return stdserve(events, listener)
}

func reuseportListenPacket(proto, addr string) (l net.PacketConn, err error) {
	return nil, errors.New("reuseport is not available")
}

func reuseportListen(proto, addr string) (l net.Listener, err error) {
	return nil, errors.New("reuseport is not available")
}
