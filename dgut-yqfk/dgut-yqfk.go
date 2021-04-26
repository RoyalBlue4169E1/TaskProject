package dgut_yqfk

import (
	"TaskProject/serverJiang"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

const (
	dgutLoginUrl      = "https://cas.dgut.edu.cn/home/Oauth/getToken/appid/illnessProtectionHome/state/home.html"
	yqfkLoginUrl      = "https://yqfk.dgut.edu.cn/auth/auth/login"
	yqfkGetDataUrl    = "https://yqfk.dgut.edu.cn/home/base_info/getBaseInfo"
	yqfkPostDataInUrl = "https://yqfk.dgut.edu.cn/home/base_info/addBaseInfo"
)

type dgutLoginRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Token        string `json:"__token__"`
	WechatVerify string `json:"wechat_verify"`
}

func getYQFKToken() (token string, err error) {
	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}

	getDgutLoginResp, err := client.Get(dgutLoginUrl)
	if err != nil {
		return "", errors.New("失败：登录页获取 err:" + err.Error())
	}

	defer getDgutLoginResp.Body.Close()
	getDgutLoginRespBody, _ := ioutil.ReadAll(getDgutLoginResp.Body)
	compile := regexp.MustCompile(`token = "([a-z0-9]+)";`)
	find := compile.FindSubmatch(getDgutLoginRespBody)
	dgutToken := string(find[1])

	// login information
	params := url.Values{}
	params.Set("username", "201841413416")
	params.Set("password", "Liliangkun991203")
	params.Set("__token__", dgutToken)

	loginRequest, _ := http.NewRequest("POST", dgutLoginUrl, strings.NewReader(params.Encode()))
	loginRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	loginRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	loginRequest.Header.Set("X-Requested-With", "XMLHttpRequest")

	loginRequest.AddCookie(&http.Cookie{
		Name:  "languageIndex",
		Value: "0",
	})

	LoginResp, err := client.Do(loginRequest)
	if err != nil {
		return "", errors.New("失败：token获取 err:" + err.Error())
	}
	defer LoginResp.Body.Close()

	LoginRespBody, err := ioutil.ReadAll(LoginResp.Body)

	// access token
	compile = regexp.MustCompile(`"info(.*?)"}`)
	res := compile.FindAllStringSubmatch(string(LoginRespBody), -1)
	token = res[0][1]
	token = token[5:]
	token = strings.ReplaceAll(token, "\\", "")
	resp3, err := client.Get(token)
	defer resp3.Body.Close()
	if err != nil {
		panic(err)
	}
	compile = regexp.MustCompile(`access_token=(.*?)$`)
	res = compile.FindAllStringSubmatch(resp3.Request.URL.String(), -1)
	token = res[0][1]
	if token == "" {
		return "", errors.New("失败：token获取 err:" + err.Error())
	}

	return token, nil
}

func getFormData(token string) (result []byte,err error) {
	var client http.Client
	client.Jar, _ = cookiejar.New(nil)

	request, _ := http.NewRequest(http.MethodGet, yqfkGetDataUrl, nil)
	request.Header.Set("authorization", "Bearer "+token)

	resp, _ := client.Do(request)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile(`"info":(.*)}`)
	res := re.FindAllStringSubmatch(string(contents), -1)
	info := res[0][1]
	if info == "" || info == "[]" {
		return nil, errors.New("失败：获取个人信息失败  err:"+err.Error())
	}
	// 已经打卡
	if strings.Contains(info, "成功") || strings.Contains(info, "已提交") {
		return []byte(info),nil
	}
	req, _ := http.NewRequest("POST", yqfkPostDataInUrl, strings.NewReader(info))
	req.Header.Set("authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		return nil, errors.New("失败：提交失败  err:"+err.Error())
	}
	contents, _ = ioutil.ReadAll(resp.Body)

	if strings.Contains(string(contents), "成功") || strings.Contains(string(contents), "已提交") {
		return contents,nil
	} else {
		return nil, errors.New("失败：疫情防控提交失败  respContent : "+string(contents))
	}

}

func YqfkRun() {
	Msg := serverJiang.NewMsg()
	Msg.SetTitle("疫情防控")
	token, err := getYQFKToken()
	if err != nil {
		Msg.AppendDesp(err.Error())
		_ = Msg.Send()
		return
	}

	result, err := getFormData(token)
	if err != nil {
		Msg.AppendDesp(err.Error())
		_ = Msg.Send()
		return
	}

	Msg.AppendDesp("疫情防控提交成功 content : "+string(result))
	_ = Msg.Send()
}
