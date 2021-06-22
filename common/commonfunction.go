package common

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"time"
)

func GetRequestParams(c *gin.Context) (map[string]interface{}, error) {
	json := make(map[string]interface{})
	err := c.ShouldBindJSON(&json)
	fmt.Println(json, err)
	if err != nil {
		return nil, err
	}
	return json, nil
}

func AutoValidate(c *gin.Context, s interface{}) error {
	if err := c.ShouldBindJSON(s); err != nil {
		return err
	}
	return nil
}

func GetMD5Data(data string) string {
	byteData := []byte(data)
	m := md5.New()
	m.Write(byteData)
	md5str := hex.EncodeToString(m.Sum(nil))
	return string(md5str)
}

//生成随机prefix的session
func GetSession(prefix string) string {
	t := time.Now().Unix()
	rand.Seed(t)
	randNum := rand.Int63()
	return GetMD5Data(prefix + strconv.Itoa(int(randNum+t)))
}

//简单判断websocket收到的参数是不是空的
func CheckWebSocketParamsIsUnEmpty(keyNames []string, m *map[string]interface{}) error {
	for _, k := range keyNames {
		_, ok := (*m)[k]
		if ok == false {
			err := errors.New("查无key : " + k)
			return err
		}
	}
	return nil
}
