package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/grid-x/modbus"
)

// Configuration structures
type RegisterConfig struct {
	Address  uint16 `json:"address"`
	Quantity uint16 `json:"quantity"`
	DataType string `json:"dataType"` // "float32", "uint16", "int16"
	Name     string `json:"name"`
}

type ModbusServerConfig struct {
	ID        string           `json:"id"`
	Host      string           `json:"host"`
	Port      int              `json:"port"`
	UnitID    byte             `json:"unitId"`
	Registers []RegisterConfig `json:"registers"`
	Enabled   bool             `json:"enabled"`
}

type RegisterValue struct {
	ServerID  string      `json:"serverId"`
	Name      string      `json:"name"`
	Address   uint16      `json:"address"`
	Value     interface{} `json:"value"`
	Timestamp time.Time   `json:"timestamp"`
	Error     string      `json:"error,omitempty"`
}

type ModbusGateway struct {
	servers   map[string]*ModbusServerConfig
	clients   map[string]modbus.Client
	handlers  map[string]*modbus.TCPClientHandler
	cache     map[string][]RegisterValue
	mu        sync.RWMutex
	wsClients map[*websocket.Conn]bool
	wsMu      sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func NewModbusGateway() *ModbusGateway {
	return &ModbusGateway{
		servers:   make(map[string]*ModbusServerConfig),
		clients:   make(map[string]modbus.Client),
		handlers:  make(map[string]*modbus.TCPClientHandler),
		cache:     make(map[string][]RegisterValue),
		wsClients: make(map[*websocket.Conn]bool),
	}
}

func main() {
	gateway := NewModbusGateway()

	// Load mock configuration
	loadMockConfig(gateway)

	// Start background polling
	go gateway.startBackgroundPolling()

	// Setup HTTP handlers
	http.HandleFunc("/ws", gateway.handleWebSocket)
	http.HandleFunc("/api/servers", gateway.handleGetServers)
	http.HandleFunc("/api/servers/add", gateway.handleAddServer)
	http.HandleFunc("/api/servers/remove", gateway.handleRemoveServer)
	http.HandleFunc("/api/read", gateway.handleRead)
	http.HandleFunc("/api/write", gateway.handleWrite)
	http.HandleFunc("/api/cache", gateway.handleGetCache)

	log.Println("🚀 Modbus Gateway Started")
	log.Println("📡 WebSocket: ws://localhost:8080/ws")
	log.Println("🌐 HTTP API: http://localhost:8080/api/*")
	log.Println("=" + string(make([]byte, 50)))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}

// Load mock configuration with multiple servers
func loadMockConfig(gw *ModbusGateway) {
	configs := []ModbusServerConfig{
		{
			ID:      "chiller-1",
			Host:    "localhost",
			Port:    502,
			UnitID:  1,
			Enabled: true,
			Registers: []RegisterConfig{
				{Address: 0, Quantity: 2, DataType: "float32", Name: "Leaving_Temp_Settings"},
				{Address: 2, Quantity: 2, DataType: "float32", Name: "Entering_Temp"},
				{Address: 4, Quantity: 2, DataType: "float32", Name: "Leaving_Temp"},
			},
		},
		{
			ID:      "chiller-2",
			Host:    "localhost",
			Port:    504, // Different port for 2nd server
			UnitID:  1,
			Enabled: false, // Disabled by default
			Registers: []RegisterConfig{
				{Address: 0, Quantity: 2, DataType: "float32", Name: "Supply_Temp"},
				{Address: 2, Quantity: 2, DataType: "float32", Name: "Return_Temp"},
			},
		},
		{
			ID:      "sensors",
			Host:    "localhost",
			Port:    505,
			UnitID:  1,
			Enabled: false,
			Registers: []RegisterConfig{
				{Address: 0, Quantity: 1, DataType: "uint16", Name: "Pressure_Sensor"},
				{Address: 1, Quantity: 1, DataType: "uint16", Name: "Flow_Sensor"},
				{Address: 2, Quantity: 1, DataType: "int16", Name: "Temperature_Sensor"},
			},
		},
	}

	for _, cfg := range configs {
		gw.addServer(cfg)
	}

	log.Printf("✓ Loaded %d mock server configurations\n", len(configs))
}

// Add server and connect
func (gw *ModbusGateway) addServer(cfg ModbusServerConfig) error {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	// Always add to servers list
	gw.servers[cfg.ID] = &cfg

	if !cfg.Enabled {
		log.Printf("➕ Server '%s' added (disabled)\n", cfg.ID)
		return nil
	}

	// Create handler
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	handler.Timeout = 5 * time.Second
	handler.SlaveID = cfg.UnitID

	// Try to connect
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := handler.Connect(ctx); err != nil {
		log.Printf("⚠️  Server '%s' connection failed: %v (will retry in background)\n", cfg.ID, err)
		// Don't return error - will retry in background
		return nil
	}

	client := modbus.NewClient(handler)

	gw.clients[cfg.ID] = client
	gw.handlers[cfg.ID] = handler
	gw.cache[cfg.ID] = make([]RegisterValue, 0)

	log.Printf("✓ Server '%s' connected to %s:%d\n", cfg.ID, cfg.Host, cfg.Port)
	return nil
}

// Background polling loop
func (gw *ModbusGateway) startBackgroundPolling() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		gw.mu.RLock()
		serverIDs := make([]string, 0, len(gw.servers))
		for id := range gw.servers {
			serverIDs = append(serverIDs, id)
		}
		gw.mu.RUnlock()

		for _, serverID := range serverIDs {
			go gw.pollServer(serverID)
		}
	}
}

// Poll single server
func (gw *ModbusGateway) pollServer(serverID string) {
	gw.mu.RLock()
	server, exists := gw.servers[serverID]
	gw.mu.RUnlock()

	if !exists || !server.Enabled {
		return
	}

	// Check if client exists, if not try to reconnect
	gw.mu.RLock()
	client, connected := gw.clients[serverID]
	gw.mu.RUnlock()

	if !connected {
		// Try to reconnect
		gw.reconnectServer(serverID)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	values := make([]RegisterValue, 0)

	for _, reg := range server.Registers {
		result, err := client.ReadHoldingRegisters(ctx, reg.Address, reg.Quantity)

		value := RegisterValue{
			ServerID:  serverID,
			Name:      reg.Name,
			Address:   reg.Address,
			Timestamp: time.Now(),
		}

		if err != nil {
			value.Error = err.Error()
			// On error, try to reconnect
			go gw.reconnectServer(serverID)
		} else {
			value.Value = parseValue(result, reg.DataType)
		}

		values = append(values, value)
	}

	// Update cache
	gw.mu.Lock()
	gw.cache[serverID] = values
	gw.mu.Unlock()

	// Broadcast to WebSocket clients
	gw.broadcastToWebSocket(map[string]interface{}{
		"type":     "data_update",
		"serverId": serverID,
		"values":   values,
	})
}

// Reconnect to server
func (gw *ModbusGateway) reconnectServer(serverID string) {
	gw.mu.RLock()
	server, exists := gw.servers[serverID]
	gw.mu.RUnlock()

	if !exists || !server.Enabled {
		return
	}

	// Close existing connection if any
	gw.mu.Lock()
	if handler, exists := gw.handlers[serverID]; exists {
		_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		handler.Close()
		cancel()
		delete(gw.clients, serverID)
		delete(gw.handlers, serverID)
	}
	gw.mu.Unlock()

	// Try to create new connection
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", server.Host, server.Port))
	handler.Timeout = 5 * time.Second
	handler.SlaveID = server.UnitID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := handler.Connect(ctx); err != nil {
		log.Printf("⚠️  Reconnect failed for '%s': %v\n", serverID, err)
		return
	}

	client := modbus.NewClient(handler)

	gw.mu.Lock()
	gw.clients[serverID] = client
	gw.handlers[serverID] = handler
	gw.mu.Unlock()

	log.Printf("✓ Server '%s' reconnected successfully\n", serverID)
}

// Parse value based on data type
func parseValue(data []byte, dataType string) interface{} {
	switch dataType {
	case "float32":
		if len(data) >= 4 {
			bits := binary.BigEndian.Uint32(data)
			return math.Float32frombits(bits)
		}
	case "uint16":
		if len(data) >= 2 {
			return binary.BigEndian.Uint16(data)
		}
	case "int16":
		if len(data) >= 2 {
			return int16(binary.BigEndian.Uint16(data))
		}
	}
	return nil
}

// WebSocket handler
func (gw *ModbusGateway) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	gw.wsMu.Lock()
	gw.wsClients[conn] = true
	gw.wsMu.Unlock()

	log.Println("📱 New WebSocket client connected")

	// Send initial cache
	gw.mu.RLock()
	cache := gw.cache
	gw.mu.RUnlock()

	conn.WriteJSON(map[string]interface{}{
		"type":  "initial_data",
		"cache": cache,
	})

	// Listen for messages from client
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		// Handle client commands
		gw.handleWebSocketCommand(conn, msg)
	}

	gw.wsMu.Lock()
	delete(gw.wsClients, conn)
	gw.wsMu.Unlock()

	log.Println("📱 WebSocket client disconnected")
}

// Handle WebSocket commands
func (gw *ModbusGateway) handleWebSocketCommand(conn *websocket.Conn, msg map[string]interface{}) {
	cmdType, _ := msg["type"].(string)

	switch cmdType {
	case "write":
		serverID, _ := msg["serverId"].(string)
		address, _ := msg["address"].(float64)
		value, _ := msg["value"].(float64)
		dataType, _ := msg["dataType"].(string)

		err := gw.writeRegister(serverID, uint16(address), value, dataType)

		response := map[string]interface{}{
			"type":     "write_response",
			"serverId": serverID,
			"address":  address,
			"success":  err == nil,
		}

		if err != nil {
			response["error"] = err.Error()
		}

		conn.WriteJSON(response)
	}
}

// Write register
func (gw *ModbusGateway) writeRegister(serverID string, address uint16, value float64, dataType string) error {
	gw.mu.RLock()
	client, exists := gw.clients[serverID]
	gw.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var data []byte

	switch dataType {
	case "float32":
		bits := math.Float32bits(float32(value))
		data = make([]byte, 4)
		binary.BigEndian.PutUint32(data, bits)
		_, err := client.WriteMultipleRegisters(ctx, address, 2, data)
		return err
	case "uint16":
		data = make([]byte, 2)
		binary.BigEndian.PutUint16(data, uint16(value))
		_, err := client.WriteMultipleRegisters(ctx, address, 1, data)
		return err
	case "int16":
		data = make([]byte, 2)
		binary.BigEndian.PutUint16(data, uint16(int16(value)))
		_, err := client.WriteMultipleRegisters(ctx, address, 1, data)
		return err
	}

	return fmt.Errorf("unsupported data type")
}

// Broadcast to all WebSocket clients
func (gw *ModbusGateway) broadcastToWebSocket(data interface{}) {
	gw.wsMu.RLock()
	defer gw.wsMu.RUnlock()

	for conn := range gw.wsClients {
		conn.WriteJSON(data)
	}
}

// HTTP API Handlers
func (gw *ModbusGateway) handleGetServers(w http.ResponseWriter, r *http.Request) {
	gw.mu.RLock()
	defer gw.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gw.servers)
}

func (gw *ModbusGateway) handleAddServer(w http.ResponseWriter, r *http.Request) {
	var cfg ModbusServerConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := gw.addServer(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "added"})
}

func (gw *ModbusGateway) handleRemoveServer(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("id")

	gw.mu.Lock()
	defer gw.mu.Unlock()

	if handler, exists := gw.handlers[serverID]; exists {
		_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		handler.Close()
	}

	delete(gw.servers, serverID)
	delete(gw.clients, serverID)
	delete(gw.handlers, serverID)
	delete(gw.cache, serverID)

	json.NewEncoder(w).Encode(map[string]string{"status": "removed"})
}

func (gw *ModbusGateway) handleRead(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("serverId")

	gw.pollServer(serverID)

	gw.mu.RLock()
	values := gw.cache[serverID]
	gw.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(values)
}

func (gw *ModbusGateway) handleWrite(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerID string  `json:"serverId"`
		Address  uint16  `json:"address"`
		Value    float64 `json:"value"`
		DataType string  `json:"dataType"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := gw.writeRegister(req.ServerID, req.Address, req.Value, req.DataType)

	response := map[string]interface{}{
		"success": err == nil,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (gw *ModbusGateway) handleGetCache(w http.ResponseWriter, r *http.Request) {
	gw.mu.RLock()
	defer gw.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gw.cache)
}
