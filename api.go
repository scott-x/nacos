/*
* @Author: scottxiong
* @Date:   2021-10-20 17:51:18
* @Last Modified by:   scottxiong
* @Last Modified time: 2021-10-20 21:17:28
 */
package nacos

import (
	"errors"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	configClient config_client.IConfigClient //nacos client
)

//nacos config
type Config struct {
	Ip          string
	Port        int
	NamespaceId string
	groups      *[]Group
}

type Group struct {
	DataId  string //data id
	Name    string //group name
	content string //the content of the config file (string)
}

//set config
func (g *Group) setConf() {
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: g.DataId,
		Group:  g.Name,
	})

	if err != nil {
		panic(err)
	}

	g.content = content
}

func (g *Group) GetConfig() string {
	if len(g.content) == 0 {
		g.setConf()
	}
	return g.content
}

func (c *Config) NewGroup(name, id string) *Group {
	groups := c.groups

	for _, v := range *groups {
		if v.DataId == id {
			panic(errors.New("error: duplicate DataId " + id))
		}
	}

	g := &Group{
		DataId: id,
		Name:   name,
	}

	*groups = append(*groups, *g)

	go watch(g) //watch config

	return g
}

func watch(g *Group) {
	//listen config change
	for {
		time.Sleep(time.Second)
		err := configClient.ListenConfig(vo.ConfigParam{
			DataId: g.DataId,
			Group:  g.Name,
			OnChange: func(namespace, group, dataId, data string) {
				g.content = data
			},
		})
		if err != nil {
			panic(err)
		}
	}
}

func initClient(c *Config) {
	var err error

	ss := []constant.ServerConfig{
		{
			IpAddr: c.Ip,           // the nacos server address
			Port:   uint64(c.Port), // the nacos server port
		},
	}

	cs := constant.ClientConfig{
		NamespaceId:         c.NamespaceId, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err = clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": ss,
		"clientConfig":  cs,
	})

	if err != nil {
		panic(err)
	}
}

func InitConfig(ip string, port int, namespaceId string) *Config {
	//set up Config
	c := &Config{
		Ip:          ip,
		Port:        port,
		NamespaceId: namespaceId,
		groups:      new([]Group), //important: must be init, otherwise will throw nil pointer error
	}

	//init client
	initClient(c)

	return c
}
