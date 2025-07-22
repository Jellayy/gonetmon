package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"

	"gopkg.in/yaml.v3"
	"github.com/influxdata/influxdb-client-go/v2"
)

type Config struct {
	PingTimes int      `yaml:"ping_times"`
	Hosts     []string `yaml:"hosts"`
}

func (c *Config) validate() error {
	// validate ping_times
	if c.PingTimes <= 0 {
		return fmt.Errorf("ping_times must be greater than 0, got: %d", c.PingTimes)
	}
	if c.PingTimes > 100 {
		return fmt.Errorf("hey pal, icmp is not a ddos tool, tone down those ping_times, got: %d", c.PingTimes)
	}

	// validate hosts
	if len(c.Hosts) == 0 {
		return fmt.Errorf("hosts list cannot be empty")
	}
	for i, host := range c.Hosts {
		if net.ParseIP(host) == nil {
			return fmt.Errorf("invalid IP address at hosts[%d]: %s", i, host)
		}
	}

	return nil
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	return &config, nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime)

	// load & validate config.yaml ping config (hosts, ping attempts)
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// load environment config for influx host, token, etc
	influx_url := os.Getenv("INFLUX_HOST")
	if influx_url == "" {
		log.Fatalf("INFLUX_HOST not found in environment")
	}
	influx_token := os.Getenv("INFLUX_TOKEN")
	if influx_token == "" {
		log.Fatalf("INFLUX_TOKEN not found in environment")
	}
	influx_org := os.Getenv("INFLUX_ORG_NAME")
	if influx_org == "" {
		log.Fatalf("INFLUX_ORG_NAME not found in environment")
	}
	influx_bucket := os.Getenv("INFLUX_BUCKET")
	if influx_bucket == "" {
		log.Fatalf("INFLUX_BUCKET not found in environment")
	}

	// influx init
	influx_client := influxdb2.NewClient(influx_url, influx_token)
	writeAPI := influx_client.WriteAPIBlocking(influx_org, influx_bucket)

	// setup event loop schedule
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// main event loop
	for {
		select {
		case <-ticker.C:
			// ping hosts & write result to influx
			for _, host := range config.Hosts {
				avgPingTime, err := ping(host, config.PingTimes)
				if err != nil {
					log.Printf("Error pinging host: %v\n", err)
				}
				log.Printf("Sending stat: %f for host: %s to influx", avgPingTime, host)
				p := influxdb2.NewPoint("stat",
					map[string]string{"unit": "ping_time", "host": host},
					map[string]interface{}{"avg": avgPingTime},
					time.Now())
				writeAPI.WritePoint(context.Background(), p)
			}
		case <-sigChan:
			log.Println("Cancel recieved, shutting down...")
			return
		case <-ctx.Done():
			return
		}
	}
}
