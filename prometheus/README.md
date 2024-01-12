# Metis Sequncer Prometheus/AlertManager Example

## Setup docker-compose

Update the command arguments for your node

```yaml
services:
  metis-sequencer-exporter:
    image: ghcr.io/metisprotocol/metis-sequencer-exporter:main
    pull_policy: always
    ports:
      - 21012:21012
    command:
      - -url.state.seq
      - http://PLACEHOLDER:9545/health
      - -url.state.node
      - http://PLACEHOLDER:1317/metis/latest-span
```

`-url.state.seq` is for your bridge node

`-url.state.node` is for your rest rpc node

## Setup your prometheus configuration

```yaml
alerting:
  alertmanagers:
    - static_configs:
        - targets: ["alertmanager:9093"]

rule_files:
  - "rules/*.yml"

scrape_configs:
  - job_name: "metis-sequencer-exporter"
    scrape_interval: 5s
    static_configs:
      - targets: ["metis-sequencer-exporter:21012"]
        labels:
          network_name: "metis-sepolia"
          sequencer_name: "seq-0"
```

You can add the sequencer_name label if you have multi-nodes.

## Setup your AlertManager configuration

```yaml
receivers:
  - name: "telegram"
    telegram_configs:
      - send_resolved: true
        api_url: https://api.telegram.org
        # It's your bot token
        bot_token: ""
        # Add `myidbot` to your telegram group，and send `/getgroupid @myidbot` to get the chat id.
        chat_id: null
        parse_mode: ""
```
