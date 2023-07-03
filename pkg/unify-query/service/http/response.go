// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package http

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/log"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/metric"
)

type response struct {
	c *gin.Context

	actionType string
	action     string
}

func (r *response) failed(ctx context.Context, err error) {
	log.Errorf(ctx, err.Error())
	metric.RequestCountInc(ctx, r.action, r.actionType, metric.StatusFailed)
	r.c.JSON(http.StatusBadRequest, ErrResponse{
		Err: err.Error(),
	})
}

func (r *response) success(ctx context.Context, data interface{}) {
	log.Infof(ctx, "query data size is %s", fmt.Sprint(unsafe.Sizeof(data)))
	metric.RequestCountInc(ctx, r.action, r.actionType, metric.StatusSuccess)
	r.c.JSON(http.StatusOK, data)
}
