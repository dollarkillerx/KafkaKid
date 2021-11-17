package main

import (
	"dollarkillerx/KafkaKid/internal/conf"
	"dollarkillerx/KafkaKid/internal/utils"
	"github.com/spf13/cobra"

	"errors"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	execute()
}

var rootCmd = &cobra.Command{}

var configFileName string

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFileName, "config_name", "c", "kafka_kid.json", "config file json")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startKafka)
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init KafkaKid Config",
	Long:  "Init KafkaKid Config",
	Run: func(cmd *cobra.Command, args []string) {
		err := ioutil.WriteFile("kafka_kid.json", []byte(conf.ConfTemp), 00666)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

var startKafka = &cobra.Command{
	Use:   "start",
	Short: "Start Kafka Cluster",
	Long:  "Start Kafka Cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化配置
		cfg, err := conf.InitConf(configFileName)
		if err != nil {
			log.Fatalln(err)
		}
		// validate
		if len(cfg.Kafka) == 0 || len(cfg.Zookeeper) == 0 || len(cfg.Nodes) == 0 {
			log.Fatalln("请认真填写配置文件")
		}

		for _, v := range cfg.Nodes {
			if v.NodeID == "" || v.Host == "" || v.Port == 0 {
				log.Fatalln("Nodes 请认真填写配置文件")
			}
		}

		for _, v := range cfg.Kafka {
			if v.Path == "" || v.KafkaID == "" || v.NodeID == "" || v.Port == 0 {
				log.Fatalln("Kafka 请认真填写配置文件")
			}
		}

		for _, v := range cfg.Zookeeper {
			if v.Path == "" || v.ZookeeperID == "" || v.NodeID == "" || v.Port == 0 {
				log.Fatalln("Zookeeper 请认真填写配置文件")
			}
		}

		nMap := initNodeManager(cfg.Nodes)

		for _, v := range cfg.Zookeeper {
			log.Println("开始部署 Zookeeper: ", v.ZookeeperID, " Node: ", v.NodeID)
			err := deployZookeeper(v, nMap)
			if err != nil {
				log.Fatalln(err)
			}
		}

	},
}

func deployZookeeper(node conf.ZookeeperConfig, nMap nodeManager) error {
	sr, ex := nMap.node[node.NodeID]
	if !ex {
		return errors.New("zookeeper find node 404")
	}

	session, err := utils.Connect(sr.User, sr.Password, sr.Host, sr.Port, sr.CertificatePath, nil)
	if err != nil {
		return err
	}

	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run("mkdir -p " + node.Path)
	if err != nil {
		return err
	}
	return nil
}

type nodeManager struct {
	node map[string]conf.NodeConfig
}

func initNodeManager(nodes []conf.NodeConfig) nodeManager {
	nMap := map[string]conf.NodeConfig{}
	for _, v := range nodes {
		nMap[v.NodeID] = v
	}

	return nodeManager{node: nMap}
}
