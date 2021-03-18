package main

import (
	"LiveServer/common"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
)

// Task Hub任务
type Task struct {
	client *Client
	data   []byte
}

// Hub 管理所有客户端连接
type Hub struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	task       chan *Task
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
		task:       make(chan *Task, 16),
	}
}

func makeRoomID(fromPhone string, fromDeviceID string, toPhone string, toDeviceID string) string {
	h := md5.New()
	h.Write([]byte(fromPhone + fromDeviceID + toPhone + toDeviceID))
	roomID := hex.EncodeToString(h.Sum(nil))
	return roomID
}

func makeClientID(phone string, deviceID string) string {
	return phone + "_" + deviceID
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			id := makeClientID(client.phone, client.deviceID)
			client.id = id
			h.clients[id] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		case broadcast := <-h.broadcast:
			for _, client := range h.clients {
				select {
				case client.send <- broadcast:
				default:
					close(client.send)
					delete(h.clients, client.id)
				}
			}
		case task := <-h.task:
			h.handleMessage(task.client, task.data)
		}
	}
}

func (h *Hub) onCall(trigger *Client, jsonMsg *common.JsonH) {
	toPhone, ok := (*jsonMsg)["to_phone"].(string)
	if !ok || len(toPhone) == 0 {
		sendData := common.JsonH{
			"type":    "call_failed",
			"message": "to_phone is failed to parse",
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
		return
	}
	toDeviceID, ok := (*jsonMsg)["to_device_id"].(string)
	if !ok || len(toDeviceID) == 0 {
		sendData := common.JsonH{
			"type":    "call_failed",
			"message": "to_device_id is failed to parse",
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
		return
	}

	toID := makeClientID(toPhone, toDeviceID)
	toClient, ok := h.clients[toID]
	if ok {
		roomID := makeRoomID(trigger.phone, trigger.deviceID, toPhone, toDeviceID)
		sendData := common.JsonH{
			"type":           "called",
			"from_phone":     trigger.phone,
			"from_device_id": trigger.deviceID,
			"room_id":        roomID,
		}
		bSendData, _ := json.Marshal(sendData)
		toClient.send <- bSendData
	} else {
		sendData := common.JsonH{
			"type":         "call_failed",
			"to_phone":     toPhone,
			"to_device_id": toDeviceID,
			"message":      "to_phone is not online",
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
	}
	return
}

func (h *Hub) onCalledResponse(messageType string, trigger *Client, jsonMsg *common.JsonH) {
	fromPhone, ok := (*jsonMsg)["from_phone"].(string)
	if !ok || len(fromPhone) == 0 {
		sendData := common.JsonH{
			"type":    "called_failed",
			"message": "from_phone is failed to parse",
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
		return
	}
	fromDeviceID, ok := (*jsonMsg)["from_device_id"].(string)
	if !ok || len(fromDeviceID) == 0 {
		sendData := common.JsonH{
			"type":    "called_failed",
			"message": "from_device_id is failed to parse",
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
		return
	}

	roomID, ok := (*jsonMsg)["room_id"].(string)
	if !ok || len(roomID) == 0 {
		sendData := common.JsonH{
			"type":    "called_failed",
			"message": "room_id is failed to parse",
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
		return
	}

	fromID := makeClientID(fromPhone, fromDeviceID)
	fromClient, ok := h.clients[fromID]
	if ok {
		sendData := common.JsonH{
			"type":         messageType,
			"to_phone":     trigger.phone,
			"to_device_id": trigger.deviceID,
			"room_id":      roomID,
		}
		bSendData, _ := json.Marshal(sendData)
		fromClient.send <- bSendData
	} else {
		sendData := common.JsonH{
			"type":           "called_failed",
			"from_phone":     fromPhone,
			"from_device_id": fromDeviceID,
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
	}
	return
}

func (h *Hub) onCallingHungup(trigger *Client, jsonMsg *common.JsonH) {
	notifiedPhone, ok := (*jsonMsg)["notified_phone"].(string)
	if !ok || len(notifiedPhone) == 0 {
		sendData := common.JsonH{
			"type":    "calling_hungup_error",
			"message": "notified_phone is failed to parse",
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
		return
	}
	notifiedDeviceID, ok := (*jsonMsg)["notified_device_id"].(string)
	if !ok || len(notifiedDeviceID) == 0 {
		sendData := common.JsonH{
			"type":    "calling_hungup_error",
			"message": "notified_device_id is failed to parse",
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
		return
	}
	roomID, ok := (*jsonMsg)["room_id"].(string)
	if !ok || len(roomID) == 0 {
		sendData := common.JsonH{
			"type":    "calling_hungup_error",
			"message": "room_id is failed to parse",
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
		return
	}

	notifiedID := makeClientID(notifiedPhone, notifiedDeviceID)
	notifiedClient, ok := h.clients[notifiedID]
	if ok {
		sendData := common.JsonH{
			"type":           "calling_be_hangup",
			"fire_phone":     trigger.phone,
			"fire_device_id": trigger.deviceID,
			"room_id":        roomID,
		}
		bSendData, _ := json.Marshal(sendData)
		notifiedClient.send <- bSendData
	}
}

// PushTask 推送消息
func (h *Hub) PushTask(trigger *Client, message []byte) {
	h.task <- &Task{client: trigger, data: message}
}

func (h *Hub) handleMessage(trigger *Client, message []byte) {
	var jsonMsg common.JsonH
	err := json.Unmarshal(message, &jsonMsg)
	if err != nil {
		fmt.Println(err.Error())
	}

	msg, ok := jsonMsg["type"]
	if !ok {
		sendData := common.JsonH{
			"type":    "json_data_invalid",
			"message": "the requested data is invalid",
		}
		bSendData, _ := json.Marshal(sendData)
		trigger.send <- bSendData
		return
	}

	switch msg.(string) {
	case "call":
		h.onCall(trigger, &jsonMsg)
	case "called_success":
		h.onCalledResponse("called_success", trigger, &jsonMsg)
	case "called_refuse":
		h.onCalledResponse("call_refuse", trigger, &jsonMsg)
	case "calling_hangup":
		h.onCallingHungup(trigger, &jsonMsg)
	default:
		log.Printf("handleMessage:%s is undefine", msg.(string))
	}
}
