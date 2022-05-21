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
	IP      string `yaml:"ip"`
	Port    uint32 `yaml:"port"`
	MaxConn uint32 `yaml:"max_conn"`
}

func Reload() error {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		fmt.Println("[Config] Reload config err", err)
		return err
	}
	cfg := GlobalCfg{
		IsTcpService: false,
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
