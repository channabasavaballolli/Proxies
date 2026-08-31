package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"path"},
	)
	dbQueryDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of Database queries in seconds.",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
		},
	)
	dbErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Total number of database errors.",
		},
	)
	vsiActiveCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "vsi_active_count",
			Help: "Current number of active Virtual Server Instances.",
		},
	)
	vsiQuotaLimit = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "vsi_quota_limit",
			Help: "Maximum allowed Virtual Server Instances.",
		},
	)
	vsiQuotaFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "vsi_quota_failures_total",
			Help: "Total number of VSI creations that failed due to quota limit.",
		},
	)
	securityFailedLogins = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "security_failed_logins_total",
			Help: "Total number of failed login attempts.",
		},
	)
	simulatedCPUSaturation = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "simulated_cpu_saturation_ratio",
			Help: "Simulated CPU utilization ratio of the ROKS cluster (0 to 100).",
		},
	)
	simulatedMemorySaturation = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "simulated_memory_saturation_ratio",
			Help: "Simulated memory utilization ratio of the ROKS cluster (0 to 100).",
		},
	)
	uiPageLoadDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "ui_page_load_seconds",
			Help:    "Simulated or reported UI page load latency in seconds.",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(dbQueryDuration)
	prometheus.MustRegister(dbErrorsTotal)
	prometheus.MustRegister(vsiActiveCount)
	prometheus.MustRegister(vsiQuotaLimit)
	prometheus.MustRegister(vsiQuotaFailures)
	prometheus.MustRegister(securityFailedLogins)
	prometheus.MustRegister(simulatedCPUSaturation)
	prometheus.MustRegister(simulatedMemorySaturation)
	prometheus.MustRegister(uiPageLoadDuration)

	// Set defaults
	vsiQuotaLimit.Set(10)
	simulatedCPUSaturation.Set(45.0)
	simulatedMemorySaturation.Set(50.0)
}

// Simulation state mutex and variables
var (
	simMutex          sync.Mutex
	latencySimulated  = false
	dbErrorSimulated  = false
	vsiQuotaSimulated = false
	simulatedCPUVal   = 45.0
	simulatedMemoryVal  = 50.0
	simulatedUILoadVal  = 1.2
)

var db *sql.DB

func main() {
	// Setup DB connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "sandbox_user")
	dbPass := getEnv("DB_PASSWORD", "sandbox_password")
	dbName := getEnv("DB_NAME", "sandbox_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	var err error
	// Retry connection because Postgres might still be starting
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to PostgreSQL")
				break
			}
		}
		log.Printf("Waiting for database connection (attempt %d/10)... error: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize tables
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS vsis (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Update active VSI count based on existing rows
	updateActiveVSICount()

	// Set up router and wrap endpoints with telemetry
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	// Core Sandbox APIs
	mux.HandleFunc("/api/vsi", instrumentHandler("/api/vsi", handleVSI))
	mux.HandleFunc("/api/login", instrumentHandler("/api/login", handleLogin))
	mux.HandleFunc("/api/db-query", instrumentHandler("/api/db-query", handleDBQuery))
	mux.HandleFunc("/api/ui-load-record", instrumentHandler("/api/ui-load-record", handleUILoadRecord))

	// Fault Injection APIs
	mux.HandleFunc("/api/simulate/latency", handleSimulateLatency)
	mux.HandleFunc("/api/simulate/db-error", handleSimulateDBError)
	mux.HandleFunc("/api/simulate/cpu-saturation", handleSimulateCPUSaturation)
	mux.HandleFunc("/api/simulate/vsi-quota", handleSimulateVSIQuota)
	mux.HandleFunc("/api/status", handleStatus)

	port := getEnv("PORT", "8080")
	log.Printf("Starting Go Sandbox Backend Service on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}

// Wrapper for telemetry and simulated latency
func instrumentHandler(path string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Inject simulated latency if active
		simMutex.Lock()
		activeLatency := latencySimulated
		simMutex.Unlock()

		if activeLatency {
			// Delays API response by > 5 seconds to trigger alert (e.g. 5.2 seconds)
			time.Sleep(5200 * time.Millisecond)
		} else {
			// Add a minor natural random latency (5ms to 50ms)
			time.Sleep(time.Duration(5+rand.Intn(45)) * time.Millisecond)
		}

		// Helper wrapper to capture status code
		writer := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		handler(writer, r)

		// Record HTTP metrics
		duration := time.Since(start).Seconds()
		statusStr := strconv.Itoa(writer.statusCode)
		httpRequestsTotal.WithLabelValues(path, r.Method, statusStr).Inc()
		httpRequestDuration.WithLabelValues(path).Observe(duration)
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (sr *statusResponseWriter) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

// --- API Handlers ---

type VSI struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func handleVSI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Validate DB state
	simMutex.Lock()
	dbFail := dbErrorSimulated
	simMutex.Unlock()

	if dbFail {
		dbErrorsTotal.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "PostgreSQL database connection reset"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, name, status, created_at FROM vsis ORDER BY id DESC")
		if err != nil {
			dbErrorsTotal.Inc()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		defer rows.Close()

		vsis := []VSI{}
		for rows.Next() {
			var v VSI
			if err := rows.Scan(&v.ID, &v.Name, &v.Status, &v.CreatedAt); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			vsis = append(vsis, v)
		}
		json.NewEncoder(w).Encode(vsis)

	case http.MethodPost:
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
			return
		}

		// Check Quota
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM vsis").Scan(&count)
		if err != nil {
			dbErrorsTotal.Inc()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		simMutex.Lock()
		quotaSim := vsiQuotaSimulated
		simMutex.Unlock()

		limit := 10
		if quotaSim {
			limit = 3 // Reduced limit for simulation
		}
		vsiQuotaLimit.Set(float64(limit))

		if count >= limit {
			vsiQuotaFailures.Inc()
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Quota limit reached: Maximum %d VSIs allowed.", limit)})
			return
		}

		// Insert VSI
		status := "Active"
		var id int
		err = db.QueryRow("INSERT INTO vsis(name, status) VALUES($1, $2) RETURNING id", req.Name, status).Scan(&id)
		if err != nil {
			dbErrorsTotal.Inc()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updateActiveVSICount()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(VSI{ID: id, Name: req.Name, Status: status, CreatedAt: time.Now()})

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM vsis WHERE id = $1", id)
		if err != nil {
			dbErrorsTotal.Inc()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updateActiveVSICount()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "VSI deleted successfully"})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// For simulation: credentials admin/password
	if req.Username == "admin" && req.Password == "password" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "Success", "token": "simulated_token_123"})
	} else {
		securityFailedLogins.Inc()
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"status": "Unauthorized", "message": "Failed login attempt: Invalid credentials"})
	}
}

func handleDBQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	simMutex.Lock()
	dbFail := dbErrorSimulated
	simMutex.Unlock()

	if dbFail {
		dbErrorsTotal.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "PostgreSQL database connection refused"})
		return
	}

	start := time.Now()
	// Run raw count query to simulate DB operation
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM vsis").Scan(&count)

	duration := time.Since(start).Seconds()
	dbQueryDuration.Observe(duration)

	if err != nil {
		dbErrorsTotal.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"query":           "SELECT COUNT(*) FROM vsis",
		"rows_count":      count,
		"execution_time":  fmt.Sprintf("%.5f seconds", duration),
		"status":          "Success",
	})
}

func handleUILoadRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		LoadTime float64 `json:"load_time"` // in seconds
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uiPageLoadDuration.Observe(req.LoadTime)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "recorded",
		"load_time": req.LoadTime,
	})
}

// --- Fault Injection Handler Endpoints ---

func handleSimulateLatency(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Enable bool `json:"enable"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	simMutex.Lock()
	latencySimulated = req.Enable
	simMutex.Unlock()

	log.Printf("Fault Injection: Latency simulation set to %v", req.Enable)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "latencySimulated": req.Enable})
}

func handleSimulateDBError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Enable bool `json:"enable"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	simMutex.Lock()
	dbErrorSimulated = req.Enable
	simMutex.Unlock()

	log.Printf("Fault Injection: DB Error simulation set to %v", req.Enable)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "dbErrorSimulated": req.Enable})
}

func handleSimulateCPUSaturation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CPU    float64 `json:"cpu"`
		Memory float64 `json:"memory"`
		UILoad float64 `json:"ui_load"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	simMutex.Lock()
	simulatedCPUVal = req.CPU
	simulatedMemoryVal = req.Memory
	simulatedUILoadVal = req.UILoad
	simMutex.Unlock()

	// Update gauges
	simulatedCPUSaturation.Set(req.CPU)
	simulatedMemorySaturation.Set(req.Memory)

	log.Printf("Fault Injection: Saturation set to CPU: %.1f%%, Memory: %.1f%%, UI Load Time: %.1fs", req.CPU, req.Memory, req.UILoad)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":            "success",
		"simulatedCPU":      req.CPU,
		"simulatedMemory":   req.Memory,
		"simulatedUILoad":   req.UILoad,
	})
}

func handleSimulateVSIQuota(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Enable bool `json:"enable"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	simMutex.Lock()
	vsiQuotaSimulated = req.Enable
	simMutex.Unlock()

	limit := 10
	if req.Enable {
		limit = 3
	}
	vsiQuotaLimit.Set(float64(limit))

	log.Printf("Fault Injection: VSI Quota Limit simulation set to limit %d (enabled: %v)", limit, req.Enable)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "vsiQuotaSimulated": req.Enable, "limit": limit})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	simMutex.Lock()
	defer simMutex.Unlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"latencySimulated":   latencySimulated,
		"dbErrorSimulated":   dbErrorSimulated,
		"vsiQuotaSimulated":  vsiQuotaSimulated,
		"simulatedCPU":       simulatedCPUVal,
		"simulatedMemory":    simulatedMemoryVal,
		"simulatedUILoad":    simulatedUILoadVal,
	})
}

// --- Helper Functions ---

func updateActiveVSICount() {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM vsis").Scan(&count)
	if err == nil {
		vsiActiveCount.Set(float64(count))
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
