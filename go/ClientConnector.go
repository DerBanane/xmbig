package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
	"github.com/DerBanane/xmbig"
	"google.golang.org/protobuf/proto"
)

type ClientConnector struct {
	ServerAddress string
	conn          net.Conn
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

func (cc *ClientConnector) SendMinerStats(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second) // Send data every 10 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Create a sample MinerStatus message
			minerStatus := &xmbig.MinerStatus{
				ClientId:    "miner-123", // Replace with the actual MinerID
				Status:      "Mining",
				Hashrate:    12345.67,
				Temperature: 60.5,
			}

			// Serialize the message
			data, err := proto.Marshal(minerStatus)
			if err != nil {
				log.Printf("Failed to marshal MinerStatus: %v", err)
				continue
			}

			// Prepend the message with its length
			length := uint32(len(data))
			err = binary.Write(cc.conn, binary.BigEndian, length)
			if err != nil {
				log.Printf("Failed to write message length: %v", err)
				continue
			}

			// Send the message
			_, err = cc.conn.Write(data)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
				continue
			}

			log.Println("MinerStatus sent successfully")

		case <-ctx.Done():
			log.Println("Stats sender stopped")
			return
		}
	}
}

func tcpClientConnector() {
	clientConnector := &ClientConnector{
		ServerAddress: "0.0.0.0.8080",
	}
	ctx := context.Background()
	go clientConnector.SendMinerStats(ctx)
}
