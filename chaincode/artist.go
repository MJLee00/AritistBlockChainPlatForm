
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"bytes"
)

//save artist
// args: Artist
func PutArtist(stub shim.ChaincodeStubInterface, art Artist) ([]byte, bool) {
	b, err := json.Marshal(art)
	if err != nil {
		return nil, false
	}


	err = stub.PutState(art.CertNo, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// 添加信息
// args: artistObj
// 证书号为 key, Artist 为 value
func (t *ArtistcationChaincode) addArtist(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	var art Artist
	err := json.Unmarshal([]byte(args[0]), &art)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	// 查重: artCertNo码必须唯一
	_, exist := GetArtistInfo(stub, art.CertNo)
	if exist {
		return shim.Error("要添加的证书码已存在")
	}

	_, bl := PutArtist(stub, art)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}


// 根据证书码查询信息状态
// args: certNo
func GetArtistInfo(stub shim.ChaincodeStubInterface, CertNo string) (Artist, bool)  {
	var art Artist

	b, err := stub.GetState(CertNo)
	if err != nil {
		return art, false
	}

	if b == nil {
		return art, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &art)
	if err != nil {
		return art, false
	}

	// 返回结果
	return art, true
}

// 根据证书编号及姓名查询信息
// args: CertNo, type
func (t *ArtistcationChaincode) queryArtistByCertNoAndType(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	CertNo := args[0]
	Type := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{ \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"CertNo\":\"%s\", \"Type\":\"%s\"}}", CertNo, Type)

	// 查询数据
	result, err := getArtByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据证书编号及类型查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的证书编号及类型没有查询到相关的信息")
	}
	return shim.Success(result)
}

// 根据指定的查询字符串实现富查询
func getArtByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer  resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}


// 根据证书号更新信息
// args: artistObject
func (t *ArtistcationChaincode) updateArtist(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	var info Artist
	err := json.Unmarshal([]byte(args[0]), &info)
	if err != nil {
		return  shim.Error("反序列化art信息失败")
	}

	// 根据证书号码查询信息
	result, bl := GetArtistInfo(stub, info.CertNo)
	if !bl{
		return shim.Error("根据证书号码查询信息时发生错误")
	}

	result.Photo = info.Photo
	result.CertNo = info.CertNo
	result.Type = info.Type
	result.Params = info.Params
	result.Date = info.Date
	_, bl = PutArtist(stub, result) 
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息更新成功"))
}


// 根据证书号号码查询详情（溯源）
// args: CertNo
func (t *ArtistcationChaincode) queryArtInfoByCertNo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	// 根据证书号号码查询art状态
	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据证书号码查询信息失败")
	}

	if b == nil {
		return shim.Error("根据证书号码没有查询到相关的信息")
	}

	// 对查询到的状态进行反序列化
	var art Artist
	err = json.Unmarshal(b, &art)
	if err != nil {
		return  shim.Error("反序列化art信息失败")
	}

	// 获取历史变更数据
	iterator, err := stub.GetHistoryForKey(art.CertNo)
	if err != nil {
		return shim.Error("根据指定的证书号码查询对应的历史变更数据失败")
	}
	defer iterator.Close()

	// 迭代处理
	var historys []HistoryItem
	var hisArt Artist
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取art的历史变更数据失败")
		}

		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisArt)

		if hisData.Value == nil {
			var empty Artist
			historyItem.Artist = empty
		}else {
			historyItem.Artist = hisArt
		}

		historys = append(historys, historyItem)

	}

	art.Historys = historys

	// 返回
	result, err := json.Marshal(art)
	if err != nil {
		return shim.Error("序列化art信息时发生错误")
	}
	return shim.Success(result)
}
