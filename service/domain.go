package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)


type Artist struct{
	CertNo	string	`json:"CertNo"`	// 证书编号 key
	Photo	string	`json:"Photo"`	// 照片
	Params string `json:"Params"`   //参数
	Type string `json:"Type"`//类型
    Date string `json:"Date"`  //交易时间
	Historys	[]HistoryItem	// 当前art的历史记录
}

type HistoryItem struct {
	TxId	string
	Artist	Artist
}
type ServiceSetup struct {
	ChaincodeID	string
	Client	*channel.Client
}
func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}
