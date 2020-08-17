package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
)

type ArtistcationChaincode struct {

}

func (t *ArtistcationChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{

	return shim.Success(nil)
}

func (t *ArtistcationChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "addArtist"{
		return t.addArtist(stub, args)		// 添加信息
	}else if fun == "queryArtistByCertNoAndType" {
		return t.queryArtistByCertNoAndType(stub,args)		// 根据证书编号及类型查询信息
	}else if fun == "updateArtist" {
		return t.updateArtist(stub, args)		// 根据证书编号更新信息
	}else if fun == "queryArtInfoByCertNo" {
		return t.queryArtInfoByCertNo(stub, args)		// 根据证书编号查询历史信息
	}

	return shim.Error("指定的函数名称错误")

}

func main(){
	err := shim.Start(new(ArtistcationChaincode))
	if err != nil{
		fmt.Printf("启动EducationChaincode时发生错误: %s", err)
	}
}
