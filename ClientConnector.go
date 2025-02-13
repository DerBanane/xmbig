package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
    "context"
	"time"
	"github.com/derBanane/xmbig"
	"google.golang.org/protobuf/proto"
)

const (
	serverAddress = "localhost:9000"
)

// Helper function to send a MinerStatus message

func StartMiner(config MinerConfig) error {
   	// Generate XMRig config file
	configFile, err := generateXMRigConfig(config)
	if err != nil {
		return fmt.Errorf("failed to generate xmrig Config: %w", err)
	}

	// Build the command
	cmdStr := fmt.Sprintf("./xmrig --config=%s", configFile)

	// Start XMRig
	cmd := exec.Command("cmd", "/C", cmdStr) // Windows

	//For the cmd, you will have to enter the command needed

	log.Printf("starting miner with command: %s", cmdStr)

	err = cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	return nil
}

// generateXMRigConfig generates XMRig config file
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