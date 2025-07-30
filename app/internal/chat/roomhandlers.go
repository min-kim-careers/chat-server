package chat

func (r *Room) handleClientRegistrations() {
	for {
		select {
		case <-r.ctx.Done():
			return
		case c := <-r.clientRegister:
			r.registerClient(c)
		case c := <-r.clientUnregister:
			r.unregisterClient(c)
		}
	}
}

func (r *Room) handleMessages() {
	defer r.pubsub.Close()

	for {
		select {
		case <-r.ctx.Done():
			return
		case m := <-r.pubsub.Channel():
			p := []byte(m.Payload)
			for _, c := range r.clients {
				if c.outbound != nil {
					c.outbound <- p
				}
			}
		}
	}
}

func (r *Room) handleClose() {
	<-r.ctx.Done()
	r.pubsub.Close()
	close(r.clientRegister)
	r.clientRegister = nil
	close(r.clientUnregister)
	r.clientUnregister = nil
}
