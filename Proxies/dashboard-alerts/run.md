# 🚀 SRE Observability & Observability Playbook

This document is your step-by-step verification script for all **11 panels** on your SRE Dashboard.

---

## 🎛️ Baseline URLs for your Postman / Browser Calls:
* **Frontend Web URL**: `https://go-frontend-route-devops-learning.sandbox-dashboard-alerts-fd8e0ef2d08fcdd9052926b491f21d24-0000.eu-de.containers.appdomain.cloud/`
* **Backend API URL**: `https://go-backend-route-devops-learning.sandbox-dashboard-alerts-fd8e0ef2d08fcdd9052926b491f21d24-0000.eu-de.containers.appdomain.cloud`

---

## 🧪 TEST CASE 1: Steady-State Traffic (Baseline)
* **Goal**: Establish a baseline of successful requests on your metrics charts.
* **Execution (Postman / Browser)**:
  * **Method**: `GET`
  * **URL**: `https://go-backend-route-devops-learning.sandbox-dashboard-alerts-fd8e0ef2d08fcdd9052926b491f21d24-0000.eu-de.containers.appdomain.cloud/`
  * Send 15-20 requests continuously (or run your background loop).
* **Verify on Dashboard**:
  * **Panel 1 (API Health)**: Blue line (`200`) rises to show traffic count.
  * **Panel 2 (API Uptime)**: Stays locked at `100%`.
  * **Panel 3 (API Latency)**: Runs flat at `0.005 seconds` (5ms).

---

## 🧪 TEST CASE 2: API Outage (Database Crash Simulation)
* **Goal**: Verify Uptime drop, HTTP 500 error mapping, and alert triggering.
* **Execution (CLI / Postman)**:
  1. Open PowerShell and shut down the database:
     ```powershell
     oc scale deployment postgres-db --replicas=0 -n devops-learning
     ```
  2. In **Postman**, send 10 requests to get database count:
     * **Method**: `GET`
     * **URL**: `https://go-backend-route-devops-learning.sandbox-dashboard-alerts-fd8e0ef2d08fcdd9052926b491f21d24-0000.eu-de.containers.appdomain.cloud/api/alerts/count`
     * *Observe*: Postman receives `500 Internal Server Error` (Database error).
* **Verify on Dashboard**:
  * **Panel 1 (API Health)**: Red line (`500`) spikes up to `10`.
  * **Panel 2 (API Uptime)**: Drops instantly from `100%` to `0%`.
  * **The Alerts**: Both the *Error Rate > 2%* and *Uptime Alert* fire Red.
* **Recovery (Cleanup)**:
  ```powershell
  oc scale deployment postgres-db --replicas=1 -n devops-learning
  ```

---

## 🧪 TEST CASE 3: Latency SLA Breach (Server Lag Simulation)
* **Goal**: Verify API processing speed alerting when thresholds are crossed.
* **Execution (Postman)**:
  * **Method**: `GET`
  * **URL**: `https://go-backend-route-devops-learning.sandbox-dashboard-alerts-fd8e0ef2d08fcdd9052926b491f21d24-0000.eu-de.containers.appdomain.cloud/`
  * **Params**: Add Key `delay` with Value `1s`.
  * Send 5 requests (each request will take exactly 1.0 second to return).
* **Verify on Dashboard**:
  * **Panel 3 (API Latency)**: Spikes up from `0.005s` to **`1.0 second`**.
  * **The Alert**: Since it crossed the `0.5s` SLA limit, the latency alert triggers.

---

## 🧪 TEST CASE 4: Browser Lag (RUM UI Latency Throttling)
* **Goal**: Measure client-side load time using real browser emulations.
* **Execution (Browser)**:
  1. Open the SRE Homepage in Chrome/Edge.
  2. Open DevTools (`F12`) -> **Network** tab -> select your custom **Demo Network** profile.
  3. Refresh the page 3 times.
* **Verify on Dashboard**:
  * **Panel 4 (Average UI Latency)** & **Panel 5 (95th UI Latency)**: Both lines spike up to match the network profile delay (e.g. `6.0 seconds`).
  * **The Alert**: UI Latency warning alert fires Red.

---

## 🧪 TEST CASE 5: Database Disk Saturation (SQL Write-Stress)
* **Goal**: Stress the storage limit, verify gauge movements, and trigger disk alarms.
* **Execution (CLI)**:
  1. Run this database write-loop:
     ```powershell
     oc exec -it deployment/postgres-db -n devops-learning -- env PGPASSWORD=devopspass psql -U devopsuser -d devopsdb -c "CREATE TABLE IF NOT EXISTS dummy_stress (id SERIAL PRIMARY KEY, data TEXT); INSERT INTO dummy_stress (data) SELECT repeat('A', 1000) FROM generate_series(1, 100000);"
     ```
* **Verify on Dashboard**:
  * **Panel 6 (PostgreSQL Storage)**: Needle rises from `319 MiB` to **`420+ MiB`** (out of `10 GiB` limit).
  * **The Alert**: Disk Overuse Alert triggers.
* **Recovery (Cleanup)**:
  ```powershell
  oc exec -it deployment/postgres-db -n devops-learning -- env PGPASSWORD=devopspass psql -U devopsuser -d devopsdb -c "DROP TABLE dummy_stress;"
  oc exec -it deployment/postgres-db -n devops-learning -- env PGPASSWORD=devopspass psql -U devopsuser -d devopsdb -c "VACUUM FULL;"
  ```

---

## 🧪 TEST CASE 6: Security Intrusion (Brute-Force Login)
* **Goal**: Detect failed login anomaly spikes.
* **Execution (Postman / Browser)**:
  * **Method**: `POST`
  * **URL**: `https://go-backend-route-devops-learning.sandbox-dashboard-alerts-fd8e0ef2d08fcdd9052926b491f21d24-0000.eu-de.containers.appdomain.cloud/api/auth/login`
  * **Body** (raw JSON): Send incorrect details (e.g. `{"username": "hacker", "password": "abc"}`) 10 times quickly.
* **Verify on Dashboard**:
  * **Panel 7 (Failed Logins)**: Shows a sharp spike representing failed logins.
  * **The Alert**: Brute-force Login Attack Warning fires.

---

## 🧪 TEST CASE 7: VM Quota Allocation (Orchestrator Check)
* **Goal**: Monitor active VM provisioning against quota allocations.
* **Execution (Postman / Browser)**:
  * **Method**: `POST`
  * **URL**: `https://go-backend-route-devops-learning.sandbox-dashboard-alerts-fd8e0ef2d08fcdd9052926b491f21d24-0000.eu-de.containers.appdomain.cloud/api/vsi/provision`
  * Send 2-3 clicks.
* **Verify on Dashboard**:
  * **Panel 8 (VSI Utilization)**: The Active VSIs line steps up from `5` to `7` or `8` (approaching the quota limit line at `25`).
* **Recovery (Cleanup)**:
  ```powershell
  oc rollout restart deployment go-backend-app -n devops-learning
  ```

---

## 🧪 TEST CASE 8: Container & Host Resource Saturation (CPU Stress)
* **Goal**: Monitor CPU and Memory consumption under peak server workloads.
* **Execution (Postman)**:
  * **Method**: `GET`
  * **URL**: `https://go-backend-route-devops-learning.sandbox-dashboard-alerts-fd8e0ef2d08fcdd9052926b491f21d24-0000.eu-de.containers.appdomain.cloud/`
  * **Params**: Add Key `stress` with Value `true`.
* **Verify on Dashboard**:
  * **Panel 9 (Pod CPU)**: Gauge pins up to **`8% to 10%`** and turns bright Red.
  * **Panel 10 (Pod Memory)**: Displays a healthy stable limit (~12 MB).
  * **Panel 11 (Host RAM)**: Displays active server RAM usage (e.g. `4.85 GiB` of `16 GiB`).
  * **The Alert**: CPU alert triggers Red.


