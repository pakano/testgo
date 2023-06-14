package elk

import (
	"log"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
)

var EsClient *elastic.Client

// InitEs 初始化es
func InitEs(eshost string, esport string) {
	esConn := "http://" + eshost + ":" + esport
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(esConn))
	if err != nil {
		log.Panic(err)
	}
	EsClient = client
}

// EsHookLog 初始化log日志
func EsHookLog(eshost string, index string) *elogrus.ElasticHook {
	hook, err := elogrus.NewElasticHook(EsClient, eshost, logrus.DebugLevel, index)
	if err != nil {
		log.Panic(err)
	}
	return hook
}
