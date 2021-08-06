package serverChat

import (
	"bytes"
	"net"
)

const BuffSize = 128

//续读功能
func ConnReadMany(conn net.Conn) (*bytes.Buffer, error) {
	buff := make([]byte, BuffSize)
	iobuff := new(bytes.Buffer)

	num, err := conn.Read(buff)

	if err != nil {
		return nil, err
	}
	iobuff.Write(buff)
	for num == BuffSize {
		num, err = conn.Read(buff)

		if err != nil {
			return nil, err
		}
		iobuff.Write(buff)
	}

	return iobuff, nil
}
