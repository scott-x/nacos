# nacos tool

对nacos进行封装，无须再手动监听就能拿到动态变化的配置, 同时也去掉了nacos自带的部分日志


### API

- `func InitConfig(ip string, port int, namespaceId string) *Config`
- `func (c *Config) NewGroup(name, id string) *Group`
- `func (g *Group) GetConfig() string`
- `func PublishConfig(g *Group, conf string) (bool, error)`
- `func DeleteConfig(g *Group) (bool, error)`

### Get Started

```go
package main

import (
    "fmt"
    "github.com/scott-x/nacos"
)

func main() {
    c := nacos.InitConfig("www.google.com", 10086, "3b1ae4c9-4895-4693-802f-21991b67f322")
    g1 := c.NewGroup("group", "wifi.json")
    g2 := c.NewGroup("group", "brand.json")
    for {
        fmt.Println(g1.GetConfig())
        fmt.Println(g2.GetConfig())
    }  
}
```