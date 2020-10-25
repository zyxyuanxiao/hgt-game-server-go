package app

import (
	"fmt"
	"github.com/robfig/cron"
	"server/exception"
)

type CronManager struct{
	deleteChan chan string
	feedbackChan chan bool
	clients map[string]*cron.Cron
}

var cronManagerVar = CronManager{
	deleteChan: make(chan string),
	feedbackChan: make(chan bool),
	clients:    make(map[string]*cron.Cron),
}
var CronCommands = make(map[string]func())

// 加载 cron
func cronLoad() {
	//CronCommands["TestAdd"] = console.TestAdd
}

// 初始化
func CronStart() {
	cronLoad()
	for {
		select {
		case cronName := <-cronManagerVar.deleteChan:
			if client, ok := cronManagerVar.clients[cronName]; ok {
				fmt.Println("cron:" + cronName + " stop")
				delete(cronManagerVar.clients, cronName)
				client.Stop()
				cronManagerVar.feedbackChan <- true
			} else {
				feedbackMsg := "cron:" + cronName + " not exist"
				cronManagerVar.feedbackChan <- false
				fmt.Println(feedbackMsg)
			}
		}
	}
}

// 注册 cron 任务
func CronRegister(cronName string, spec string, handle func()) bool {
	if _, ok := cronManagerVar.clients[cronName]; ok {
		exception.Logic("cron already exist")
	}
	c := cron.New()
	c.AddFunc(spec, handle)

	cronManagerVar.clients[cronName] = c
	go c.Start()

	return true
}

// 停止 cron 任务
func CronStop(cronName string) bool {
	cronManagerVar.deleteChan <- cronName

	select {
	case feedbackMsg := <- cronManagerVar.feedbackChan:
		return feedbackMsg
	}

	return false
}