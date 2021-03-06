package main

import (
	"encoding/base64"
	"encoding/json"
	"hash/crc32"
	"io/ioutil"
	"net/url"
	"strconv"
)

type DataSource struct {
	Name         string              `json:"name"`
	SubOption    SubOpt              `json:"sub_option"`
	UpdateOption UpdOpt              `json:"update_option"`
	Shadowsocks  []ShadowsocksEntity `json:"shadowsocks"`
}

func (ds DataSource) GetShadowsocksNames() []string {
	ret := make([]string, len(ds.Shadowsocks))
	for i, v := range ds.Shadowsocks {
		ret[i] = v.Name
	}
	return ret
}

type SubOpt struct {
	Path []string `json:"path"`
}

type UpdOpt struct {
	Path string `json:"path"`
}

type ShadowsocksEntity struct {
	Name       string `json:"name"`
	Server     string `json:"server"`
	Port       int    `json:"port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
	Plugin     string `json:"plugin"`
	PluginOpts string `json:"plugin_opts"`
}

func (se ShadowsocksEntity) ToURI() string {
	uri := "ss://"
	uri += base64.RawURLEncoding.EncodeToString([]byte(se.Method + ":" + se.Password))
	uri += "@"
	uri += se.Server
	uri += ":"
	uri += strconv.Itoa(se.Port)
	if se.Plugin != "" {
		uri += "?plugin="
		uri += url.PathEscape(se.Plugin + ";" + se.PluginOpts)
	}
	uri += "#"
	uri += url.PathEscape(se.Name)
	return uri
}
func (se ShadowsocksEntity) ServerID() string {
	return strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(se.Server+":"+strconv.Itoa(se.Port)))), 16)
}

var datasource DataSource

func LoadDataSource(filepath string) error {
	bJSON, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bJSON, &datasource)
	if err != nil {
		return err
	}
	return nil
}
