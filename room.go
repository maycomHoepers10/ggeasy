package main

import (
	"fmt"
	"math/rand"
)

type Room struct {
	Name      string
	Clients   map[*Client]bool
	Broadcast chan []byte
}

func generate(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, n)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func NewRoom(client *Client) *Room {

	room := &Room{
		Name:      string(generate(10)),
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan []byte),
	}

	room.Clients[client] = true //Adiciona o cliente na lista de clientes da SALA

	return room
}

/*Vai entrar em uma sala aleatória ou criar sala*/
func EnterRoom(client *Client, h *Hub) {
	fmt.Println("Enter room!")

	var room *Room
	fmt.Println(room)
	fmt.Println(len(h.rooms) == 0)

	if len(h.rooms) == 0 {
		fmt.Println("Registrando a primeira sala!")
		room = NewRoom(client)
	} else {
		fmt.Println(h.rooms)
		for rm := range h.rooms {

			if len(rm.Clients) < 2 {
				fmt.Println("Entrou em uma sala existente!", string(rm.Name))
				rm.Clients[client] = true
				room = rm
				break
			}
		}
	}

	if room == nil {
		fmt.Println("Registrando sala já que todas estão cheias!")
		room = NewRoom(client)
	}

	client.room = room
	h.rooms[room] = true

	//start mensagens da sala
	go room.broadcast()
}

func (r *Room) broadcast() {

	for {
		select {
		case message := <-r.Broadcast:
			println("Enviando mensagem do cliente")
			for client := range r.Clients {
				select {
				case client.send <- message:
					println("enviou mensagem para pessoa da sala...")
				default:
					close(client.send)
					delete(r.Clients, client)
				}
			}
		}
	}
}
