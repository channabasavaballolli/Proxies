package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	backendURL := os.Getenv("BACKEND_API_URL")
	if backendURL == "" {
		// Fallback baseline for local testing or auto-binding
		backendURL = "https://go-backend-route-devops-learning.sandbox-dashboard-alerts-fd8e0ef2d08fcdd9052926b491f21d24-0000.eu-de.containers.appdomain.cloud"
	}

	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		data := struct {
			BackendURL string
		}{
			BackendURL: backendURL,
		}
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Template execution error: %s", err)
		}
	})

	log.Printf("Starting Frontend Web Client on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>DevOps Enterprise SRE Dashboard</title>
	<style>
		body { 
			font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
			background-color: #0c0c0e; 
			color: #e1e1e6; 
			margin: 0; 
			padding: 40px; 
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: center;
			min-height: 90vh;
		}
		.card { 
			background: linear-gradient(135deg, #1b1b1f 0%, #111114 100%);
			padding: 35px; 
			border-radius: 12px; 
			border: 1px solid #2e2e33; 
			max-width: 500px; 
			width: 100%;
			box-shadow: 0 8px 32px rgba(0,0,0,0.7); 
		}
		h1 { color: #04d361; margin-top: 0; font-size: 1.8em; text-align: center; }
		h2 { color: #a8a8b3; font-size: 1.1em; margin-bottom: 25px; text-align: center; }
		.metric-group { 
			margin: 20px 0; 
			padding: 15px;
			background-color: #121214;
			border-radius: 6px;
			border: 1px solid #232326; 
		}
		.label { font-weight: bold; color: #8d8d99; font-size: 0.9em; text-transform: uppercase; letter-spacing: 0.5px; }
		.value { font-size: 1.5em; color: #ffffff; font-weight: bold; margin-top: 5px; }
		
		/* Login Form Styles */
		.auth-section {
			margin-top: 25px;
			padding-top: 20px;
			border-top: 1px solid #2e2e33;
		}
		.form-group {
			margin-bottom: 15px;
		}
		label {
			display: block;
			margin-bottom: 5px;
			font-size: 0.85em;
			color: #a8a8b3;
		}
		input {
			width: 100%;
			padding: 10px;
			background-color: #121214;
			border: 1px solid #2e2e33;
			border-radius: 4px;
			color: #fff;
			box-sizing: border-box;
			font-size: 0.95em;
		}
		input:focus {
			border-color: #04d361;
			outline: none;
		}
		.btn { 
			background-color: #04d361; 
			color: #0c0c0e; 
			border: none; 
			padding: 12px 20px; 
			border-radius: 4px; 
			font-weight: bold; 
			cursor: pointer; 
			width: 100%;
			font-size: 1em;
			transition: background-color 0.2s;
			margin-top: 10px;
		}
		.btn:hover { background-color: #00b351; }
		.btn-sec {
			background-color: #202024;
			color: #e1e1e6;
			border: 1px solid #323238;
		}
		.btn-sec:hover { background-color: #29292e; border-color: #04d361; }
		.status-msg {
			margin-top: 10px;
			font-size: 0.9em;
			text-align: center;
			font-weight: bold;
		}
	</style>
</head>
<body>
	<div class="card">
		<h1>DevOps SRE Dashboard</h1>
		<h2>Distributed Client (Frontend Tier)</h2>
		
		<div class="metric-group">
			<div class="label">Live Alert Count (PostgreSQL API)</div>
			<div id="alerts-count" class="value">Loading...</div>
		</div>

		<!-- Interactive VSI Provisioning Block -->
		<div class="metric-group">
			<div class="label">Cloud Orchestrator Actions</div>
			<button id="provision-btn" class="btn btn-sec">Provision Mock VSI</button>
			<div id="vsi-status" class="status-msg" style="color: #a8a8b3;"></div>
		</div>

		<!-- Real Login Form Block (Security & Access Verification) -->
		<div class="auth-section">
			<h3 style="color: #ffffff; margin-top: 0; font-size: 1.1em;">🔐 Authenticate Session</h3>
			<form id="login-form">
				<div class="form-group">
					<label for="username">Username</label>
					<input type="text" id="username" required placeholder="admin">
				</div>
				<div class="form-group">
					<label for="password">Password</label>
					<input type="password" id="password" required placeholder="••••••••">
				</div>
				<button type="submit" class="btn">Access Panel</button>
			</form>
			<div id="auth-status" class="status-msg"></div>
		</div>
	</div>

	<script>
		const API_BASE = "{{.BackendURL}}";

		// 1. Fetch live DB alerts count from backend API every 5 seconds
		async function fetchAlertsCount() {
			try {
				const res = await fetch(API_BASE + "/api/alerts/count");
				if (res.ok) {
					const data = await res.json();
					document.getElementById("alerts-count").innerText = data.count.toLocaleString();
					document.getElementById("alerts-count").style.color = "#04d361";
				} else {
					throw new Error("HTTP error " + res.status);
				}
			} catch (err) {
				document.getElementById("alerts-count").innerText = "Connection Failed";
				document.getElementById("alerts-count").style.color = "#e53e3e";
				console.error("Failed to fetch alerts count:", err);
			}
		}
		
		fetchAlertsCount();
		setInterval(fetchAlertsCount, 5000);

		// 2. Provision Mock VSI action (calls backend to increment VSI count)
		document.getElementById("provision-btn").addEventListener("click", async () => {
			const statusDiv = document.getElementById("vsi-status");
			statusDiv.innerText = "Provisioning VSI...";
			statusDiv.style.color = "#a8a8b3";
			try {
				const res = await fetch(API_BASE + "/api/vsi/provision", {
					method: "POST"
				});
				if (res.ok) {
					const data = await res.json();
					statusDiv.innerText = "VSI Successfully Provisioned! Active: " + data.active_vsis;
					statusDiv.style.color = "#04d361";
				} else {
					throw new Error("HTTP Status " + res.status);
				}
			} catch (err) {
				statusDiv.innerText = "Provisioning Request Failed";
				statusDiv.style.color = "#e53e3e";
				console.error("VSI provisioning failure:", err);
			}
		});

		// 3. Login submission handler (triggers real security failed-login metric logic)
		document.getElementById("login-form").addEventListener("submit", async (e) => {
			e.preventDefault();
			const statusDiv = document.getElementById("auth-status");
			statusDiv.innerText = "Verifying...";
			statusDiv.style.color = "#a8a8b3";

			const u = document.getElementById("username").value;
			const p = document.getElementById("password").value;

			try {
				const res = await fetch(API_BASE + "/api/auth/login", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ username: u, password: p })
				});

				if (res.ok) {
					statusDiv.innerText = "Access Granted! Welcome back.";
					statusDiv.style.color = "#04d361";
					document.getElementById("password").value = "";
				} else if (res.status === 401) {
					statusDiv.innerText = "Access Denied: Invalid Credentials";
					statusDiv.style.color = "#e53e3e";
				} else {
					throw new Error("HTTP error " + res.status);
				}
			} catch (err) {
				statusDiv.innerText = "Auth Service Offline";
				statusDiv.style.color = "#e53e3e";
				console.error("Authentication check error:", err);
			}
		});

		// 4. Real User Monitoring (RUM) Client Latency measurement
		window.addEventListener('load', () => {
			setTimeout(() => {
				const perf = window.performance.timing;
				const loadTimeMs = perf.loadEventEnd - perf.navigationStart;
				
				fetch(API_BASE + "/api/ui-latency", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ latency_ms: loadTimeMs })
				}).catch(err => console.error("RUM Latency reporting failure:", err));
			}, 100);
		});
	</script>
</body>
</html>
`
