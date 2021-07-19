package backtask

func AllBackTask() {
	go TaskConsumeMessage()
	go CleanClient()
	go MonitoringMain()
}
