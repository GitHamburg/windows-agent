package g

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/toolkits/file"
	"net"
	"fmt"
	"os/exec"
	"io/ioutil"
	"strings"
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
	if hostname != "" && hostname != "__HOSTNAME__" {
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
		log.Println("本地hostname为空，退出，请手动配置hostname！")
		os.Exit(0)
	}

	return hostname, err
}

func IP() string {
	ip := Config().IP
	if ip != "" && ip != "__HOSTNAME__" {
		// use ip in configuration

		if !CheckLocalIp(ip) {
			log.Println("本地ip不匹配",ip,"准备退出！")
			os.Exit(0)
		}

		return ip
	}

	if len(LocalIps) > 0 {
		ip = LocalIps[0]
	}

	if ip == "" || ip == "__HOSTNAME__" {
		ip = GetLocalIp()
	}

	if ip == "" {
		log.Println("本地ip为空，退出，请手动配置ip！")
		os.Exit(0)
	}

	if len(LocalIps) == 0 {
		LocalIps = append(LocalIps, ip)
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
	if Config().IP != "" && Config().IP != "__HOSTNAME__" {
		return Config().IP
	}else {
		var finalIp string
		cmd := exec.Command("cmd","/c","ipconfig")
		if out, err := cmd.StdoutPipe(); err != nil {
			fmt.Println(err)
		}else{
			defer out.Close()
			if err := cmd.Start(); err != nil {
				fmt.Println(err)
			}

			if opBytes, err := ioutil.ReadAll(out); err != nil {
				log.Fatal(err)
			} else {
				str:= string(opBytes)

				var strs = strings.Split(str,"\r\n")

				if 0 != len(strs){
					var havingFinalIp4 bool = false
					var cnt int = 0
					for index,value := range strs{
						vidx := strings.Index(value,"IPv4")
						//说明已经找到该ip
						if vidx != -1{
							ip4lines := strings.Split(value,":")
							if len(ip4lines) == 2{
								cnt = index
								havingFinalIp4 = true
								ip4str := ip4lines[1]
								finalIp = strings.TrimSpace(ip4str)
							}

						}
						if havingFinalIp4 && index == cnt +2{
							lindex := strings.Index(value,":")
							if -1 != lindex{
								lines := strings.Split(value,":")
								if len(lines) == 2{
									fip := lines[1]
									if strings.TrimSpace(fip) != ""{
										break
									}
								}
							}
							havingFinalIp4 = false
							finalIp = ""
						}
					}
				}
			}
		}
		return finalIp
	}
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