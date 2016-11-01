package backend

import (
	_ "encoding/base64"
	"errors"
	_ "fmt"
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
)

var errBackend = errors.New("backend create fail!")

type udpPacket struct {
	addr *net.UDPAddr
	data []byte
}

type Backend struct {
	conn   *net.UDPConn
	rxChan chan udpPacket
	txChan chan udpPacket
	closed bool
	wg     sync.WaitGroup
}

func NewBackend(bind string) (*Backend, error) {

	addr, err := net.ResolveUDPAddr("udp", bind)
	if err != nil {
		return nil, errBackend
	}

	connect, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, errBackend
	}

	b := &Backend{
		conn:   connect,
		rxChan: make(chan udpPacket),
		txChan: make(chan udpPacket),
		closed: false,
	}

	//handle receivePacket
	go func() {
		b.wg.Add(1)
		err := b.receivePackets()
		if !b.closed {
			log.Fatal(err)
		}
		b.wg.Done()
	}()

	//handle sendPacket
	go func() {
		b.wg.Add(1)
		err := b.sendPacket()
		if !b.closed {
			log.Fatal(err)
		}
		b.wg.Done()
	}()
	return b, nil
}

func (b *Backend) Close() error {
	log.Info("backend: closing gateway backend")
	b.closed = true
	if err := b.conn.Close(); err != nil {
		return err
	}
	log.Info("backend: handling last packets")
	b.wg.Wait()
	return nil
}

func (b *Backend) receivePackets() error {
	//	buf := make([]byte, 65507) // max udp data size
	//	for {
	//		i, addr, err := b.conn.ReadFromUDP(buf)
	//		if err != nil {
	//			return fmt.Errorf("gateway: read from udp error: %s", err)
	//		}
	//		data := make([]byte, i)
	//		copy(data, buf[:i])
	//		go func(data []byte) {
	//			if err := b.handlePacket(addr, data); err != nil {
	//				log.WithFields(log.Fields{
	//					"data_base64": base64.StdEncoding.EncodeToString(data),
	//					"addr":        addr,
	//				}).Errorf("gateway: could not handle packet: %s", err)
	//			}
	//		}(data)
	//	}

	return nil
}

func (b *Backend) sendPacket() error {
	return nil
}
