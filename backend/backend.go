package backend

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	log "github.com/fatih/color"
	"github.com/panyingyun/vmloragateway/gateway"
)

var errBackend = errors.New("backend create fail!")

type Backend struct {
	conn   *net.UDPConn
	addr   *net.UDPAddr
	txChan chan gateway.TXPK
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
		txChan: make(chan gateway.TXPK),
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
	log.Cyan("backend: closing gateway backend")
	b.closed = true
	if err := b.conn.Close(); err != nil {
		return err
	}
	log.Cyan("backend: handling last packets")
	b.wg.Wait()
	return nil
}

//Send Join or comfirm/uncomfirm uplink data from lora node
func (b *Backend) SendPushData(rxpk gateway.RXPK) error {
	//generate gateway stat ...
	var rxPkt gateway.PushDataPacket
	rxPkt.ProtocolVersion = 2
	rxPkt.RandomToken = uint16(rand.Uint32())

	rxPkt.GatewayMAC = b.mac
	var payload gateway.PushDataPayload
	payload.RXPK = append(payload.RXPK, rxpk)
	rxPkt.Payload = payload

	//send by udp
	stBytes, err := rxPkt.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = b.conn.Write(stBytes)
	if err != nil {
		return err
	}
	log.Green("backend: PushData(Join or Uplink) --> token = %v", rxPkt.RandomToken)
	return nil
}

//Send TX_ACK When received PULL_RESP packet
func (b *Backend) sendTXACK(token uint16) error {
	//generate heartbeat...
	var txack gateway.TXACKPacket
	txack.ProtocolVersion = 2
	txack.RandomToken = token
	txack.GatewayMAC = b.mac

	//send by udp
	hbyte, err := txack.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = b.conn.Write(hbyte)
	if err != nil {
		return err
	}

	log.Green("backend: TXACK --> token = %v ", txack.RandomToken)
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

	log.Green("backend: PullData --> token = %v ", heartbeat.RandomToken)
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

	log.Green("backend: PushData(Stat) --> token = %v ", stPkt.RandomToken)
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
				log.Red("backend handle packet error = ", err)
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

	log.Green("backend: PushAck --> token = %v ", p.RandomToken)

	return nil
}

func (b *Backend) handlePullAck(data []byte) error {
	//log.Println("handlePullAck")
	var p gateway.PullACKPacket
	err := p.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	log.Green("backend: PullAck --> token = %v ", p.RandomToken)
	return nil
}

func (b *Backend) handlePullResp(data []byte) error {
	//log.Println("handlePullResp")
	var p gateway.PullRespPacket
	err := p.UnmarshalBinary(data)
	if err != nil {
		return err
	}
	log.Green("backend: PullResp --> token = %v TXPK = %v", p.RandomToken, p.Payload.TXPK)
	b.sendTXACK(p.RandomToken)
	b.txChan <- p.Payload.TXPK
	return nil
}

func (b *Backend) TXChan() chan gateway.TXPK {
	return b.txChan
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
