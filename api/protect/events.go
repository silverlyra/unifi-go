package protect

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type UpdateReceiver struct {
	C    chan ReceivedUpdate
	conn *websocket.Conn
	stop chan struct{}
	done chan struct{}
}

func (c *Client) ReceiveUpdates(ctx context.Context) (*UpdateReceiver, error) {
	u := c.URL("/proxy/protect/ws/updates")
	u.Scheme = "wss"

	if c.Bootstrap == nil {
		if err := c.Load(ctx); err != nil {
			return nil, err
		}
	}
	q := u.Query()
	q.Set("lastUpdateId", c.Bootstrap.LastUpdateID)
	u.RawQuery = q.Encode()

	log.Printf("Receiving updates via %v", u)
	conn, _, err := c.WebSocket.DialContext(ctx, u.String(), nil)
	if err != nil {
		return nil, err
	}

	receiver := newUpdateReceiver(conn)
	go receiver.attend()

	return receiver, nil
}

func newUpdateReceiver(conn *websocket.Conn) *UpdateReceiver {
	return &UpdateReceiver{
		conn: conn,
		C:    make(chan ReceivedUpdate, 1),
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
}

func (r *UpdateReceiver) Stop() {
	close(r.stop)
}

func (r *UpdateReceiver) attend() {
	go r.receive()

	for {
		select {
		case <-r.done:
			return
		case <-r.stop:
			r.conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)

			select {
			case <-r.done:
			case <-time.After(time.Second):
				log.Println("force-closing WebSocket connection")
				r.conn.Close()
			}
			return
		}
	}
}

func (r *UpdateReceiver) receive() {
	defer close(r.done)
	defer close(r.C)

	for {
		_, msg, err := r.conn.ReadMessage()
		if err != nil {
			log.Printf("error receiving update from WebSocket: %v\n", err)
			return
		}

		update, err := readUpdate(msg)
		if err != nil {
			log.Printf("received invalid update from WebSocket: %v\n", err)
			continue
		}

		r.C <- update
	}
}

type ReceivedUpdate struct {
	Action  *Action
	Payload interface{}
}

type Action struct {
	ID          string `json:"id"`
	Action      string `json:"action"`
	ModelKey    string `json:"modelKey"`
	NewUpdateID string `json:"newUpdateId"`
}

type CameraUpdate struct {
	Camera *Camera

	IsMotionDetected bool `json:"isMotionDetected"`
	LastMotion       Time `json:"lastMotion"`
	LastRing         Time `json:"lastRing"`
}

type Event struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	Camera *Camera

	CameraID          string   `json:"camera"`
	ModelKey          string   `json:"modelKey"`
	Start             Time     `json:"start"`
	Score             float64  `json:"score"`
	SmartDetectEvents []string `json:"smartDetectEvents"`
	SmartDetectTypes  []string `json:"smartDetectTypes"`
}

const (
	updateFrameAction  byte = 1
	updateFramePayload byte = 2
)

func readUpdate(data []byte) (ReceivedUpdate, error) {
	var action Action
	var payload interface{}

	data, err := readUpdateFrame(updateFrameAction, data, &action)
	if err != nil {
		return ReceivedUpdate{}, err
	}

	if action.ModelKey == "camera" && action.Action == "update" {
		payload = &CameraUpdate{}
	} else if action.ModelKey == "event" {
		payload = &Event{}
	} else {
		m := make(map[string]interface{})
		payload = &m
	}

	_, err = readUpdateFrame(updateFramePayload, data, payload)
	return ReceivedUpdate{Action: &action, Payload: payload}, err
}

func readUpdateFrame(kind byte, data []byte, dest interface{}) ([]byte, error) {
	if len(data) < 8 {
		return data, fmt.Errorf("data smaller than 8-byte frame (got %d bytes)", len(data))
	}
	if data[0] != kind {
		return data, fmt.Errorf("expected %d frame; got %d", kind, data[0])
	}
	if data[1] != 1 {
		return data, fmt.Errorf("expected JSON frame; got %d", data[1])
	}
	if data[2] > 1 {
		return data, fmt.Errorf("unexpected frame compression %d", data[2])
	}

	deflated := data[2] == 1

	size := binary.BigEndian.Uint32(data[4:])
	contents := data[8:]
	payload := contents[:size]
	rest := contents[size:]

	if deflated {
		zr, err := zlib.NewReader(bytes.NewReader(payload))
		if err != nil {
			return data, err
		}
		defer zr.Close()

		var b bytes.Buffer
		if _, err := io.Copy(&b, zr); err != nil {
			return rest, err
		}

		payload = b.Bytes()
	}

	if err := json.Unmarshal(payload, dest); err != nil {
		return rest, err
	}

	return rest, nil
}
