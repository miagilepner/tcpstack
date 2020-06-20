package socket

import "io"

type Socket interface {
	io.ReadWriteCloser
	Attach(iface string) error
}
