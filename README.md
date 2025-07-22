![GitHub Release](https://img.shields.io/github/v/release/jellayy/gonetmon)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/jellayy/gonetmon/release.yml)
![Docker Pulls](https://img.shields.io/docker/pulls/jellayy/gonetmon)

# GoNetMon

Golang Network Monitor

Simple golang-based script + docker container for monitoring the health of your network. Pings a list of hosts every minute and pushes the response time in milliseconds to an InfluxDB2 bucket (-1 for timeouts).

# Usage

## config.yaml

GoNetMon is powered by a simple YAML configuration file that declares which hosts you want to ping and how many times you want to ping them (results averaged). If you don't pass your own config file to the container, the following default config will be loaded:

```yaml
ping_times: 5
hosts:
  - 192.168.1.1
  - 8.8.8.8
  - 1.1.1.1
  - 192.0.43.10
```

## Docker

### Docker Compose

An example docker-compose file can be found at [docker_compose.yaml](https://github.com/Jellayy/gonetmon/blob/main/docker_compose.yaml)

```yaml
services:
  gonetmon:
    image: Jellayy/gonetmon
    container_name: gonetmon
    restart: unless-stopped
    environment:
      INFLUX_HOST: https://xxx.com
      INFLUX_TOKEN: xxx
      INFLUX_ORG_NAME: yourorg
      INFLUX_BUCKET: gonetmon
    # volumes:
      # - config.yaml:config.yaml # custom ping config pass
```
