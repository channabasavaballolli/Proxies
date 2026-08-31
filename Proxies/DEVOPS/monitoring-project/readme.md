# Monitoring Project – NGINX + Prometheus + Grafana

## Overview
This project demonstrates a complete monitoring stack using Docker Compose.

Components:
- NGINX – Sample web application
- NGINX Prometheus Exporter – Converts NGINX statistics into Prometheus metrics
- Prometheus – Scrapes and stores time-series metrics
- Grafana – Visualizes metrics and creates alerts

---

# Architecture

```text
Browser
   |
localhost:8080
   |
NGINX
   |
/status (stub_status)
   |
NGINX Prometheus Exporter
   |
/metrics
   |
Prometheus (Scrapes every 15s)
   |
Grafana
```

---

# Project Structure

```text
monitoring-project/
│
├── docker-compose.yml
├── nginx/
│   ├── Dockerfile
│   ├── nginx.conf
│   └── index.html
│
└── prometheus/
    └── prometheus.yml
```

---

# Prerequisites

- Docker Desktop
- WSL Ubuntu
- Docker Compose

Verify:

```bash
docker --version
docker compose version
```

---

# Docker Commands

Start stack

```bash
docker compose up -d
```

Stop stack

```bash
docker compose down
```

Stop and remove volumes

```bash
docker compose down -v
```

Running containers

```bash
docker ps
```

All containers

```bash
docker ps -a
```

Logs

```bash
docker logs <container-name>
```

Enter container

```bash
docker exec -it prometheus sh
```

Restart

```bash
docker compose restart
```

Volumes

```bash
docker volume ls
```

---

# Dockerfile

Purpose:
- Build a custom NGINX image
- Copy custom webpage
- Copy nginx.conf

Example:

```dockerfile
FROM nginx:latest

COPY index.html /usr/share/nginx/html/index.html
COPY nginx.conf /etc/nginx/nginx.conf
```

---

# docker-compose.yml

Responsibilities

- Builds custom NGINX image
- Pulls Prometheus image
- Pulls Grafana image
- Pulls NGINX Exporter image
- Creates Docker network
- Creates persistent volumes

---

# Prometheus Configuration

prometheus.yml

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets: ["prometheus:9090"]

  - job_name: "nginx"
    static_configs:
      - targets: ["nginx-exporter:9113"]
```

---

# NGINX Configuration

Enable stub_status

```nginx
location /status {
    stub_status;
    allow all;
}
```

---

# URLs

NGINX

http://localhost:8080

NGINX Status

http://localhost:8080/status

Exporter

http://localhost:9113/metrics

Prometheus

http://localhost:9090

Grafana

http://localhost:3000

---

# Prometheus Queries

Health

```promql
nginx_up
```

Active Connections

```promql
nginx_connections_active
```

Reading

```promql
nginx_connections_reading
```

Writing

```promql
nginx_connections_writing
```

Waiting

```promql
nginx_connections_waiting
```

Requests/sec

```promql
irate(nginx_http_requests_total[5m])
```

Accepted Connections/sec

```promql
irate(nginx_connections_accepted[5m])
```

Handled Connections/sec

```promql
irate(nginx_connections_handled[5m])
```

---

# Graph Meaning

## nginx_up

1 = Healthy

0 = Down

---

## Active Connections

Current connected clients.

---

## Reading

NGINX is reading client requests.

---

## Writing

NGINX is sending responses.

---

## Waiting

Idle Keep-Alive connections.

---

## Requests/sec

Number of HTTP requests processed every second.

---

# Generate Traffic

```bash
while true
do
    curl http://localhost:8080 > /dev/null
    sleep 0.2
done
```

Stop

```text
Ctrl + C
```

---

# Alerting

Example Query

```promql
nginx_up
```

Condition

```
Last() IS BELOW 1
```

Pending

```
30s
```

Testing

```bash
docker stop nginx
```

Recover

```bash
docker start nginx
```

Alert Flow

```
Normal
↓
Pending
↓
Firing
↓
Normal
```

---

# Persistent Storage

```yaml
volumes:
  grafana-storage:
  prometheus-data:
```

Never use

```bash
docker compose down -v
```

unless you want to remove dashboards and metrics.

---

# Troubleshooting

Prometheus target status

Status → Targets

NGINX health

```promql
nginx_up
```

Container status

```bash
docker ps
```

View logs

```bash
docker logs prometheus
docker logs grafana
docker logs nginx
```

---

# Demo Steps

1. Run docker compose up -d
2. Show docker ps
3. Open NGINX
4. Show /status
5. Show /metrics
6. Show Prometheus Targets
7. Run PromQL queries
8. Open Grafana dashboard
9. Generate traffic
10. Observe graphs
11. Stop NGINX
12. Show alert firing
13. Start NGINX
14. Alert returns to Normal

---

# Interview Questions

## Why Prometheus?

Collects and stores time-series metrics.

## Why Grafana?

Visualizes metrics and supports alerting.

## Why Exporter?

Converts application statistics into Prometheus metrics.

## Why Docker Compose?

Runs multiple containers using one configuration.

## What is PromQL?

Query language used by Prometheus.

## What is Scraping?

Prometheus periodically pulls metrics from targets.

## Difference between Metrics and Logs?

Metrics are numerical time-series values.
Logs are textual event records.

---

# Future Improvements

- Node Exporter
- Alertmanager
- Email Alerts
- Microsoft Teams/Slack Notifications
- Application Metrics (200,404,500)
- Response Time Monitoring
- CPU, RAM and Disk Dashboards
