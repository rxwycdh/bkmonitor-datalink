package pre_calculate

import (
	"context"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

const dataIdFilePath = "./connections_test.yaml"

func TestApmPreCalculateViaFile(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	op, err := Initial(ctx)
	if err != nil {
		t.Fatal(err)
	}

	go op.Run()
	go op.WatchConnections(dataIdFilePath)

	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		switch <-s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cancel()
			logger.Infof("Bye")
			os.Exit(0)
		}
	}
}
