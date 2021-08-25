package settings

// 用viper做读取
import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf *AppConfig

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	Port      int    `mapstructure:"port"`
	StartTime string `mapstructure:"start_time"`
	MachineId int64  `mapstructure:"machine_id"`

	*LogConfig   `mapstructure:"log"`
	*RedisConfig `mapstructure:"redis"`
	*MongoConfig `mapstructure:"mongo"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type MongoConfig struct {
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Database    string `mapstructure:"database"`
	MaxPoolSize int    `mapstructure:"max_pool_size"`
}

func Init() (err error) {
	Conf = &AppConfig{}
	// 方式一： 直接指定配置文件
	viper.SetConfigFile("./conf/conf.yaml") // 直接指定配置文件名称

	// 方式二：指定文件名（不含文件类型），已经对应的加载路径
	// 受pwd影响
	// viper.SetConfigName("conf")    // 配置文件名称（不含文件类型）
	// viper.AddConfigPath("./conf/") // 配置文件路径， 可以添加多个路径

	// 方式三：远程配置中心传入字节流，告诉viper当前的配置使用什么格式去解析
	// viper.SetConfigType("yaml")    // 配置文件类型(只在指定远程获取配置文件时生效)

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		fmt.Printf("viper read config failed: %s \n", err)
		return
	}

	// 读取成功后，反序列化
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper unmarshal config failed: %s \n", err)
		return
	}

	// 持续监听配置
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("config file change...\n")
		// 文件一旦概念，重新unmarshal
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper unmarshal config failed: %s \n", err)
			return
		}
	})
	return
}
