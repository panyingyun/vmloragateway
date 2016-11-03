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
	mac    [8]byte
	stat   *gateway.Stat
	wg     sync.WaitGroup
}

func NewBackend(bind string, mac string, longtitude float64, latitude float64, altitude int32) (*Backend, error) {

	addr, err := net.ResolveUDPAddr("udp", bind)
	if err != nil {
		return nil, errBackend
	}

	connect, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, errBackend
	}

	macbytes, err := covertStrToByte(mac)
	if err != nil {
		return nil, err
	}
	var gwstat gateway.Stat
	gwstat.ACKR = 100.0

	gwstat.Alti = altitude
	gwstat.Lati = latitude
	gwstat.Long = longtitude

	gwstat.RXFW = 0
	gwstat.RXNb = 0
	gwstat.RXOK = 0

	gwstat.DWNb = 0

	b := &Backend{
		addr:   addr,
		conn:   connect,
		rxChan: make(chan gateway.PullRespPacket),
		closed: false,
		mac:    macbytes,
		stat:   &gwstat,
	}

	//handle receivePacket
	go func() {
		b.wg.Add(1)
		b.receivePackets()
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
func (b *Backend) SendHeartbeat() error {
	//generate heartbeat...
	var heartbeat gateway.PullDataPacket
	heartbeat.ProtocolVersion = 2
	heartbeat.RandomToken = uint16(rand.Uint32())
	heartbeat.GatewayMAC = b.mac

	//send by udp
	hbyte, err := heartbeat.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = b.conn.Write(hbyte)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"token": heartbeat.RandomToken,
	}).Info("backend: PullData -->")
	return nil
}

//Send (Just a stat every 30ms)
func (b *Backend) SendStatData() error {
	//generate gateway stat ...
	var stPkt gateway.PushDataPacket
	stPkt.ProtocolVersion = 2
	stPkt.RandomToken = uint16(rand.Uint32())

	stPkt.GatewayMAC = b.mac
	var payload gateway.PushDataPayload
	b.stat.Time = gateway.ExpandedTime(time.Now())
	payload.Stat = b.stat
	stPkt.Payload = payload

	//send by udp
	stBytes, err := stPkt.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = b.conn.Write(stBytes)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"token": stPkt.RandomToken,
	}).Info("backend: PushData -->")
	return nil
}

//Send PushData Command(transmitt lora node data here)
func (b *Backend) SendPushData() error {
	return nil
}

func (b *Backend) receivePackets() error {
	buf := make([]byte, 65507) // max udp data size
	for {
		i, err := b.conn.Read(buf)
		if err != nil {
			return fmt.Errorf("backend: read from udp error: %s", err)
		}
		data := make([]byte, i)
		copy(data, buf[:i])
		go func(data []byte) {
			if err := b.handlePacket(data); err != nil {
				log.Errorf("backend handle packet error = ", err)
			}
		}(data)
	}
	return nil
}

func (b *Backend) handlePacket(data []byte) error {
	pt, err := gateway.GetPacketType(data)
	if err != nil {
		return err
	}
	//	log.WithFields(log.Fields{
	//		"type":             pt,
	//		"protocol_version": data[0],
	//	}).Info("backend: received udp packet from bridge")

	switch pt {
	case gateway.PullACK:
		return b.handlePullAck(data)
	case gateway.PushACK:
		return b.handlePushAck(data)
	case gateway.PullResp:
		return b.handlePullResp(data)
	default:
		return fmt.Errorf("backend: unknown packet type=%s", pt)
	}
}

func (b *Backend) handlePushAck(data []byte) error {
	//log.Println("handlePushAck")
	var p gateway.PushACKPacket
	err := p.UnmarshalBinary(data)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"token": p.RandomToken,
	}).Info("backend: PushAck -->")
	return nil
}

func (b *Backend) handlePullAck(data []byte) error {
	//log.Println("handlePullAck")
	var p gateway.PullACKPacket
	err := p.UnmarshalBinary(data)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"token": p.RandomToken,
	}).Info("backend: PullAck -->")
	return nil
}

func (b *Backend) handlePullResp(data []byte) error {
	log.Println("handlePullResp")
	return nil
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
