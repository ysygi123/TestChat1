package common

import (
	"crypto/md5"
	"encoding/hex"
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

func GetSession(prefix string) string {
	t := time.Now().Unix()
	rand.Seed(t)
	randNum := rand.Int63()
	return GetMD5Data(prefix + strconv.Itoa(int(randNum+t)))
}
