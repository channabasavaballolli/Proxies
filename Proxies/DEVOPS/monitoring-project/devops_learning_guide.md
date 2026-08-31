# 🚀 Beginner's Guide: Understanding Your Monitoring & Observability Project

Welcome to your DevOps learning journey! Since you are at the beginning stage, this guide is written in simple, plain language to explain exactly what was built, why we built it this way, and how all the components work together.

---

## 🧭 The Big Picture: What is this project?

Imagine you are running a cloud hosting service (like IBM Cloud) where users can create temporary testing environments called **Sandboxes**. Inside these sandboxes, they can boot up virtual computers called **VSIs (Virtual Server Instances)**.

To make sure this service runs smoothly, you (the DevOps Engineer) must monitor it. You need to know:
1. Is the website up?
2. Is it slow?
3. Is the database broken?
4. Are we running out of computer resources (CPU and Memory)?

This project is a **local simulation** of that monitoring system. It runs 6 virtual services on your computer using Docker to show you how a real cloud monitoring stack is configured.

---

## 🏗️ The 6 Components of Your Stack

Here is what each container in the project does:

| Component Name | What it is | Why we use it (Its Job) |
| :--- | :--- | :--- |
| **1. Nginx** | Web Server / Proxy | The front door of your project. It serves the web page (UI Console) to your browser and forwards API commands to the backend. |
| **2. PostgreSQL** | Relational Database | The warehouse. It stores the names and statuses of the Virtual Servers (VSIs) created by the user. |
| **3. Go Backend** | Application Server | The brain. Written in Go, it connects to PostgreSQL, implements the APIs, and exposes statistics (metrics) about the system. |
| **4. Nginx Exporter** | Helper Exporter | Converts Nginx internal stats (like active connection count) into a format Prometheus understands. |
| **5. Prometheus** | Metrics Database | The time-series recorder. Every 15 seconds, it asks the Go Backend and Nginx Exporter for their metrics and saves them. |
| **6. Grafana** | Visualization tool | The painter. It fetches metrics from Prometheus and draws beautiful, real-time charts and meters. |

---

## 🔄 How the Data Flows

When you use the project, this is how data moves:

1. **User Action:** You open `http://localhost:8080` in your browser. Nginx serves you the web page.
2. **Operations:** You click "Provision VSI". 
   - A command goes to the **Go Backend**.
   - The Go Backend writes a new row in **PostgreSQL**.
   - The Go Backend records that a VSI was created and measures how long it took in milliseconds.
3. **Scraping:** **Prometheus** periodically pulls these speed and volume statistics from the Go Backend's `/metrics` page.
4. **Visualizing:** **Grafana** reads those statistics from Prometheus and draws them on charts at `http://localhost:3000`.

---

## 🛠️ What We Implemented For Your Requirements

You shared a document outlining the monitoring requirements for the Sandbox platform. Here is how we built them:

### 1. The Four Golden Signals of Monitoring
DevOps best practices look at 4 key signals of service health:
*   **Latency (Delay):** We measure how long API requests take.
    *   *The Alert:* If the API takes **> 5 seconds**, it is too slow.
    *   *Our Simulation:* We added an **"Inject API Latency"** toggle. When you turn it on, the Go backend purposely sleeps for 5.2 seconds, triggering the alert on the Grafana graph.
*   **Traffic (Demand):** We count the total number of HTTP requests processed per second.
*   **Errors (Failures):** We measure the rate of HTTP errors (like 500 Internal Server Errors).
    *   *The Alert:* If the error rate goes **above 2%**, we alert.
    *   *Our Simulation:* We added an **"Inject Postgres Failure"** toggle. When enabled, it disconnects the database, forcing queries to fail and immediately spiking the error rate to 100%.
*   **Saturation (Fullness):** We measure how close the computer resources are to being full.
    *   *The Alert:* CPU utilization **> 80%** or Memory usage **> 75%**.
    *   *Our Simulation:* We created slider bars on the web page. Adjusting them sends CPU/Memory levels (like 85%) to Prometheus, lighting up red alert indicators in Grafana.

### 2. VSI Quota Checking
*   *Requirement:* Ensure users do not create too many VSIs.
*   *Our Simulation:* We built a **VSI Quota Constraint** toggle. When enabled, it drops the maximum allowed VSIs to `3`. If you try to create a 4th VSI, it fails, returns a `403 Forbidden` error, and increases the **Quota Violations** counter on the Grafana board.

### 3. UI Load Time Monitoring
*   *Requirement:* Alert if the user interface takes **> 5 seconds** to load.
*   *Our Simulation:* The web browser measures the page load time and sends it back to the backend. We added a slider to simulate slow client loading speeds so you can see it graph in Grafana.

### 4. Security & Access
*   *Requirement:* Audit failed logins.
*   *Our Simulation:* Submitting incorrect passwords on the portal logs a failure, incrementing the **Security** panel count in Grafana.

---

## 📖 Essential Commands (And What They Do)

To operate this stack, you use `docker compose` commands. Here is what they mean:

### 1. `docker compose up --build -d`
*   **What it does:** 
    *   `up`: Starts all the containers.
    *   `--build`: Recompiles the Go backend code and rebuilds Nginx configurations.
    *   `-d` (detached mode): Runs the containers in the background, freeing up your terminal.

### 2. `docker compose ps`
*   **What it does:** Lists all running containers, their ports, and tells you if they are healthy.

### 3. `docker compose down`
*   **What it does:** Stops all containers, removes them, and cleans up the virtual Docker network.

---

## 🎯 How to Learn From This Project (Playground Steps)

1. Run the stack: `docker compose up --build -d`
2. Open the **Sandbox Portal Console** at: http://localhost:8080
3. Open the **Grafana Dashboard** at: http://localhost:3000
   *(You will be logged in automatically as an Admin with no password needed)*
4. Under "Sandbox Platform", click the **"Sandbox Services Health & Performance"** dashboard.
5. Place the Sandbox UI and Grafana side-by-side on your screen.
6. Click **"Provision VSI"** 5 times. Watch the **VSI Active Utilization** gauge rise in Grafana.
7. Click **"Inject API Latency"** on the UI. Provision a VSI. Observe the **API Average Latency** chart spike past the 5-second red line in Grafana.
8. Click **"Inject Postgres Failure"** on the UI. Try to click "Execute PostgreSQL Stats Query". You will see an "Outage Detected" error in the UI, and the **PostgreSQL Errors** chart in Grafana will tick upward.
9. Drag the **ROKS CPU Saturation** slider to **90%**. Check Grafana's **Saturation** gauge to see it shift into the red zone.
10. Click **"Clear All Fault Injections"** on the UI and see Grafana return to normal.
