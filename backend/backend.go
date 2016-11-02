package backend

import (
	_ "encoding/base64"
	"errors"
	_ "fmt"
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/panyingyun/vmloragateway/gateway"
)

var errBackend = errors.New("backend create fail!")

type Backend struct {
	conn   *net.UDPConn
	rxChan chan gateway.PullRespPacket
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
		rxChan: make(chan gateway.PullRespPacket),
		closed: false,
	}

	//handle receivePacket
	go func() {
		b.wg.Add(1)
		b.receivePackets()
		//		if !b.closed {
		//			log.Fatal(err)
		//		}
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

//Send PullData Command(Just a heartbeat every 10ms)
func (b *Backend) SendPullData() error {

	return nil
}

//Send (Just a stat every 30ms)
func (b *Backend) SendStatData() error {
	return nil
}

//Send PushData Command(transmitt lora node data here)
func (b *Backend) SendPushData() error {
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

func (b *Backend) handlePushAck() {

}

func (b *Backend) handlePullAck() {

}
func (b *Backend) handlePullResp() {

}
