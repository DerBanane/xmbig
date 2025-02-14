package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"time"

	"github.com/derbanane/xmbig/xmbig"
	"google.golang.org/protobuf/proto"
	"gopkg.in/ini.v1"
)

type MinerConfig struct {
	PoolAddress string
	Username    string
	Password    string
	Algorithm   string
}

func StartMiner(config MinerConfig) error {
	configFile, err := createXMRigConfigFile(config)
	if err != nil {
		return fmt.Errorf("failed to generate xmrig Config: %w", err)
	}

	// Build the command
	cmdStr := fmt.Sprintf("./xmrig --config=%s", configFile)

	// Start XMRig
	cmd := exec.Command("cmd", "/C", cmdStr)
	log.Printf("starting miner with command: %s", cmdStr)
	err = cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	return nil
}


func createXMRigConfigFile(config MinerConfig) (string, error) {
	// Load a default config
	cfg, err := ini.Load("default_config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return "", err
	}
	// General Configuration
	cfg.Section("pool").Key("url").SetValue(config.PoolAddress)
	cfg.Section("pool").Key("user").SetValue(config.Username)
	cfg.Section("pool").Key("pass").SetValue(config.Password)
	cfg.Section("pool").Key("algo").SetValue(config.Algorithm)

	// Store the configuration in the file
	err = cfg.SaveTo("config.ini")
	if err != nil {
		fmt.Printf("Fail to create file: %v", err)
		return "", err
	}
	return "config.ini", nil
}

func RunClient(serverAddress string, minerId string) {
	
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to server %s", serverAddress)

	// 1. Sende die Miner-ID
	_, err = conn.Write([]byte(minerId))
	if err != nil {
		log.Fatalf("Failed to send miner ID: %v", err)
	}
	// Run the function
	go readMessages(conn)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Erstelle eine Beispiel MinerStatus-Nachricht
		minerStatus := &xmbig.MinerStatus{
			ClientId:  minerId,
			Status:    "Mining",
			Hashrate:  12345.67,
			Temperature: 60.5,
		}

		// Serialisiere die Nachricht
		data, err := proto.Marshal(minerStatus)
		if err != nil {
			log.Printf("Failed to marshal MinerStatus: %v", err)
			continue // Skip this iteration
		}
		err = sendProtobufMessage(conn, data)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			continue
		}
		log.Println("MinerStatus sent successfully")
	
	}
}

func sendProtobufMessage(conn net.Conn, data []byte) error {
	length := uint32(len(data))
	err := binary.Write(conn, binary.BigEndian, length)
	if err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func readMessages(conn net.Conn) {
	for {
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

		var cmd xmbig.ControlCommand
		err = proto.Unmarshal(buffer, &cmd)
		if err != nil {
			fmt.Println("Error during unmarshal", err)
			return
		}

		fmt.Println("Received Command", cmd.Command, cmd.Payload)
	}
}