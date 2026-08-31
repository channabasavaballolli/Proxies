package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// LogEntry defines our structured JSON logging format
type LogEntry struct {
	Timestamp string  `json:"timestamp"`
	Level     string  `json:"level"`
	Component string  `json:"component"`
	Message   string  `json:"message"`
	Duration  float64 `json:"duration_ms,omitempty"`
	Query     string  `json:"query,omitempty"`
	Status    int     `json:"status,omitempty"`
	Error     string  `json:"error,omitempty"`
}

// Custom Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_http_requests_total",
			Help: "Total number of HTTP requests processed by the application.",
		},
		[]string{"path", "status", "method"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_app_http_request_duration_seconds",
			Help:    "Histogram of response latency for HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	dbQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_db_queries_total",
			Help: "Total number of database queries executed.",
		},
		[]string{"query_type", "status"},
	)

	dbQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_app_db_query_duration_seconds",
			Help:    "Histogram of database query execution time in seconds.",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0},
		},
		[]string{"query_type"},
	)

	vsiConnectionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_vsi_connections_total",
			Help: "Total number of connection attempts to VSI / external API.",
		},
		[]string{"endpoint", "status"},
	)

	vsiConnectionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_app_vsi_connection_duration_seconds",
			Help:    "Histogram of connection latency to external VSI endpoint in seconds.",
			Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
		},
		[]string{"endpoint"},
	)

	uiClientDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "go_app_ui_client_load_duration_seconds",
			Help:    "Real client-side browser page load latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
	)

	securityFailedLogins = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_security_failed_logins_total",
			Help: "Total number of failed user authentication attempts.",
		},
		[]string{"username", "status"},
	)

	vsiActiveCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_app_vsi_active_count",
			Help: "Number of active virtual server instances.",
		},
	)

	vsiQuotaLimit = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_app_vsi_quota_limit",
			Help: "Maximum quota limit of virtual server instances.",
		},
	)
)

var (
	vsiMutex   sync.Mutex
	activeVSIs = 5
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(dbQueriesTotal)
	prometheus.MustRegister(dbQueryDuration)
	prometheus.MustRegister(vsiConnectionsTotal)
	prometheus.MustRegister(vsiConnectionDuration)
	prometheus.MustRegister(uiClientDuration)
	prometheus.MustRegister(securityFailedLogins)
	prometheus.MustRegister(vsiActiveCount)
	prometheus.MustRegister(vsiQuotaLimit)
}

// writeLog helper to output structured JSON logs to stdout
func writeLog(level, component, message string, durationMs float64, query string, status int, errVal error) {
	errStr := ""
	if errVal != nil {
		errStr = errVal.Error()
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Component: component,
		Message:   message,
		Duration:  durationMs,
		Query:     query,
		Status:    status,
		Error:     errStr,
	}

	jsonBytes, err := json.Marshal(entry)
	if err == nil {
		fmt.Println(string(jsonBytes))
	} else {
		log.Printf("Failed to marshal log entry: %s", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize VSI Quota telemetry metrics
	vsiQuotaLimit.Set(25.0)
	vsiActiveCount.Set(5.0)

	// Retrieve Database connection details
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres-service" // Matches the internal service name in OpenShift
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "devopsuser"
	}
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "devopspass"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "devopsdb"
	}

	// Connection string format for github.com/lib/pq
	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
		dbHost, dbUser, dbPass, dbName)

	var db *sql.DB
	var err error

	// 1. Connection Retry Loop (DevOps Best Practice)
	writeLog("INFO", "system", "Connecting to PostgreSQL database...", 0, "", 0, nil)
	for i := 1; i <= 15; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			writeLog("INFO", "system", fmt.Sprintf("Successfully connected to database after %d attempts", i), 0, "", 0, nil)
			break
		}

		writeLog("WARNING", "system", fmt.Sprintf("Database connection attempt %d/15 failed. Retrying in 2 seconds...", i), 0, "", 0, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		writeLog("ERROR", "system", "Failed to connect to database after 15 attempts. Proceeding in disconnected state.", 0, "", 0, err)
	}
	if db != nil {
		defer db.Close()
	}

	// 2. Database Schema Migration
	migrateQuery := `
	CREATE TABLE IF NOT EXISTS system_alerts (
		id SERIAL PRIMARY KEY,
		source VARCHAR(50) NOT NULL,
		level VARCHAR(10) NOT NULL,
		message TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if db != nil {
		_, err = db.Exec(migrateQuery)
		if err != nil {
			writeLog("ERROR", "postgres", "Failed to execute database migration schema (will retry dynamically)", 0, migrateQuery, 0, err)
		} else {
			writeLog("INFO", "postgres", "Database schema migration executed successfully", 0, migrateQuery, 0, nil)
		}
	}

	// 3. PostgreSQL Simulated Load Loop (Real Queries & Timings)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			
			// Insert operation
			insertStart := time.Now()
			insertQuery := "INSERT INTO system_alerts (source, level, message) VALUES ('scheduler', 'info', 'Routine system health check passed')"
			_, insertErr := db.Exec(insertQuery)
			duration := float64(time.Since(insertStart).Microseconds()) / 1000.0 // convert to milliseconds
			dbQueryDuration.WithLabelValues("INSERT").Observe(time.Since(insertStart).Seconds())

			if insertErr != nil {
				dbQueriesTotal.WithLabelValues("INSERT", "failed").Inc()
				writeLog("ERROR", "postgres", "Failed to insert record", duration, insertQuery, 0, insertErr)
			} else {
				dbQueriesTotal.WithLabelValues("INSERT", "success").Inc()
				writeLog("INFO", "postgres", "Alert record inserted", duration, insertQuery, 0, nil)
			}

			// Select operation
			selectStart := time.Now()
			selectQuery := "SELECT COUNT(*) FROM system_alerts"
			var count int
			selectErr := db.QueryRow(selectQuery).Scan(&count)
			duration = float64(time.Since(selectStart).Microseconds()) / 1000.0
			dbQueryDuration.WithLabelValues("SELECT").Observe(time.Since(selectStart).Seconds())

			if selectErr != nil {
				dbQueriesTotal.WithLabelValues("SELECT", "failed").Inc()
				writeLog("ERROR", "postgres", "Failed to query database count", duration, selectQuery, 0, selectErr)
			} else {
				dbQueriesTotal.WithLabelValues("SELECT", "success").Inc()
				writeLog("INFO", "postgres", fmt.Sprintf("Alerts table count: %d", count), duration, selectQuery, 0, nil)
			}
		}
	}()

	// 4. VSI / External Connection Loop (Real HTTP Requests & Latency)
	go func() {
		client := http.Client{
			Timeout: 4 * time.Second, // Timeout to record connection failures
		}
		externalURL := "https://www.githubstatus.com/api/v2/status.json"

		for {
			time.Sleep(8 * time.Second)
			
			connStart := time.Now()
			vsiConnectionsTotal.WithLabelValues(externalURL, "attempt").Inc()
			
			resp, httpErr := client.Get(externalURL)
			duration := float64(time.Since(connStart).Microseconds()) / 1000.0
			vsiConnectionDuration.WithLabelValues(externalURL).Observe(time.Since(connStart).Seconds())

			if httpErr != nil {
				vsiConnectionsTotal.WithLabelValues(externalURL, "failed").Inc()
				writeLog("ERROR", "vsi", "Failed to reach external VSI API endpoint", duration, "", 0, httpErr)
			} else {
				resp.Body.Close()
				vsiConnectionsTotal.WithLabelValues(externalURL, "success").Inc()
				writeLog("INFO", "vsi", fmt.Sprintf("Successfully checked external VSI status. HTTP Status: %d", resp.StatusCode), duration, "", resp.StatusCode, nil)
			}
		}
	}()

	// Define HTTP Route Handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues("/", r.Method))
		defer timer.ObserveDuration() //

		// 1. Chaos Latency Delay Injection
		if delayStr := r.URL.Query().Get("delay"); delayStr != "" {
			if duration, err := time.ParseDuration(delayStr); err == nil {
				writeLog("WARNING", "system", fmt.Sprintf("Chaos testing: Injecting artificial delay of %s", delayStr), 0, "", 0, nil)
				time.Sleep(duration)
			}
		}

		// 2. Chaos CPU Stress Ingress
		if stressStr := r.URL.Query().Get("stress"); stressStr != "" {
			writeLog("WARNING", "system", "Chaos testing: Spawning CPU stress busy-loop for 15 seconds", 0, "", 0, nil)
			// Run busy CPU loops on multiple cores
			for c := 0; c < 2; c++ {
				go func() {
					endTime := time.Now().Add(15 * time.Second)
					for time.Now().Before(endTime) {
						// Intensive mathematical operations
						_ = 3.14159 * 2.71828 / 1.41421
					}
				}()
			}
		}

		// Return simple API status instead of monolithic HTML
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "online", "service": "backend-api"})
		httpRequestsTotal.WithLabelValues("/", "200", r.Method).Inc()
	})

	// Helper to handle CORS preflight and headers
	enableCors := func(w *http.ResponseWriter, r *http.Request) bool {
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
		(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			(*w).WriteHeader(http.StatusOK)
			return true
		}
		return false
	}

	// Endpoint to retrieve alert count
	http.HandleFunc("/api/alerts/count", func(w http.ResponseWriter, r *http.Request) {
		if enableCors(&w, r) {
			return
		}
		var totalAlerts int
		err := db.QueryRow("SELECT COUNT(*) FROM system_alerts").Scan(&totalAlerts)
		if err != nil {
			httpRequestsTotal.WithLabelValues("/api/alerts/count", "500", r.Method).Inc()
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"count": totalAlerts})
		httpRequestsTotal.WithLabelValues("/api/alerts/count", "200", r.Method).Inc()
	})

	// Endpoint to handle login verification (Security & Access telemetry)
	http.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if enableCors(&w, r) {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Real validation check
		if req.Username == "admin" && req.Password == "devopspassword" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "success", "token": "mock-jwt-token"})
			writeLog("INFO", "auth", fmt.Sprintf("Successful login for user: %s", req.Username), 0, "", http.StatusOK, nil)
			httpRequestsTotal.WithLabelValues("/api/auth/login", "200", r.Method).Inc()
		} else {
			// Increment security failed logins metric
			securityFailedLogins.WithLabelValues(req.Username, "failed").Inc()
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			writeLog("WARNING", "auth", fmt.Sprintf("Failed login attempt for user: %s", req.Username), 0, "", http.StatusUnauthorized, nil)
			httpRequestsTotal.WithLabelValues("/api/auth/login", "401", r.Method).Inc()
		}
	})

	// Endpoint to provision dummy VSIs (Virtual Server Instances quota verification)
	http.HandleFunc("/api/vsi/provision", func(w http.ResponseWriter, r *http.Request) {
		if enableCors(&w, r) {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		vsiMutex.Lock()
		activeVSIs++
		vsiActiveCount.Set(float64(activeVSIs))
		current := activeVSIs
		vsiMutex.Unlock()

		writeLog("INFO", "vsi", fmt.Sprintf("Provisioned virtual instance. Total active: %d", current), 0, "", http.StatusOK, nil)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      200,
			"active_vsis": current,
		})
		httpRequestsTotal.WithLabelValues("/api/vsi/provision", "200", r.Method).Inc()
	})

	// Endpoint for browser to report client-side page load latency (Real User Monitoring)
	http.HandleFunc("/api/ui-latency", func(w http.ResponseWriter, r *http.Request) {
		if enableCors(&w, r) {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type LatencyReport struct {
			LatencyMs float64 `json:"latency_ms"`
		}

		var report LatencyReport
		err := json.NewDecoder(r.Body).Decode(&report)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Convert milliseconds to seconds for Prometheus
		latencySec := report.LatencyMs / 1000.0
		uiClientDuration.Observe(latencySec)

		writeLog("INFO", "http", fmt.Sprintf("Reported real client-side UI page load latency: %.3f seconds", latencySec), report.LatencyMs, "", http.StatusOK, nil)
		
		w.WriteHeader(http.StatusNoContent)
		httpRequestsTotal.WithLabelValues("/api/ui-latency", "204", r.Method).Inc()
	})

	// Expose raw Prometheus scrapable metrics page
	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("Starting Go Telemetry App on port %s...\n", port)
	writeLog("INFO", "system", fmt.Sprintf("Starting Go Telemetry App on port %s...", port), 0, "", 0, nil)

	addr := ":" + port
	if err := http.ListenAndServe(addr, nil); err != nil {
		writeLog("CRITICAL", "system", "Web server crashed", 0, "", 0, err)
		log.Fatalf("Error starting server: %s", err)
	}
}
