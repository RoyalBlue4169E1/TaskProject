package main

import (
	"TaskProject/MoGuDing"
	"TaskProject/dgut-yqfk"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	MoGuDing.MoGuDingRun()


	c := cron.New()
	c.AddFunc("CRON_TZ=Asia/Shanghai 10 6 * * *", func() { dgut_yqfk.YqfkRun() })
	c.AddFunc("CRON_TZ=Asia/Shanghai 40 8 * * *", func() { MoGuDing.MoGuDingRun() })
	c.AddFunc("CRON_TZ=Asia/Shanghai 15 18 * * *", func() { MoGuDing.MoGuDingRun() })

	c.Start()

	location, _ := time.LoadLocation("Asia/Shanghai")
	fmt.Println("start work : "+time.Now().In(location).Format("2006-01-02 15:04:05"))


	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-ch:
		cancel()
		fmt.Printf("\nreceived signal %s, exit.\n", s)
	}
}
