package scheduler

import (
	"encoding/hex"
	"fmt"
	"net/netip"
	"sync"
	"time"

	"github.com/millken/golog"
	"github.com/millken/mkdns/internal/wire"
)

type Scheduler struct {
	// contains filtered or unexported fields
	txChan chan []byte
	rxChan chan []byte
	/* Stop channel */
	stop     chan struct{}
	stopOnce sync.Once
}

func New() *Scheduler {
	return &Scheduler{
		txChan: make(chan []byte, 512),
		rxChan: make(chan []byte, 512),
	}
}

func (s *Scheduler) Start() {
	s.stop = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case tx := <-s.txChan:
				s.handleTx(tx)
			case <-s.stop:
				return
			}
		}
	}()
}

func (s *Scheduler) TxChan() chan []byte {
	return s.txChan
}

func (s *Scheduler) RxChan() chan []byte {
	return s.rxChan
}

func (s *Scheduler) handleTx(txBytes []byte) {
	//packet := layers.NewPacket(txBytes, layers.Default)
	fmt.Println(hex.Dump(txBytes))
	fmt.Printf("%x\n", txBytes[:])
	return
	req, resp := wire.AcquireMessage(), wire.AcquireMessage()
	defer wire.ReleaseMessage(req)
	defer wire.ReleaseMessage(resp)
	if err := wire.ParseMessage(req, txBytes, false); err != nil {
		golog.Errorf("parse req message error : %s", err)
		return
	}
	golog.Infof("%s", req)

	client := &wire.Client{
		AddrPort:    netip.AddrPortFrom(netip.AddrFrom4([4]byte{8, 8, 8, 8}), 53),
		ReadTimeout: 1 * time.Second,
		MaxConns:    1000,
	}

	if err := client.Exchange(req, resp); err != nil {
		golog.Errorf("Exchange error : %s", err)
		return
	}
	wire.Short(resp)
}

func (s *Scheduler) handleRx(txBytes []byte) {
	fmt.Println(hex.Dump(txBytes))
}

func (s *Scheduler) Stop() {
	s.stopOnce.Do(func() {
		s.stop <- struct{}{}
	})
}
