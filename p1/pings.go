package p1

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// icmp package struct
type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	Identifier  uint16
	SequenceNum uint16
}
type Res struct {
	From       string      `json:"from"`
	Seq        uint16      `json:"seq"`
	TotalTime  int64       `json:"total_time"`
	FailIp     *net.IPAddr `json:"fail_ip"`
	ReceiveCnt int         `json:"receive_cnt"`
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func getICMP(seq uint16) ICMP {
	icmp := ICMP{
		Type:        8,
		Code:        0,
		CheckSum:    0,
		Identifier:  0,
		SequenceNum: seq,
	}

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.CheckSum = CheckSum(buffer.Bytes())
	buffer.Reset()

	return icmp
}

func sendICMPRequest(icmp ICMP, destAddr *net.IPAddr) (Res, error) {
	conn, err := net.DialIP("ip4:icmp", nil, destAddr)
	if err != nil {
		fmt.Printf("Fail to connect to remote host: %s\n", err)
		return Res{FailIp: destAddr}, err
	}
	defer conn.Close()

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return Res{FailIp: destAddr}, err
	}

	tStart := time.Now()

	conn.SetReadDeadline((time.Now().Add(time.Second * 2)))

	recv := make([]byte, 1024)
	receiveCnt, err := conn.Read(recv)
	if err != nil {
		return Res{FailIp: destAddr}, err
	}

	tEnd := time.Now()
	duration := tEnd.Sub(tStart).Nanoseconds() / 1e6

	var res Res
	res.ReceiveCnt = receiveCnt
	res.From = destAddr.String()
	res.Seq = icmp.SequenceNum
	res.TotalTime = duration

	return res, err
}

func Goping(ip string) {

	host := ip
	raddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		fmt.Printf("Fail to resolve %s, %s\n", host, err)
		return
	}

	for i := 1; i < 2; i++ {
		res1, _ := sendICMPRequest(getICMP(uint16(i)), raddr)
		if res1.FailIp != nil {
			fmt.Println(res1.FailIp, "ping is  fail ")
		} else {
			fmt.Println(res1.From, "ping is successful ")
		}

		time.Sleep(1 * time.Microsecond)
	}
}
