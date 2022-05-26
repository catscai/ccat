package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// 默认配置路径
var cfgPath = "./conf/ccat.yaml"

type GlobalCfg struct {
	TcpCfg       map[string]*TcpServiceCfg `yaml:"tcp_service"`
	IsTcpService bool
	LogCfg       *CatLogCfg `yaml:"log"`
}

var AppCfg *GlobalCfg

type TcpServiceCfg struct {
	Base        *BaseServiceCfg
	MaxPackLen  uint32 `yaml:"max_pack_len"` // 最大包长度,0-表示不限制
	Auto        bool   `yaml:"auto"`         // 是否自动创建服务监听
	Name        string
	IP          string          `yaml:"ip"`
	Port        uint32          `yaml:"port"`
	MaxConn     uint32          `yaml:"max_conn"`
	WorkerGroup *WorkerGroupCfg `yaml:"worker_group"`
}

type CatLogCfg struct {
	AppName    string `yaml:"app_name"`
	LogDir     string `yaml:"log_dir"`
	Level      string `yaml:"level"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
	Console    bool   `yaml:"console"`
	IsFuncName bool   `yaml:"func"`
}

type BaseServiceCfg struct {
	MaxPackLen uint32 `yaml:"max_pack_len"` // 最大包长度,0-表示不限制
	Auto       bool   `yaml:"auto"`         // 是否自动创建服务监听
}

type WorkerGroupCfg struct {
	Size        uint32 `yaml:"size"`
	QueueLength uint32 `yaml:"queue_length"`
}

func Reload() error {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		fmt.Println("[Config] Reload config err", err)
		return err
	}
	cfg := GlobalCfg{
		IsTcpService: false,
		LogCfg: &CatLogCfg{
			AppName:    "project",
			LogDir:     "./logs/",
			Level:      "debug",
			MaxSize:    128,
			MaxAge:     7,
			MaxBackups: 30,
			Console:    false,
		},
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Println("[Config] yaml Unmarshal err", err)
		return err
	}
	for name := range cfg.TcpCfg {
		cfg.TcpCfg[name].Name = name
		if cfg.TcpCfg[name].WorkerGroup == nil {
			cfg.TcpCfg[name].WorkerGroup = &WorkerGroupCfg{
				Size:        10,
				QueueLength: 10,
			}
		}
		cfg.TcpCfg[name].Base = &BaseServiceCfg{
			MaxPackLen: cfg.TcpCfg[name].MaxPackLen,
			Auto:       cfg.TcpCfg[name].Auto,
		}
		fmt.Printf("[Config] Tcp info %+v\n", *cfg.TcpCfg[name])
	}
	if len(cfg.TcpCfg) > 0 {
		cfg.IsTcpService = true
	}
	fmt.Printf("[Config] global %+v\nlog:%+v\n", cfg, *cfg.LogCfg)
	AppCfg = &cfg
	return nil
}

func GetTcpServiceCfg(name string) *TcpServiceCfg {
	if _, ok := AppCfg.TcpCfg[name]; ok {
		return AppCfg.TcpCfg[name]
	}

	return nil
}

func GetBaseServiceCfg(name string) *BaseServiceCfg {
	if cfg, ok := AppCfg.TcpCfg[name]; ok {
		return cfg.Base
	}

	return nil
}
