package cache

func ToCacheMessageValue(m *CacheMessage) map[string]any {
	return map[string]any{
		"id":        m.ID.String(),
		"roomId":    m.RoomID.String(),
		"clientId":  m.ClientID,
		"createdAt": m.Content,
		"content":   m.CreatedAt,
	}
}
