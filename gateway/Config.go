package gateway

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	WsPort               int    `json:"wsPort"`
	WsReadTimeout        int    `json:"wsReadTimeout"`
	WsWriteTimeout       int    `json:"wsWriteTimeout"`
	WsInChannelSize      int    `json:"wsInChannelSize"`
	WsOutChannelSize     int    `json:"wsOutChannelSize"`
	WsHeartbeatInterval  int    `json:"wsHeartbeatInterval"`
	MaxMergerDelay       int    `json:"maxMergerDelay"`
	MaxMergerBatchSize   int    `json:"maxMergerBatchSize"`
	MergerWorkerCount    int    `json:"mergerWorkerCount"`
	MergerChannelSize    int    `json:"mergerChannelSize"`
	ServicePort          int    `json:"servicePort"`
	ServiceReadTimeout   int    `json:"serviceReadTimeout"`
	ServiceWriteTimeout  int    `json:"serviceWriteTimeout"`
	ServerPem            string `json:"serverPem"`
	ServerKey            string `json:"serverKey"`
	BucketCount          int    `json:"bucketCount"`
	BucketWorkerCount    int    `json:"bucketWorkerCount"`
	MaxJoinRoom          int    `json:"maxJoinRoom"`
	DispatchChannelSize  int    `json:"dispatchChannelSize"`
	DispatchWorkerCount  int    `json:"dispatchWorkerCount"`
	BucketJobChannelSize int    `json:"bucketJobChannelSize"`
	BucketJobWorkerCount int    `json:"bucketJobWorkerCount"`
}

var (
	G_config *Config
)

func InitConfig(fileName string) (err error) {
	var (
		content []byte
		config  Config
	)
	if content, err = ioutil.ReadFile(fileName); err != nil {
		return
	}
	if err = json.Unmarshal(content, &config); err != nil {
		return
	}
	G_config = &config
	return
}
