package messagevalidate

type Message struct {
	SendUid        int    `json:"send_uid" binding:"required"`
	ReceiveUid     int    `json:"receive_uid" binding:"required"`
	MessageContent string `json:"message_content" binding:"required"`
}
