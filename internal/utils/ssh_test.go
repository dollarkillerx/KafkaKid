package utils

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/yahoo/vssh"
)

func TestConnect(t *testing.T) {
	vs := vssh.New().Start()
	auth := vssh.GetConfigUserPass("mv-user", "rimepevc2021")
	vs.AddClient("192.168.88.39:22", auth, vssh.SetMaxSessions(4))
	vs.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//session, err := Connect("mv-user", "rimepevc2021", "192.168.88.39", 22, "", nil)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//defer session.Close()
	////var stdoutBuf bytes.Buffer
	////session.Stdout = &stdoutBuf
	//session.Stdout = os.Stdout
	//session.Stderr = os.Stderr
	for {

		var cmd string
		fmt.Print(">")
		fmt.Scan(&cmd)

		timeout, _ := time.ParseDuration("5s")
		respChan := vs.Run(ctx, cmd, timeout)

		resp := <-respChan
		if err := resp.Err(); err != nil {
			log.Fatal(err)
		}

		stream := resp.GetStream()
		defer stream.Close()

		for stream.ScanStdout() {
			txt := stream.TextStdout()
			fmt.Println(txt)
		}

		if err := stream.Err(); err != nil {
			log.Fatal(err)
		}
	}
}
