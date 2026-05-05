# Proxy Architecture & Load Balancing

This project demonstrates the implementation of Forward Proxy, Reverse Proxy, and Load Balancing using Go and Nginx.

---

## Overview

The system consists of three main components:

### 1. Forward Proxy (Client-Side)

A custom proxy built in Go that forwards client requests to external servers.
Used to understand how client-side proxying works.

### 2. Reverse Proxy & Load Balancer (Server-Side)

Nginx acts as a reverse proxy that receives client requests and distributes them across multiple backend servers.

### 3. Backend Application

A Go-based backend server designed using layered architecture and interfaces for clean and maintainable code.

---

## Key Features

* Reverse proxy using Nginx
* Load balancing across multiple backend instances
* Custom forward proxy implementation in Go
* Layered backend architecture (Handler + Service)
* Interface-based design for loose coupling
* Header forwarding using `X-Forwarded-For`

---

## Tech Stack

* Go (Golang)
* Nginx
* Curl / PowerShell

---

## Project Structure

```
/forward_proxy              # Go forward proxy
/reverse_proxy/backend     # Backend server (layered architecture)
/reverse_proxy/proxy       # Custom Go reverse proxy
/nginx                     # Nginx configuration
```

---

## How to Run

### 1. Start Backend Servers

```powershell
go run ./reverse_proxy/backend/main.go -port 9001
go run ./reverse_proxy/backend/main.go -port 9002
```

---

### 2. Start Nginx Reverse Proxy

```powershell
.\nginx-1.30.0\nginx.exe -p . -c .\nginx\nginx.conf
```

---

### 3. Test Load Balancing

```powershell
curl.exe http://localhost:9005
```

You should see responses alternating between port 9001 and 9002.

---

### 4. Run Forward Proxy (Optional)

```powershell
go run ./forward_proxy/main.go
```

Test:

```powershell
curl.exe -x http://localhost:9000 http://example.com
```

---

## Core Concepts

### Reverse Proxy

Nginx sits in front of backend servers and forwards requests, hiding internal infrastructure and improving scalability.

### Load Balancing

Incoming requests are distributed across backend instances (9001 and 9002) to improve performance and reliability.

### Interfaces in Go

The backend uses a `Greeter` interface to separate business logic from HTTP handling, enabling flexibility and easier testing.

### Dependency Injection

The service layer is injected into the handler, making the system modular and maintainable.

---

## Summary

This project demonstrates both configuration-based proxying using Nginx and code-based proxying using Go, along with scalable backend design using interfaces and layered architecture.
