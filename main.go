package main

import (
	"os"
	"fmt"
	"github.com/kongyixueyuan.com/education/sdkInit"
	"github.com/kongyixueyuan.com/education/service"
	"encoding/json"
	"github.com/kongyixueyuan.com/education/web/controller"
	"github.com/kongyixueyuan.com/education/web"
	"time"
)

const (
	configFile = "config.yaml"
	initialized = false
	SimpleCC = "simplecc"
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID: "kevinkongyixueyuan",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/kongyixueyuan.com/education/fixtures/artifacts/channel.tx",

		OrgAdmin:"Admin",
		OrgName:"Org1",
		OrdererOrgName: "orderer.kevin.kongyixueyuan.com",

		ChaincodeID: SimpleCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath: "github.com/kongyixueyuan.com/education/chaincode/",
		UserName:"User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()

	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)

	//===========================================//

	serviceSetup := service.ServiceSetup{
		ChaincodeID:SimpleCC,
		Client:channelClient,
		
	}


	t := time.Now() 
	ret :=t.Format("2006-01-02 15:04:05")

	art := service.Artist{
		CertNo: "111",
		Photo: "/static/photo/11.png",
		Type: "jade",
		Params: "this is test one",
		Date:ret,
	}

	art2 := service.Artist{
		CertNo: "222",
		Photo: "/static/photo/22.png",
		Type: "fan",
		Params: "this is test two ",
		Date: ret,
	}

	msg, err := serviceSetup.SaveArt(art)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}

	msg1, err := serviceSetup.SaveArt(art2)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息发布成功, 交易编号为: " + msg1)
	}
	
	// 根据证书编号与类型查询信息
	result, err := serviceSetup.FindArtByCertNoAndType("111","jade")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var art service.Artist
		json.Unmarshal(result, &art)
		fmt.Println("根据证书编号与类型查询信息成功：")
		fmt.Println(art)
	}

	// 根据certno号码查询信息
	result, err = serviceSetup.FindArtInfoByCerNo("111")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var art service.Artist
		json.Unmarshal(result, &art)
		fmt.Println("根据证书号码查询信息成功：")
		fmt.Println(art)
	}

	ts := time.Now() 
	rets  := ts.Format("2006-01-02 15:04:05")
	// 修改/添加信息
	info := service.Artist{
		CertNo: "111",
		Photo: "/static/photo/11.png",
		Type: "jade",
		Params: "this is test one after update",
		Date : rets,
	}
	msg, err = serviceSetup.ModifyArt(info)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息操作成功, 交易编号为: " + msg)
	}


	// 根据certno号码查询信息
	result, err = serviceSetup.FindArtInfoByCerNo("111")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var art service.Artist
		json.Unmarshal(result, &art)
		fmt.Println("根据证书号码查询信息成功：")
		fmt.Println(art)
	}


	app := controller.Application{
		Setup: &serviceSetup,
	}
	web.WebStart(app)

}