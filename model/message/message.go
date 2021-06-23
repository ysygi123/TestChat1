package message

type Message struct {
	Id             int    `db:"id" json:"id"`
	SendUid        int    `db:"send_uid" json:"send_uid"`
	ReceiveUid     int    `db:"receive_uid" json:"receive_uid"`
	CreatedTime    uint64 `db:"created_time" json:"created_time"`
	MessageContent string `db:"message_content" json:"message_content"`
}

type MessageList struct {
	MessageType    uint8  `db:"message_type" json:"message_type"`
	MessageNum     uint8  `db:"message_num" json:"message_num"`
	IsDel          uint8  `db:"is_del" json:"is_del"`
	Id             int    `db:"id" json:"id"`
	Uid            int    `db:"uid" json:"uid"`
	FromId         int    `db:"from_id" json:"from_id"`
	CreatedTime    uint64 `db:"created_time" json:"created_time"`
	UpdateTime     uint64 `db:"created_time" json:"created_time"`
	MessageContent string `db:"message_content" json:"message_content"`
}