// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	rooms map[*Room]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[*Room]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	println("Server started...")
	for {
		select {
		case client := <-h.register:
			println("Registrando cliente em uma sala")
			println(client.name)
			h.clients[client] = true
			EnterRoom(client, h)
			println("Está é minha sala...")
			println(client.room.Name)
		case client := <-h.unregister:
			println("Descadastrando cliente")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			println("Enviando mensagem do cliente")
			for client := range h.clients {
				select {
				case client.send <- message:
					println("send...")
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
