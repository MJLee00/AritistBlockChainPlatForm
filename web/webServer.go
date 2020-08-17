package web

import (
	"net/http"
	"fmt"
	"github.com/kongyixueyuan.com/education/web/controller"

)
// 启动Web服务并指定路由信息
func WebStart(app controller.Application)  {

	fs:= http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 指定路由信息(匹配请求)
	http.HandleFunc("/addArt", app.AddArt)	// 提交信息请求

	http.HandleFunc("/query", app.FindInfoByCertAndType)	// 根据证书编号与类型查询信息

	http.HandleFunc("/query2", app.FindByCertNo)	// 根据证书号查询信息

	http.HandleFunc("/modify", app.Modify)	//  修改信息

	http.HandleFunc("/upload", app.UploadFile)

	fmt.Println("启动Web服务, 监听端口号为: 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("Web服务启动失败: %v", err)
	}

	
}
