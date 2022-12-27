package app

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/wgo-admin/backend/internal/app/store"
	"github.com/wgo-admin/backend/internal/pkg/log"
	"github.com/wgo-admin/backend/pkg/db"
)

// 初始化配置，从配置文件、环境变量，并读取配置文件的内容到viper
func initConfig() {
	fmt.Println("cfg", cfgFile)
	if cfgFile != "" {
		// 从命令行参数获取配置文件路径
		viper.SetConfigFile(cfgFile)
	} else {

	}

	// 读取匹配的环境变量
	viper.AutomaticEnv()

	// 读取环境变量的前缀为 MINIBLOG，如果是 miniblog，将自动转变为大写。
	viper.SetEnvPrefix("MINIBLOG")

	// 以下 2 行，将 viper.Get(key) key 字符串中 '.' 和 '-' 替换为 '_'
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 读取配置文件。如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Failed to read viper configuration file", "err", err)
	}

	// 打印 viper 当前使用的配置文件，方便 Debug.
	log.Debugw("Using config file", "file", viper.ConfigFileUsed())
}

func logConfig() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            log.FormatLog(viper.GetString("log.format")),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}

func initStore() error {
	dbOptions := &db.MySQLOptions{
		Host:                  viper.GetString("db.host"),
		Username:              viper.GetString("db.username"),
		Password:              viper.GetString("db.password"),
		Database:              viper.GetString("db.database"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
		LogLevel:              viper.GetInt("db.log-level"),
	}

	db, err := db.NewMySQL(dbOptions)
	if err != nil {
		return err
	}

	// 初始化存储层
	_ = store.NewStore(db)

	return nil
}
