package MoGuDing

import (
	"TaskProject/serverJiang"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	phone     = "13728600493"
	password  = "Lsw981029"
	loginUrl  = "https://api.moguding.net:9000/session/user/v1/login"
	planIdUrl = "https://api.moguding.net:9000/practice/plan/v1/getPlanByStu"
	signInUrl = "https://api.moguding.net:9000/attendence/clock/v1/save"
)

type LoginRequestBody struct {
	Phone string `json:"phone"`

	Password string `json:"password"`

	LoginType string `json:"loginType"`
}

func login() (token string, err error) {
	loginStruct := LoginRequestBody{
		Phone:     phone,
		Password:  password,
		LoginType: "ios",
	}
	loginRequestBody, _ := json.Marshal(loginStruct)
	reader := bytes.NewReader(loginRequestBody)
	req, _ := http.NewRequest("POST", loginUrl, reader)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   120 * time.Second,
				KeepAlive: 120 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("失败：登录失败1 err : " + err.Error())
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var mapResult map[string]interface{}

	json.Unmarshal(body, &mapResult)

	if mapResult["code"].(float64) != 200 {
		return "", errors.New("\"失败：登录失败2 " + mapResult["code"].(string) + mapResult["msg"].(string))
	}
	mapResultData := mapResult["data"].(map[string]interface{})
	return mapResultData["token"].(string), nil
}

func getPlanId(token string) (planId string, err error) {
	requestBody := map[string]string{
		"paramsType": "student",
	}
	requestBodyByte, _ := json.Marshal(requestBody)
	reader := bytes.NewReader(requestBodyByte)
	getPlanIdRequest, _ := http.NewRequest("POST", planIdUrl, reader)
	getPlanIdRequest.Header.Set("Authorization", token)
	client := http.Client{}
	response, err := client.Do(getPlanIdRequest)
	if err != nil {
		return "", errors.New("失败：获取planId失败 " + err.Error())
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var respMap map[string]interface{}
	json.Unmarshal(body, &respMap)
	if respMap["code"].(float64) != 200 {
		return "", errors.New("失败：获取planId失败 " + string(respMap["code"].([]byte)) +string(respMap["msg"].([]byte) ))
	}

	return ((respMap["data"].([]interface{}))[0]).(map[string]interface{})["planId"].(string), nil
}

func signIn(token, planId string) (signInType string,err error) {
	signInType = "START"
	location, err := time.LoadLocation("Asia/Shanghai")
	hour :=time.Now().In(location).Hour()
	if hour >= 17 {
		signInType = "END"
	}

	requestBody := map[string]string{
		"device":         "ios",
		"planId":         planId,
		"country":        "中国",
		"state":          "NORMAL",
		"attendanceType": "",
		"address":        "海岸城（东座）",
		"type":           signInType,
		"longitude":      "113.937586",
		"city":           "深圳市",
		"province":       "广东省",
		"latitude":       "22.517891",
	}
	requestBodyByte, _ := json.Marshal(requestBody)
	reader := bytes.NewReader(requestBodyByte)
	signInRequest, _ := http.NewRequest("POST", signInUrl, reader)
	signInRequest.Header.Set("Authorization", token)
	signInRequest.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := http.Client{}
	response, err := client.Do(signInRequest)
	if err != nil {
		return "",errors.New("失败：签到提交失败 " + err.Error())
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	var respMap map[string]interface{}

	json.Unmarshal(body, &respMap)

	if respMap["code"].(float64) != 200 {
		return "",errors.New("失败：签到检验失败 " + err.Error())
	}

	return signInType,nil
}

func MoGuDingRun() {
	msg := serverJiang.NewMsg()
	msg.SetTitle("蘑菇丁签到")
	token, err := login()
	if err != nil {
		msg.AppendDesp(err.Error())
		_ = msg.Send()
		return
	}

	planId, err := getPlanId(token)
	if err != nil {
		msg.AppendDesp(err.Error())
		_ = msg.Send()
		return
	}

	signInType,err := signIn(token, planId)
	if err != nil {
		msg.AppendDesp(signInType+err.Error())
		_ = msg.Send()
		return
	}

	msg.AppendDesp(signInType+"成功：签到成功")
	_ = msg.Send()
}
