// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

// Code generated by "stringer -type=WatchEventType -trimprefix WatchEvent"; DO NOT EDIT.

package define

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[WatchEventUnknown-0]
	_ = x[WatchEventAdded-1]
	_ = x[WatchEventDeleted-2]
	_ = x[WatchEventModified-3]
	_ = x[WatchEventNoChange-4]
}

const _WatchEventType_name = "UnknownAddedDeletedModifiedNoChange"

var _WatchEventType_index = [...]uint8{0, 7, 12, 19, 27, 35}

func (i WatchEventType) String() string {
	if i < 0 || i >= WatchEventType(len(_WatchEventType_index)-1) {
		return "WatchEventType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _WatchEventType_name[_WatchEventType_index[i]:_WatchEventType_index[i+1]]
}
