package messagevalidate

type Message struct {
	MessageType    uint8  `db:"message_type" json:"message_type"`
	SendUid        int    `json:"send_uid" binding:"required"`
	ReceiveUid     int    `json:"receive_uid" binding:"required"`
	GroupId        int    `json:"group_id"`
	MessageContent string `json:"message_content" binding:"required"`
}
