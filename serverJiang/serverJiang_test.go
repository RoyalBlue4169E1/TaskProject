package serverJiang

import (
	"testing"
)

func TestServerJiang(t *testing.T) {
	msg := NewMsg()
	msg.SetTitle("测试Server酱")
	msg.AppendDesp("测试字符串1")
	msg.AppendDesp("测试字符串2")
	msg.AppendDesp("测试字符串3")

	err := msg.Send()
	if err != nil {
		panic(err)
	}

}