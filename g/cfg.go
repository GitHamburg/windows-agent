package g

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/toolkits/file"
	"net"
	"fmt"
)

type PluginConfig struct {
	Enabled bool   `json:"enabled"`
	Dir     string `json:"dir"`
	Git     string `json:"git"`
	LogDir  string `json:"logs"`
	Python  string `json:"python"`
}

type MsSQLConfig struct {
	Enabled  bool     `json:"enabled"`
	Addr     string   `json:"addr"`
	Port     int      `json:"port"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Encrypt  string   `json:"encrypt"`
	Instance []string `json:"instance"`
}

type IIsConfig struct {
	Enabled  bool     `json:"enabled"`
	Websites []string `json:"websites"`
}

type HeartbeatConfig struct {
	Enabled  bool   `json:"enabled"`
	Addr     string `json:"addr"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type TransferConfig struct {
	Enabled  bool     `json:"enabled"`
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
}

type HttpConfig struct {
	Enabled  bool   `json:"enabled"`
	Listen   string `json:"listen"`
	Backdoor bool   `json:"backdoor"`
}

type CollectorConfig struct {
	IfacePrefix []string `json:"ifacePrefix"`
}

type GlobalConfig struct {
	Debug         bool              `json:"debug"`
	Hostname      string            `json:"hostname"`
	IP            string            `json:"ip"`
	Plugin        *PluginConfig     `json:"plugin"`
	IIs           *IIsConfig        `json:"iis"`
	MsSQL         *MsSQLConfig      `json:"mssql"`
	Logfile       string            `json:"logfile"`
	Heartbeat     *HeartbeatConfig  `json:"heartbeat"`
	Transfer      *TransferConfig   `json:"transfer"`
	Http          *HttpConfig       `json:"http"`
	Collector     *CollectorConfig  `json:"collector"`
	DefaultTags   map[string]string `json:"default_tags"`
	IgnoreMetrics map[string]bool   `json:"ignore"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func Hostname() (string, error) {
	hostname := Config().Hostname
	if hostname != "" {
		return hostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		logger.Println("ERROR: os.Hostname() fail", err)
	}

	if hostname == "" || hostname == "__HOSTNAME__" {
		hostname = GetLocalIp()
	}

	if hostname == "" {
		log.Println("本地hostname为空，退出，请手动配置hostname")
		os.Exit(0)
	}

	return hostname, err
}

func IP() string {
	ip := Config().IP
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIps) > 0 {
		ip = LocalIps[0]
	}

	if ip == "" || ip == "__HOSTNAME__" {
		ip = GetLocalIp()
	}else {
		if !CheckLocalIp(ip) {
			log.Println("本地ip不匹配",ip)
			os.Exit(0)
		}
	}

	if ip == "" {
		log.Println("本地ip为空，退出，请手动配置ip")
		os.Exit(0)
	}

	return ip
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}

func GetLocalIp() (string)  {
	if Config().IP != "" {
		return Config().IP
	}else {
		addrs, err := net.InterfaceAddrs()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, address := range addrs {

			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String();
				}

			}
		}
	}
	return "";
}


func CheckLocalIp(ip string) (bool)  {
	if net.ParseIP(ip) == nil {
		//非ip
		return true;
	}else {
		addrs, err := net.InterfaceAddrs()

		if err != nil {
			fmt.Println(err)
			return true;
		}

		for _, address := range addrs {

			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil && ipnet.IP.String() == ip {
					return true;
				}

			}
		}
	}
	return false;
}