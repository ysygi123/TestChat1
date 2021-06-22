package message

type Message struct {
	Id          int    `db:"id" json:"id"`
	SendUid     int    `db:"send_uid" json:"send_uid"`
	ReceiveUid  int    `db:"receive_uid" json:"receive_uid"`
	CreatedTime uint64 `db:"created_time" json:"created_time"`
	Message     string `db:"message" json:"message"`
}
