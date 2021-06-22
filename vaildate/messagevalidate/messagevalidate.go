package messagevalidate

type Message struct {
	SendUid    int    `json:"send_uid" binding:"required"`
	ReceiveUid int    `json:"receive_uid" binding:"required"`
	Message    string `json:"message" binding:"required"`
}
