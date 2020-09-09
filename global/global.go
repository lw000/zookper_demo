package global

import (
	"demo/zookper_demo/zkserve"
)

var (
	ZKPServe *zkserve.Zkp
	ZKHosts  = []string{
		"192.168.0.115:2182",
		"192.168.0.115:2183",
		"192.168.0.115:2181",
	}
)

func init() {
	ZKPServe = zkserve.New()
}
