package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type Message struct {
	Method string      `json:"_method"`
	To     string      `json:"to,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func RunWebRTC(stuns []string, signaling, uuid, to, file string, maxSize int64, timeout time.Duration, verbose bool) error {
	if len(stuns) == 0 {
		return errors.New("no stun server")
	}
	if uuid == "" || to == "" {
		return errors.New("no pair(uuid,to)")
	}
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	s, err := f.Stat()
	if err != nil {
		return err
	}
	if s.IsDir() {
		return errors.New("not a file")
	}

	size := s.Size()
	if size > maxSize {
		return errors.New("file too large")
	}

	conn, _, err := websocket.DefaultDialer.Dial(signaling, nil)
	if err != nil {
		return err
	}

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{ICEServers: []webrtc.ICEServer{
		{URLs: stuns},
	}})
	if err != nil {
		conn.Close()
		return err
	}

	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			conn.WriteJSON(&Message{"sendto", to, map[string]interface{}{"type": "candidate", "candidate": candidate.ToJSON()}})
		}
	})

	var Ordered = true
	var MaxRetransmits uint16 = 30
	if _, err = pc.CreateDataChannel(uuid, &webrtc.DataChannelInit{
		Ordered:        &Ordered,
		MaxRetransmits: &MaxRetransmits,
	}); err != nil {
		conn.Close()
		return err
	}

	channelOpened := false

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		dc.OnOpen(func() {
			channelOpened = true
			fmt.Println("transport file ...", file, size)
			conn.Close()
			b := make([]byte, 32*1024)
			for {
				n, err := f.Read(b)
				dc.Send(b[:n])
				if err != nil {
					break
				}
			}
			f.Close()
		})
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			dc.Close()
			pc.Close()
			conn.Close()
		})
		dc.OnClose(func() {
			fmt.Println("dc closed")
		})
	})

	go func() {
		for {
			conn.SetReadDeadline(time.Now().Add(timeout))
			m := &Message{}
			if err := conn.ReadJSON(m); err != nil {
				if !channelOpened {
					f.Close()
				}
				conn.Close()
				return
			}
			if verbose {
				log.Printf("%#v\n", m)
			}
			switch m.Method {
			case "connect":
				offer, err := pc.CreateOffer(nil)
				if err != nil {
					f.Close()
					pc.Close()
					conn.Close()
					return
				}
				pc.SetLocalDescription(offer)
				conn.WriteJSON(&Message{"sendto", to, offer})
			case "sendto":
				data := m.Data.(map[string]interface{})
				switch data["type"].(string) {
				case "answer":
					pc.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: data["sdp"].(string)})
				case "candidate":
					pc.AddICECandidate(webrtc.ICECandidateInit{Candidate: data["candidate"].(string)})
				}
			default:
				log.Println("unknown method:", m.Method)
			}
		}
	}()
	conn.WriteJSON(&Message{"connect", "", uuid})

	return nil
}

func main() {
	stuns := []string{"stun:39.108.177.67:13478"}
	signaling := "wss://sv.mi-cloud.cn:5002"
	deviceuuid := "10001-1655461605272-device"
	browseruuid := "10001-1655461605272-browser"
	file := "/usr/local/openwrt_r21.02.01_sdk.tgz"

	err := RunWebRTC(stuns, signaling, deviceuuid, browseruuid, file, 2000*1024*1024, 600*time.Second, false)
	fmt.Println(err)
	select {}
}
