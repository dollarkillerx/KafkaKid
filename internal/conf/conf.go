package conf

import (
	"encoding/json"
	"io/ioutil"
)

type KidConfig struct {
	Nodes     []NodeConfig      `json:"Nodes" yaml:"Nodes"`
	Zookeeper []ZookeeperConfig `json:"Zookeeper" yaml:"Zookeeper"`
	Kafka     []KafkaConfig     `json:"Kafka" yaml:"Kafka"`
}

type ZookeeperConfig struct {
	ZookeeperID string `json:"ZookeeperID" yaml:"ZookeeperID"`
	NodeID      string `json:"NodeID" yaml:"NodeID"`
	Port        int    `json:"Port" yaml:"Port"`
	Path        string `json:"Path" yaml:"Path"`
}

type KafkaConfig struct {
	KafkaID string `json:"KafkaID" yaml:"KafkaID"`
	NodeID  string `json:"NodeID" yaml:"NodeID"`
	Port    int    `json:"Port" yaml:"Port"`
	Path    string `json:"Path" yaml:"Path"`
}

type NodeConfig struct {
	NodeID          string `json:"NodeID" yaml:"NodeID"`
	Host            string `json:"Host" yaml:"Host"`
	Port            int    `json:"Port" yaml:"Port"`
	User            string `json:"User" yaml:"User"`
	Password        string `json:"Password" yaml:"Password"`
	CertificatePath string `json:"CertificatePath" yaml:"CertificatePath"`
}

var ConfTemp = `
{
  "Nodes": [
    {
      "NodeID": "NodeID",
      "Host": "Host",
      "Port": 22,
      "User": "User",
      "Password": "xx",
      "CertificatePath": "密码登录请删除该行"
    }
  ],
  "Zookeeper": [
    {
      "ZookeeperID": "ZookeeperID 不能重复",
      "NodeID": "上面的 NodeID",
      "Port": 2222,
      "Path": "安装的位置"
    }
  ],
  "KafkaConfig": [
    {
      "KafkaID": "KafkaID 不能重复",
      "NodeID": "上面的 NodeID",
      "Port": 6666,
      "Path": "安装的位置"
    }
  ]
}
`

func InitConf(path string) (KidConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return KidConfig{}, err
	}

	var kid KidConfig
	err = json.Unmarshal(file, &kid)
	if err != nil {
		return KidConfig{}, err
	}

	return kid, nil
}
