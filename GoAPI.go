package main

import (
	"context"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/proto"
	"gopkg.in/ini.v1"
	"github.com/derbanane/xmbig/xmbig" // Korrekter Importpfad
)

// Globale Variablen
var (
	db          *sql.DB
	config      Config
	minerConns  sync.Map // Verwaltet die TCP-Verbindungen zu den Minern
	tcpListener net.Listener
)

// MinerConfig represents the configuration for a miner
type MinerConfig struct {
	MinerID     string `json:"minerID"`
	PoolAddress string `json:"poolAddress"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Algorithm   string `json:"algorithm"`
	AutoSwitch  bool   `json:"autoSwitch"`
	TorEnabled  bool   `json:"torEnabled"`
	ExtraParams string `json:"extraParams"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Name     string `ini:"name"`
}

// Config represents the application configuration
type Config struct {
	Database DatabaseConfig `ini:"database"`
}

// globals
var (
	//Globals
	clientConnector ClientConnector
)
type ClientConnector struct {
    ServerAddress string
}

func (cc *ClientConnector) SendCommand(minerID string, command *xmbig.ControlCommand) error {
		// Get Connection for this miner Id
	connRaw, ok := minerConns.Load(minerID)
	if !ok {
		return fmt.Errorf("unable to find miner with this MinerId")
	}
	conn, _ := connRaw.(net.Conn)

    // Serialize the message
    data, err := proto.Marshal(command)
    if err != nil {
        return fmt.Errorf("failed to marshal ControlCommand: %w", err)
    }

    // Prepend the message with its length
    length := uint32(len(data))
    err = binary.Write(conn, binary.BigEndian, length)
    if err != nil {
        return fmt.Errorf("failed to write message length: %w", err)
    }

    // Send the message
    _, err = conn.Write(data)
    if err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }

    fmt.Printf("Successfully send command to miner %s", minerID)

    return nil
}
func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Could not load .env file: %v", err)
	}

	// Load configuration from file
	cfg := Config{}
	err = loadConfig("config.ini", &cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	// Connect to the database
	db, err = connectToDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()
	println("Connected to the Database")

	   // Initialisiere die TCP Verbindung
	go tcpServerConnector()

	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.POST("/api/miner/config", func(c *gin.Context) {
		setMinerConfig(c, cfg.Database, &minerConns)
	})
	router.POST("/api/miner/command", sendMinerCommand)
	router.GET("/api/client/log", getClientLog)

	// Use port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}

func setMinerConfig(c *gin.Context, dbConfig DatabaseConfig, minerConns *sync.Map) {
	var config MinerConfig
	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Store config in database (PostgreSQL)
	err := storeMinerConfig(config, dbConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store config in database: " + err.Error()})
		return
	}

	// Generate XMRig config file
	configFile, err := generateXMRigConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Miner config updated successfully", "configFile": configFile})
}

func sendMinerCommand(c *gin.Context) {
    var command struct {
        MinerID string `json:"minerID"`
        Command string `json:"command"`
        Payload string `json:"payload"`
    }

    if err := c.BindJSON(&command); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get Connection for this miner Id
    connRaw, ok := minerConns.Load(command.MinerID)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "unable to find miner with this MinerId"})
        return
    }
    conn, _ := connRaw.(net.Conn)

    // Create ControlCommand Protobuf Message
    controlCommand := &xmbig.ControlCommand{
        Command: command.Command,
        Payload: command.Payload,
    }

    err := sendCommand(conn, controlCommand)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send command via TCP: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Command send successfully"})
}

func getClientLog(c *gin.Context) {
    clientId := c.Query("clientId")

    if clientId == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ClientId is required"})
        return
    }

    logContent, err := getMinerLogContent(clientId)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch the logfile"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"client_log": logContent})

}

func getMinerLogContent(clientId string) (string, error) {
    path := "logs/" + clientId + ".log"
    content, err := ioutil.ReadFile(path)
    if err != nil {
        return "", err
    }
    return string(content), nil
}

func storeMinerConfig(config MinerConfig, dbConfig DatabaseConfig) error {
    connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
    db, err := sql.Open("mysql", connStr)
    if err != nil {
        return err
    }
    defer db.Close()

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS miner_configs (
            miner_id VARCHAR(255) PRIMARY KEY,
            poolAddress VARCHAR(255) NOT NULL,
            username VARCHAR(255) NOT NULL,
            password VARCHAR(255),
            algorithm VARCHAR(255),
            autoSwitch BOOLEAN,
            torEnabled BOOLEAN,
            extraParams VARCHAR(255)
        )
    `)
    if err != nil {
        return err
    }

    // Prepare the insert statement
    stmt, err := db.Prepare(`
        INSERT INTO miner_configs (miner_id, poolAddress, username, password, algorithm, autoSwitch, torEnabled, extraParams)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE 
            poolAddress = ?, 
            username = ?, 
            password = ?, 
            algorithm = ?, 
            autoSwitch = ?, 
            torEnabled = ?, 
            extraParams = ?
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute the insert statement
    _, err = stmt.Exec(
        config.MinerID,
        config.PoolAddress,
        config.Username,
        config.Password,
        config.Algorithm,
        config.AutoSwitch,
        config.TorEnabled,
        config.ExtraParams,
        config.PoolAddress,
        config.Username,
        config.Password,
        config.Algorithm,
        config.AutoSwitch,
        config.TorEnabled,
        config.ExtraParams,
    )

    return err
}

func generateXMRigConfig(config MinerConfig) (string, error) {
	// Load a default config
	cfg, err := ini.Load("default_config.ini")
	if (err != nil) {
		fmt.Printf("Fail to read file: %v", err)
		return "", err
	}
	// General Configuration
	cfg.Section("").Key("url").SetValue(config.PoolAddress)
	cfg.Section("").Key("user").SetValue(config.Username)
	cfg.Section("").Key("pass").SetValue(config.Password)
	cfg.Section("").Key("algo").SetValue(config.Algorithm)

	// Store the configuration in the file
	err = cfg.SaveTo("config.ini")
	if err != nil {
		fmt.Printf("Fail to create file: %v", err)
		return "", err
	}
	return "config.ini", nil
}

// loadConfig reads configuration from file
func loadConfig(file string, cfg *Config) error {
    iniCfg, err := ini.Load(file)
    if err != nil {
        return err
    }

    err = iniCfg.MapTo(cfg)
    if err != nil {
        return err
    }

    return nil
}

// connectToDatabase establishes a connection to the database
func connectToDatabase(dbConfig DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

//TCP Code
func tcpServerConnector() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	fmt.Println("Listening on :9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting:", err)
			continue
		}

		// Lese die MinerID am Anfang der Verbindung.
		var minerIdBuffer [256]byte
		n, err := conn.Read(minerIdBuffer[:])
		if err != nil {
			log.Printf("Failed to read miner ID: %v", err)
			conn.Close()
			continue
		}

		minerId := string(minerIdBuffer[:n])
		minerId = strings.TrimSpace(minerId) // Entferne Leerzeichen

		minerConns.Store(minerId, conn)
		go handleConnection(conn, minerConns)

	}
}

func handleConnection(conn net.Conn, minerConns *sync.Map) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	for {
		// 1. Lese die Länge der Nachricht (als uint32)
		var length uint32
		err := binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				return
			}
			fmt.Println("Error reading message length:", err)
			return
		}

		// 2. Lese die Nachricht selbst
		buffer := make([]byte, length)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		if n != int(length) {
			fmt.Printf("Incomplete message read: expected %d, got %d\n", length, n)
			return
		}

		// 3. Deserialisiere die Protobuf-Nachricht
		var minerStatus xmbig.MinerStatus
		err = proto.Unmarshal(buffer, &minerStatus)
		if err != nil {
			fmt.Println("Error unmarshaling message:", err)
			return
		}

		// 4. Verarbeite die Nachricht
		fmt.Printf("Received MinerStatus: %+v from %s\n", minerStatus, conn.RemoteAddr())
		// Hier kannst du die Miner-Statusdaten in deiner Datenbank speichern oder andere Aktionen ausführen
	}
}

// Hilfsfunktion zum Senden von Befehlen
func sendCommand(conn net.Conn, cmd *xmbig.ControlCommand) error {
	data, err := proto.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	length := uint32(len(data))
	err = binary.Write(conn, binary.BigEndian, length)
	if err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	return nil
}