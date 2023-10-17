package core

import (
	monitorLogger "github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Uuid() string {
	r, _ := uuid.NewRandom()
	return r.String()
}

var logger = monitorLogger.With(
	zap.String("location", "core"),
)
