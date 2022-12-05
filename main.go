package main

import (
	"log"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	cron "github.com/robfig/cron/v3"
)

func main() {
	// set scheduler timezone
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakartaTime))
	// stop scheduler exact before end of function
	defer scheduler.Stop()
	// set scheduler task, please adjust as your needs
	scheduler.AddFunc("13 12 * * *", PushObject)
	// start scheduler
	go scheduler.Start()
	// trap SIGINT for shutdown trigger
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
// push to bucket
func PushObject() {
  currentTime := time.Now()
	app := "gsutil"
	opt := "cp"
	src := "<your log file source>"
	bucket := "<your google cloud storage uri>"
  dest := fmt.Sprintf("%s%s", bucket, currentTime.Format("2006-01-02"))

	cmd, err := exec.Command(app, "-m", opt, "-r", src, dest).Output()
	if err != nil {
		log.Fatalln("Error excute:", err)
		return
	}
	fmt.Println(string(cmd))
}
