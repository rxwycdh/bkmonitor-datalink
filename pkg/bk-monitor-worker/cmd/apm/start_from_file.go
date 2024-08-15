package apm

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/config"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/apm/pre_calculate/tools"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/log"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
)

var filePath string

func StartFromFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start_from_file",
		Short: "start apm task from file",
		Run:   listenFile,
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "connection file")
	return cmd
}

func listenFile(cmd *cobra.Command, args []string) {
	config.InitConfig()
	log.InitLogger()

	ctx, cancel := context.WithCancel(context.Background())
	if err := tools.StartListenerFromFile(ctx, filePath); err != nil {
		panic(err)
	}
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
