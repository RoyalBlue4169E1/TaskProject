package serverJiang

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	ServerJiangUrl="https://sc.ftqq.com/SCU153450Tbf04ca15a918626969ac31b984a498c1600a83fad8b34.send"
)

type ServerJiangMsg struct {
	Text string `form:"text"`
	Desp string	`form:"desp"`
}

func NewMsg()(Msg *ServerJiangMsg){
	return &ServerJiangMsg{}
}

func (Msg *ServerJiangMsg)SetTitle(text string){
	Msg.Text = text
}

func (Msg *ServerJiangMsg)AppendDesp(desp string){
	Msg.Desp=Msg.Desp+desp+"->"
}

func (Msg *ServerJiangMsg)Send() error {
	fmt.Println(Msg.Desp)

	if len(Msg.Text)<=0|| len(Msg.Desp)<=0{
		return errors.New("标题或内容为空")
	}

	values := url.Values{}
	values.Set("text",Msg.Text)
	values.Set("desp",Msg.Desp[:len(Msg.Desp)-2])

	response, err := http.PostForm(ServerJiangUrl, values)

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer response.Body.Close()

	respBody, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(respBody))

	return nil

}

