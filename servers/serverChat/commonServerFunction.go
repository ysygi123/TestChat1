package serverChat

import (
	"bytes"
	"net"
)

const BuffSize = 128

//续读功能
func ConnReadMany(conn net.Conn) ([]byte, error) {
	buff := make([]byte, BuffSize)
	iobuff := new(bytes.Buffer)

	allnum := 0

	num, err := conn.Read(buff)
	allnum += num
	if err != nil {
		return nil, err
	}
	iobuff.Write(buff)
	for num == BuffSize {
		num, err = conn.Read(buff)
		allnum += num
		if err != nil && err.Error() != "EOF" {
			return nil, err
		}
		iobuff.Write(buff)
	}
	//必须要这么写，，不然就会出错  invalid character '\x00' after top-level value
	returnByte := iobuff.Bytes()[:allnum]
	return returnByte, nil
}
