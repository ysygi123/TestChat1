package configStart

//配置文件
import (
	"fmt"
	"github.com/spf13/viper"
)

func ConfigStart() {
	viper.SetConfigName("test")
	viper.AddConfigPath("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
