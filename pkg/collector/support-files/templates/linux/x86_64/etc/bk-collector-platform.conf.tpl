type: "platform"
processor:
{% if apdex_config is defined %}
  # ApdexCalculator: 健康度状态计算器
  - name: "{{ apdex_config.name }}"
    config:
      calculator:
        type: "{{ apdex_config.type }}"
      rules:
        {%- for rule_config in apdex_config.rules %}
        - kind: "{{ rule_config.kind }}"
          predicate_key: "{{ rule_config.predicate_key }}"
          metric_name: "{{ rule_config.metric_name }}"
          destination: "{{ rule_config.destination }}"
          apdex_t: {{ rule_config.apdex_t }} # ms
        {%- endfor %}
{%- endif %}

{% if sampler_config is defined %}
  # Sampler: 采样处理器
  - name: "{{ sampler_config.name }}"
    config:
      type: "{{ sampler_config.type }}"
      sampling_percentage: {{ sampler_config.sampling_percentage }}
{%- endif %}

{% if qps_config is defined %}
  # Qps: Qps限流
  - name: "{{ qps_config.name }}"
    config:
      type: "{{ qps_config.type }}"
      qps: {{ qps_config.qps }}
      burst: {{ qps_config.qps }}
{%- endif %}

{% if token_checker_config is defined %}
  # TokenChecker: 权限校验处理器
  - name: "{{ token_checker_config.name }}"
    config:
      type: "{{ token_checker_config.type }}"
      resource_key: "{{ token_checker_config.resource_key }}"
      salt: "{{ token_checker_config.salt }}"
      decoded_key: "{{ token_checker_config.decoded_key }}"
      decoded_iv: "{{ token_checker_config.decoded_iv }}"
{%- endif %}

{% if resource_filter_config is defined %}
  # ResourceFilter: 资源过滤处理器
  - name: "{{ resource_filter_config.name }}"
    config:
      assemble:
        {%- for as_config in  resource_filter_config.assemble %}
        - destination: "{{ as_config.destination }}"
          separator: "{{ as_config.separator }}"
          keys:
            {%- for key in as_config.get("keys", []) %}
            - "{{ key }}"
            {%- endfor %}
        {%- endfor %}
      drop:
        keys:
          {%- for drop_key in resource_filter_config.get("drop", {}).get("keys", []) %}
          - "{{ drop_key }}"
          {%- endfor %}
{%- endif %}

{% if metric_configs is defined %}
  # bk_apm_count
  {% if metric_configs.metric_bk_apm_count_config is defined %}
  - name: "{{ metric_configs.metric_bk_apm_count_config.name }}"
    config:
      operations:
        {%- for operation in metric_configs.metric_bk_apm_count_config.operations %}
        - type: "{{ operation.type }}"
          metric_name: "{{ operation.metric_name }}"
          rules:
            {%- for rule_config in operation.rules %}
            - kind: "{{ rule_config.kind }}"
              predicate_key: "{{ rule_config.predicate_key }}"
              dimensions:
                {%- for dimension_key in rule_config.dimensions %}
                - "{{ dimension_key }}"
                {%- endfor %}
            {%- endfor %}
        {%- endfor %}
  {%- endif %}

  # bk_apm_total
  {% if metric_configs.metric_bk_apm_total_config is defined %}
  - name: "{{ metric_configs.metric_bk_apm_total_config.name }}"
    config:
      operations:
        {%- for operation in metric_configs.metric_bk_apm_total_config.operations %}
        - type: "{{ operation.type }}"
          metric_name: "{{ operation.metric_name }}"
          rules:
            {%- for rule_config in operation.rules %}
            - kind: "{{ rule_config.kind }}"
              predicate_key: "{{ rule_config.predicate_key }}"
              dimensions:
                {%- for dimension_key in rule_config.dimensions %}
                - "{{ dimension_key }}"
                {%- endfor %}
            {%- endfor %}
        {%- endfor %}
  {%- endif %}

  # bk_apm_duration
  {% if metric_configs.metric_bk_apm_duration_config is defined %}
  - name: "{{ metric_configs.metric_bk_apm_duration_config.name }}"
    config:
      operations:
        {%- for operation in metric_configs.metric_bk_apm_duration_config.operations %}
        - type: "{{ operation.type }}"
          metric_name: "{{ operation.metric_name }}"
          rules:
            {%- for rule_config in operation.rules %}
            - kind: "{{ rule_config.kind }}"
              predicate_key: "{{ rule_config.predicate_key }}"
              dimensions:
                {%- for dimension_key in rule_config.dimensions %}
                - "{{ dimension_key }}"
                {%- endfor %}
            {%- endfor %}
        {%- endfor %}
  {%- endif %}

  # bk_apm_duration_max
  {% if metric_configs.metric_bk_apm_duration_max_config is defined %}
  - name: "{{ metric_configs.metric_bk_apm_duration_max_config.name }}"
    config:
      operations:
        {%- for operation in metric_configs.metric_bk_apm_duration_max_config.operations %}
        - type: "{{ operation.type }}"
          metric_name: "{{ operation.metric_name }}"
          rules:
            {%- for rule_config in operation.rules %}
            - kind: "{{ rule_config.kind }}"
              predicate_key: "{{ rule_config.predicate_key }}"
              dimensions:
                {%- for dimension_key in rule_config.dimensions %}
                - "{{ dimension_key }}"
                {%- endfor %}
            {%- endfor %}
          {%- endfor %}
  {%- endif %}

  # bk_apm_duration_min
  {% if metric_configs.metric_bk_apm_duration_min_config is defined %}
  - name: "{{ metric_configs.metric_bk_apm_duration_min_config.name }}"
    config:
      operations:
        {%- for operation in metric_configs.metric_bk_apm_duration_min_config.operations %}
        - type: "{{ operation.type }}"
          metric_name: "{{ operation.metric_name }}"
          rules:
            {%- for rule_config in operation.rules %}
            - kind: "{{ rule_config.kind }}"
              predicate_key: "{{ rule_config.predicate_key }}"
              dimensions:
                {%- for dimension_key in rule_config.dimensions %}
                - "{{ dimension_key }}"
                {%- endfor %}
            {%- endfor %}
          {%- endfor %}
  {%- endif %}

  # bk_apm_duration_sum
  {% if metric_configs.metric_bk_apm_duration_sum_config is defined %}
  - name: "{{ metric_configs.metric_bk_apm_duration_sum_config.name }}"
    config:
      operations:
        {%- for operation in metric_configs.metric_bk_apm_duration_sum_config.operations %}
        - type: "{{ operation.type }}"
          metric_name: "{{ operation.metric_name }}"
          rules:
            {%- for rule_config in operation.rules %}
            - kind: "{{ rule_config.kind }}"
              predicate_key: "{{ rule_config.predicate_key }}"
              dimensions:
                {%- for dimension_key in rule_config.dimensions %}
                - "{{ dimension_key }}"
                {%- endfor %}
            {%- endfor %}
        {%- endfor %}
  {%- endif %}

   # bk_apm_duration_bucket
  {% if metric_configs.metric_bk_apm_duration_bucket_config is defined %}
  - name: "{{ metric_configs.metric_bk_apm_duration_bucket_config.name }}"
    config:
      operations:
        {%- for operation in metric_configs.metric_bk_apm_duration_bucket_config.operations %}
        - type: "{{ operation.type }}"
          metric_name: "{{ operation.metric_name }}"
          buckets: {{ operation.buckets }}
          rules:
            {%- for rule_config in operation.rules %}
            - kind: "{{ rule_config.kind }}"
              predicate_key: "{{ rule_config.predicate_key }}"
              dimensions:
                {%- for dimension_key in rule_config.dimensions %}
                - "{{ dimension_key }}"
                {%- endfor %}
            {%- endfor %}
        {%- endfor %}
  {%- endif %}

{%- endif %}
