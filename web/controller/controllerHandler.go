
package controller

import (
	"net/http"
	"encoding/json"
	"github.com/kongyixueyuan.com/education/service"
	"fmt"
)

// 添加信息
func (app *Application) AddArt(w http.ResponseWriter, r *http.Request)  {

	art := service.Artist{
		CertNo:r.FormValue("certNo"),
		Photo:r.FormValue("photo"),
		Type: r.FormValue("type"),
		Params:r.FormValue("params"),
		Date: r.FormValue("date"),
	}
    
	txId ,err :=app.Setup.SaveArt(art)

	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json") //返回数据格式是json
    if err != nil{
		w.WriteHeader(500)
		w.Write([]byte("the certNo has existed!"))
	}else{
		r,err := json.Marshal(txId)
		if err !=nil{
			w.WriteHeader(500)
			w.Write([]byte("error"))
		}else{
			w.WriteHeader(200)
			w.Write(r)
		}
	}
//	r.Form.Set("certNo", art.CertNo)
	//r.Form.Set("type", art.Type)
//	app.FindInfoByCertAndType(w, r)

}

// 根据证书编号与类型查询信息
func (app *Application) FindInfoByCertAndType(w http.ResponseWriter, r *http.Request)  {
	certNo := r.FormValue("certNo")
	Type := r.FormValue("type")
	result, err := app.Setup.FindArtByCertNoAndType(certNo, Type)
	var art = service.Artist{}
	json.Unmarshal(result, &art)

	fmt.Println("根据证书编号与类型查询信息成功：")
	fmt.Println(art)

	data := &struct {
		Art service.Artist
		Msg string
		Flag bool
		History bool
	}{
		Art:art,
		Msg:"",
		Flag:false,
		History:false,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}
	var buf []byte

	buf, err = json.Marshal(data)
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json") //返回数据格式是json
	if buf != nil {
		w.WriteHeader(200)
		w.Write(buf)
	}else {
		w.WriteHeader(500)
		w.Write([]byte("error"))
	}
}

// 根据证书号码查询信息
func (app *Application) FindByCertNo(w http.ResponseWriter, r *http.Request)  {
	certNo := r.FormValue("certNo")
	result, err := app.Setup.FindArtInfoByCerNo(certNo)
	var art = service.Artist{}
	json.Unmarshal(result, &art)

	data := &struct {
		Art service.Artist
		Msg string
		Flag bool
		History bool
	}{
		Art:art,
		Msg:"",
		Flag:false,
		History:true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	var buf []byte
	buf, err = json.Marshal(data)
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json") //返回数据格式是json
	if buf != nil {
		w.WriteHeader(200)
		w.Write(buf)
	}else {
		w.WriteHeader(500)
		w.Write([]byte("error"))
	}
}

// 修改/添加新信息
func (app *Application) Modify(w http.ResponseWriter, r *http.Request) {
	art := service.Artist{
		CertNo:r.FormValue("certNo"),
		Photo:r.FormValue("photo"),
		Type: r.FormValue("type"),
		Params: r.FormValue("params"),
		Date:r.FormValue("date"),
	}

	t,err := app.Setup.ModifyArt(art)
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json") //返回数据格式是json
    if err != nil{
		w.WriteHeader(500)
		w.Write([]byte("error"))
	}else{
		r,err := json.Marshal(t)
		if err !=nil{
			w.WriteHeader(500)
			w.Write([]byte("error"))
		}else{
			w.WriteHeader(200)
			w.Write(r)
		}
	}
	//r.Form.Set("certNo", art.CertNo)
//	r.Form.Set("type", art.Type)
//	app.FindInfoByCertAndType(w, r)
}
