package main

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v2"
)

var registry = map[string]func(DataSource) string{
	"plain":  plainSubs,
	"base64": base64Subs,
	"sip008": sip008Subs,
	"clash":  clashSubs,
}

func plainSubs(ds DataSource) string {
	ssUri := make([]string, len(ds.Shadowsocks))
	for i, v := range ds.Shadowsocks {
		ssUri[i] = v.ToURI()
	}
	return strings.Join(ssUri, "\n")
}

func base64Subs(ds DataSource) string {
	return base64.StdEncoding.EncodeToString([]byte(plainSubs(ds)))
}

type SIP008 struct {
	Version int            `json:"version"`
	Remark  string         `json:"remark"`
	Servers []SIP008Server `json:"servers"`
}
type SIP008Server struct {
	ID         string `json:"id"`
	Remark     string `json:"remark"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
	Plugin     string `json:"plugin,omitempty"`
	PluginOpts string `json:"plugin_opts,omitempty"`
}

func sip008Subs(ds DataSource) string {
	servers := make([]SIP008Server, len(ds.Shadowsocks))
	for i, v := range ds.Shadowsocks {
		servers[i] = SIP008Server{
			ID:         v.ServerID(),
			Remark:     v.Name,
			Server:     v.Server,
			ServerPort: v.Port,
			Password:   v.Password,
			Method:     v.Method,
			Plugin:     v.Plugin,
			PluginOpts: v.PluginOpts,
		}
	}
	sip008 := SIP008{
		Version: 1,
		Remark:  ds.Name,
		Servers: servers,
	}
	bjson, err := json.Marshal(sip008)
	if err != nil {
		return ""
	} else {
		return string(bjson)
	}
}

type ClashFile struct {
	Proxies       []Proxies          `yaml:"proxies"`
	ProxyGroups   []ProxyGroups      `yaml:"proxy-groups"`
	RuleProviders map[string]Ruleset `yaml:"rule-providers"`
	Rules         []string           `yaml:"rules"`
}
type PluginOpts struct {
	Mode           string `yaml:"mode"`
	TLS            bool   `yaml:"tls"`
	SkipCertVerify bool   `yaml:"skip-cert-verify"`
	Host           string `yaml:"host"`
	Path           string `yaml:"path"`
	Mux            bool   `yaml:"mux"`
}
type Proxies struct {
	Name       string     `yaml:"name"`
	Type       string     `yaml:"type"`
	Server     string     `yaml:"server"`
	Port       int        `yaml:"port"`
	Cipher     string     `yaml:"cipher"`
	Password   string     `yaml:"password"`
	Plugin     string     `yaml:"plugin,omitempty"`
	PluginOpts PluginOpts `yaml:"plugin-opts,omitempty"`
}
type ProxyGroups struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Proxies   []string `yaml:"proxies"`
	Tolerance int      `yaml:"tolerance"`
	URL       string   `yaml:"url"`
	Interval  int      `yaml:"interval"`
}
type Ruleset struct {
	Behavior string `yaml:"behavior"`
	Type     string `yaml:"type"`
	URL      string `yaml:"url"`
	Interval int    `yaml:"interval"`
	Path     string `yaml:"path"`
}

func clashSubs(ds DataSource) string {
	clash := ClashFile{
		ProxyGroups: []ProxyGroups{
			{
				Name:      "AutoSelect",
				Type:      "url-test",
				Proxies:   ds.GetShadowsocksNames(),
				Tolerance: 100,
				URL:       "http://www.gstatic.com/generate_204",
				Interval:  600,
			},
		},
		RuleProviders: map[string]Ruleset{
			"whitelist_all_in_one": {
				Behavior: "domain",
				Type:     "http",
				URL:      "https://raw.githubusercontent.com/IceCodeNew/4Share/master/Clash/DOMAIN_DIRECT_ALLINONE.yaml",
				Interval: 10800,
				Path:     "custom_rules/DOMAIN_DIRECT_ALLINONE.yaml",
			},
			"blacklist_all_in_one": {
				Behavior: "domain",
				Type:     "http",
				URL:      "https://raw.githubusercontent.com/IceCodeNew/4Share/master/Clash/DOMAIN_ADS_ALLINONE.yaml",
				Interval: 10800,
				Path:     "custom_rules/DOMAIN_ADS_ALLINONE.yaml",
			},
			"recommend_to_proxy": {
				Behavior: "domain",
				Type:     "http",
				URL:      "https://raw.githubusercontent.com/IceCodeNew/4Share/master/Clash/ICN/DOMAIN_ICN_PROXY.yaml",
				Interval: 10800,
				Path:     "custom_rules/DOMAIN_ICN_PROXY.yaml",
			},
			"speedtest": {
				Behavior: "domain",
				Type:     "http",
				URL:      "https://raw.githubusercontent.com/IceCodeNew/4Share/master/Clash/ICN/DOMAIN_ICN_SPEEDTEST.yaml",
				Interval: 10800,
				Path:     "custom_rules/DOMAIN_ICN_SPEEDTEST.yaml",
			},
			"china_ip_list": {
				Behavior: "domain",
				Type:     "http",
				URL:      "https://raw.githubusercontent.com/IceCodeNew/4Share/master/Clash/ICN/CHINA_IP_LIST.yaml",
				Interval: 10800,
				Path:     "custom_rules/CHINA_IP_LIST.yaml",
			},
		},
		Rules: []string{
			"RULE-SET,blacklist_all_in_one,REJECT",
			"RULE-SET,recommend_to_proxy,AutoSelect",
			"RULE-SET,whitelist_all_in_one,DIRECT",
			"DOMAIN-SUFFIX,cn,DIRECT",
			"DOMAIN-SUFFIX,ip6-localhost,DIRECT",
			"DOMAIN-SUFFIX,ip6-loopback,DIRECT",
			"DOMAIN-SUFFIX,lan,DIRECT",
			"DOMAIN-SUFFIX,local,DIRECT",
			"DOMAIN-SUFFIX,localhost,DIRECT",
			"IP-CIDR,0.0.0.0/8,DIRECT",
			"IP-CIDR,10.0.0.0/8,DIRECT",
			"IP-CIDR,100.64.0.0/10,DIRECT",
			"IP-CIDR,127.0.0.0/8,DIRECT",
			"IP-CIDR,169.254.0.0/16,DIRECT",
			"IP-CIDR,172.16.0.0/12,DIRECT",
			"IP-CIDR,192.0.0.0/24,DIRECT",
			"IP-CIDR,192.0.2.0/24,DIRECT",
			"IP-CIDR,192.88.99.0/24,DIRECT",
			"IP-CIDR,192.168.0.0/16,DIRECT",
			"IP-CIDR,198.18.0.0/15,DIRECT",
			"IP-CIDR,198.51.100.0/24,DIRECT",
			"IP-CIDR,203.0.113.0/24,DIRECT",
			"IP-CIDR,224.0.0.0/4,DIRECT",
			"IP-CIDR,240.0.0.0/4,DIRECT",
			"IP-CIDR,255.255.255.255/32,DIRECT",
			"IP-CIDR,::1/128,DIRECT",
			"IP-CIDR,fc00::/7,DIRECT",
			"IP-CIDR,fe80::/10,DIRECT",
			"RULE-SET,china_ip_list,DIRECT",
			"GEOIP,CN,DIRECT",
			"IP-CIDR,115.27.0.0/16,DIRECT",
			"IP-CIDR,162.105.0.0/16,DIRECT",
			"IP-CIDR,202.112.7.0/24,DIRECT",
			"IP-CIDR,202.112.8.0/24,DIRECT",
			"IP-CIDR,222.29.0.0/17,DIRECT",
			"IP-CIDR,222.29.128.0/19,DIRECT",
			"IP-CIDR,2001:da8:201::/48,DIRECT",
			"MATCH,AutoSelect",
		},
	}
	clash.Proxies = make([]Proxies, len(ds.Shadowsocks))
	for i, v := range ds.Shadowsocks {
		proxy := Proxies{
			Name:     v.Name,
			Type:     "ss",
			Server:   v.Server,
			Port:     v.Port,
			Cipher:   v.Method,
			Password: v.Password,
		}
		if v.Plugin == "v2ray-plugin" {
			var pluginHost, pluginPath string
			for _, v := range strings.Split(v.PluginOpts, ";") {
				if strings.HasPrefix(v, "host=") {
					pluginHost = v[5:]
				}
				if strings.HasPrefix(v, "path=") {
					pluginPath = v[5:]
				}
			}
			proxy.Plugin = "v2ray-plugin"
			proxy.PluginOpts = PluginOpts{
				Mode:           "websocket",
				TLS:            true,
				SkipCertVerify: false,
				Mux:            false,
				Host:           pluginHost,
				Path:           pluginPath,
			}
		}
		clash.Proxies[i] = proxy
	}
	byaml, err := yaml.Marshal(clash)
	if err != nil {
		return ""
	} else {
		return string(byaml)
	}
}

func GenerateSubscription(typeName string) string {
	if fn, ok := registry[typeName]; ok {
		return fn(datasource)
	} else {
		return ""
	}
}
