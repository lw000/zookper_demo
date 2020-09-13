package global

var (
	ZookeeperHosts = []string{
		"192.168.0.115:2182",
		"192.168.0.115:2183",
		"192.168.0.115:2181",
	}

	SharedData = make(map[string]interface{})
)

func init() {

}
