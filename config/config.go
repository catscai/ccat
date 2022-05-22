package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func init() {
	if err := Reload(); err != nil {
		panic(fmt.Sprintf("[Config] Reload err:%+v", err))
	}
}

// 默认配置路径
var cfgPath = "./conf/ccat.yaml"

type GlobalCfg struct {
	TcpCfg       TcpServiceCfg `yaml:"tcp_service"`
	IsTcpService bool
}

var AppCfg *GlobalCfg

type TcpServiceCfg struct {
	Name       string `yaml:"name"` // 服务名称用于区分多个服务
	IP         string `yaml:"ip"`
	Port       uint32 `yaml:"port"`
	MaxConn    uint32 `yaml:"max_conn"`
	MaxPackLen uint32 `yaml:"max_pack_len"` // 最大包长度,0-表示不限制
}

func Reload() error {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		fmt.Println("[Config] Reload config err", err)
		return err
	}
	cfg := GlobalCfg{
		IsTcpService: false,
		TcpCfg: TcpServiceCfg{
			MaxConn:    1024, // 默认最大连接数
			MaxPackLen: 0,    // 默认最大包长度
		},
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Println("[Config] yaml Unmarshal err", err)
		return err
	}
	if len(cfg.TcpCfg.IP) > 0 && cfg.TcpCfg.Port > 0 {
		cfg.IsTcpService = true
	}
	AppCfg = &cfg
	return nil
}
