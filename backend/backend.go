package backend

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/panyingyun/vmloragateway/gateway"
)

var errBackend = errors.New("backend create fail!")

type Backend struct {
	conn   *net.UDPConn
	addr   *net.UDPAddr
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
		addr:   addr,
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

func (b *Backend) sendPullData() error {

	return nil
}

//Send PullData Command(Just a heartbeat every 10ms)
func (b *Backend) SendHeartbeat(mac string) error {
	//generate heartbeat...
	var heartbeat gateway.PullDataPacket
	heartbeat.ProtocolVersion = 2
	heartbeat.RandomToken = uint16(rand.Uint32())
	macbt, err := covertStrToByte(mac)
	if err != nil {
		return err
	}
	heartbeat.GatewayMAC = macbt

	//send by udp
	hbyte, err := heartbeat.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = b.conn.Write(hbyte)
	if err != nil {
		return err
	}
	return nil
}

//Send (Just a stat every 30ms)
func (b *Backend) SendStatData(latitude float64, longtitude float64, altitude int32, mac string) error {
	//generate gateway stat ...
	var stat gateway.PushDataPacket
	stat.ProtocolVersion = 2
	stat.RandomToken = uint16(rand.Uint32())
	macbt, err := covertStrToByte(mac)
	//log.Info("macbt = ", macbt)
	if err != nil {
		return err
	}
	stat.GatewayMAC = macbt
	var payload gateway.PushDataPayload
	var st gateway.Stat
	st.Time = gateway.ExpandedTime(time.Now())
	st.ACKR = 100.0
	st.Alti = altitude
	st.Lati = latitude
	st.Long = longtitude
	st.RXFW = 100
	st.RXNb = 100
	st.RXOK = 100
	st.DWNb = 0
	payload.Stat = &st
	stat.Payload = payload

	//send by udp
	stBytes, err := stat.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = b.conn.Write(stBytes)
	if err != nil {
		return err
	}
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

func covertStrToByte(mac string) ([8]byte, error) {
	var macbytes [8]byte
	b, err := hex.DecodeString(mac)
	if err != nil {
		return macbytes, err
	}
	if len(b) != len(macbytes) {
		return macbytes, fmt.Errorf("macbytes: exactly %d bytes are expected", len(b))
	}
	copy(macbytes[:], b)
	return macbytes, nil
}
