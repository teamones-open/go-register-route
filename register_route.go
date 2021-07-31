package register_route

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type RegisterParam struct {
	BelongSystem string `json:"belong_system"` // 所属系统
	Record       string `json:"record"`        // 路由记录
	Method       string `json:"method"`        // 请求方式
}

type UpdateData struct {
	Data []RegisterParam `json:"data"`
}

// 注册上报路由
func Register(routes gin.RoutesInfo, url string, BelongSystem string) (err error) {

	if len(routes) > 0 {
		var routesParam []RegisterParam

		// 整理成添加数据集
		for _, value := range routes {
			routesParam = append(routesParam, RegisterParam{
				BelongSystem: BelongSystem,
				Record:       value.Path,
				Method:       strings.ToUpper(value.Method),
			})
		}

		// 将日志数据转换为JSON
		updateData := UpdateData{
			Data: routesParam,
		}
		payload, err := json.Marshal(updateData)

		if err != nil {
			return err
		}

		body := bytes.NewBuffer(payload)

		// 创建POST请求
		req, err := http.NewRequest("POST", url, body)

		if err != nil {
			return err
		}

		// 设置header头
		req.Header.Add("content-type", "application/json")

		// 发送请求
		var httpClient = http.Client{}
		resp, err := httpClient.Do(req)

		if err != nil {
			return err
		}

		// 如果状态代码大于201，则返回错误
		if resp.StatusCode > http.StatusCreated {
			return errors.New(fmt.Sprintf("failed to post payload, the server responded with a status of %v", resp.StatusCode))
		}

	} else {
		return errors.New("Route does not exist")
	}

	return
}
