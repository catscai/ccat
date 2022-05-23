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
	TcpCfg       map[string]*TcpServiceCfg `yaml:"tcp_service"`
	IsTcpService bool
}

var AppCfg *GlobalCfg

type TcpServiceCfg struct {
	BaseServiceCfg
	Name        string
	IP          string          `yaml:"ip"`
	Port        uint32          `yaml:"port"`
	MaxConn     uint32          `yaml:"max_conn"`
	WorkerGroup *WorkerGroupCfg `yaml:"worker_group"`
}

type BaseServiceCfg struct {
	MaxPackLen uint32 `yaml:"max_pack_len"` // 最大包长度,0-表示不限制
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
		fmt.Printf("[Config] Tcp info %+v\n", *cfg.TcpCfg[name])
	}
	if len(cfg.TcpCfg) > 0 {
		cfg.IsTcpService = true
	}
	fmt.Printf("[Config] global %+v\n", cfg)
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
		return &cfg.BaseServiceCfg
	}

	return nil
}
