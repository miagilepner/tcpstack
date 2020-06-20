package socket

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"github.com/rs/zerolog/log"
)

type bpfSocket struct {
	fd        int
	filename  string
	f         *os.File
	closeOnce sync.Once
}

func NewSocket() Socket {
	return &bpfSocket{}
}

func (b *bpfSocket) Attach(iface string) error {
	for i := 0; i < 5; i++ {
		filename := fmt.Sprintf("/dev/bpf%d", i)
		fd, err := syscall.Open(filename, syscall.O_WRONLY, 0666)
		if fd != -1 || !errors.Is(err, syscall.EBUSY) {
			b.fd = fd
			b.filename = filename
			break
		}
	}
	if b.fd == 0 {
		return fmt.Errorf("failed to find open bpf file")
	}
	b.f = os.NewFile(uintptr(b.fd), b.filename)
	var ifaceRequest [16]byte
	for i, r := range iface {
		ifaceRequest[i] = byte(r)
	}
	_, _, errSetIface := syscall.Syscall(syscall.SYS_IOCTL, uintptr(b.fd), syscall.BIOCSETIF, uintptr(unsafe.Pointer(&ifaceRequest[0])))
	if errSetIface != 0 {
		closeErr := b.Close()
		if closeErr != nil {
			log.Error().Err(closeErr).Msg("closing")
		}
		return os.NewSyscallError("ioctl", errSetIface)
	}
	var setHeader uint32 = 1
	_, _, errSetHeader := syscall.Syscall(syscall.SYS_IOCTL, uintptr(b.fd), syscall.BIOCSHDRCMPLT, uintptr(unsafe.Pointer(&setHeader)))
	if errSetHeader != 0 {
		closeErr := b.Close()
		if closeErr != nil {
			log.Error().Err(closeErr).Msg("closing")
		}
		return os.NewSyscallError("ioctl", errSetHeader)
	}
	return nil
}

func (b *bpfSocket) Read(data []byte) (int, error) {
	return 0, nil
}

func (b *bpfSocket) Write(data []byte) (int, error) {
	return 0, nil
}

func (b *bpfSocket) Close() error {
	var err error
	b.closeOnce.Do(func() {
		err = b.f.Close()
	})
	return err
}
