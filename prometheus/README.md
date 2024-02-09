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
    extra_hosts:
      - "host.docker.internal:host-gateway"
    command:
      # deploy the exporter on the same instance with your sequencer
      - -url.state.seq
      - http://host.docker.internal:9545/health
      - -url.state.node
      - http://host.docker.internal:1317/metis/latest-span
      - -url.state.l1dtl
      - http://host.docker.internal:7878/eth/context/latest
```

`-url.state.seq` is for your bridge node

`-url.state.node` is for your rest rpc node

`-url.state.l1dtl` is for your dtl service

if you want to deploy the exporter on the same instance with your sequencer,

you can use `host.docker.internal` as the host of the url.

## Setup your prometheus configuration with `./config/prometheus.yml` file

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
          sequencer_name: "__YOUR_SEQUENCER_NAME__"
```

You can add the sequencer_name label if you have multi-nodes.

## Add basic auth for for your prometheus in `./config/prometheus-web.yml` file

```yml
basic_auth_users:
  # Usernames and hashed passwords that have full access to the web server via basic authentication.
  # you can use https://bcrypt-generator.com/ to generate the password
  admin: "$2a$12$/URNoeIEVn0IdOqjO3QcxOmOH6lVSJ61uGMKmfkXD/.rq2rbNjFYe"
```

## Setup your AlertManager configuration

if you want to use your custom receiver like slack and email, please refer to https://prometheus.io/docs/alerting/latest/configuration/

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
