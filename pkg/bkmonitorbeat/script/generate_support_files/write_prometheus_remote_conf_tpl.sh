#!/bin/bash
# Tencent is pleased to support the open source community by making
# 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
# Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
# Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
# You may obtain a copy of the License at http://opensource.org/licenses/MIT
# Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
# an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
# specific language governing permissions and limitations under the License.


write_prometheus_remote_conf_tpl() {
  path="$1"
  system="$2"
  arch="$3"
  cat <<EOF > "$path"
EOF
  cat <<EOF >> "$path"
# 子配置信息
type: metricbeat
name: {{ config_name }}
version: {{ config_version }}

# 最大超时时间
max_timeout: {{ max_timeout | default(100, true) }}s
# 最小检测间隔
min_period: {{ min_period | default(3, true) }}s

dataid: {{ dataid }}

tasks: {% for task in tasks %}
  - task_id: {{ task.task_id }}
    bk_biz_id: {{ task.bk_biz_id }}
    # 周期
    period: {{ task.period }}s
    # 超时
    timeout: {{ task.timeout | default(60, true) }}s
    module:
      module: prometheus
      metricsets: ["collector"]
      enabled: true
      hosts: ["{{ task.metric_url }}"]
      metrics_path: ''
      namespace: {{ task.config_name }}
      dataid: {{ task.dataid }}
      {% if task.diff_metrics %}diff_metrics:
      {% for metric in task.diff_metrics %}- {{ metric }}
      {% endfor %}
      {% endif %}
    {% if task.labels %}labels:
    {% for label in task.labels %}{% for key, value in label.items() %}{{ "-" if loop.first else " "  }} {{ key }}: "{{ value }}"
    {% endfor %}{% endfor %}
    {% endif %}{% endfor %}
EOF
}
