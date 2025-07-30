package chat

func (h *Hub) handleRoomRegistrations() {
	for {
		select {
		case r := <-h.roomRegister:
			h.registerRoom(r)
		case r := <-h.roomUnregister:
			h.unregisterRoom(r)
		}
	}
}

func (h *Hub) handleClientRegistrations() {
	for {
		select {
		case c := <-h.clientRegister:
			h.registerClient(c)
		case c := <-h.clientUnregister:
			h.unregisterClient(c)
		}
	}
}
