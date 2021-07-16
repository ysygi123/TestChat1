package message

//这个是用于丢到队列里面去 处理多种消息用的
type PipelineMessage struct {
	MessageType uint8       `json:"message_type"`
	MessageBody interface{} `json:"message_body"`
}

type Message struct {
	IsDel          uint8  `db:"is_del" json:"is_del"`
	MessageType    uint8  `db:"message_type" json:"message_type"`
	Id             int    `db:"id" json:"id"`
	SendUid        int    `db:"send_uid" json:"send_uid"`
	ReceiveUid     int    `db:"receive_uid" json:"receive_uid"`
	GroupId        int    `db:"group_id" json:"group_id"`
	CreatedTime    uint64 `db:"created_time" json:"created_time"`
	ChatId         uint64 `db:"chat_id" json:"chat_id"`
	MessageContent string `db:"message_content" json:"message_content"`
}

type MessageList struct {
	MessageType    uint8  `db:"message_type" json:"message_type"`
	IsDel          uint8  `db:"is_del" json:"is_del"`
	MessageNum     int    `db:"message_num" json:"message_num"`
	Id             int    `db:"id" json:"id"`
	Uid            int    `db:"uid" json:"uid"`
	FromId         int    `db:"from_id" json:"from_id"`
	CreatedTime    uint64 `db:"created_time" json:"created_time"`
	UpdateTime     uint64 `db:"update_time" json:"update_time"`
	ChatId         uint64 `db:"chat_id" json:"chat_id"`
	MessageContent string `db:"message_content" json:"message_content"`
}
