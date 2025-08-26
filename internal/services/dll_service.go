package services

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"

	"github.com/NOTMKW/DLLBEL/internal/dto"
	"github.com/NOTMKW/DLLBEL/internal/models"
)

type DLLService struct {
	connections map[string]*models.DLLConnection
	mu          sync.RWMutex
	eventChan   chan *models.MT5Event
}

func NewDLLService(eventChan chan *models.MT5Event) *DLLService {
	return &DLLService{
		connections: make(map[string]*models.DLLConnection),
		eventChan:   eventChan,
	}
}

func (s *DLLService) StartListener(dllID string) error {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	log.Printf("DLL %s listening on port %d", dllID, port)

	go func() {
		defer listener.Close()
		
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			dllConn := &models.DLLConnection{
				ID:          dllID,
				Conn:        conn,
				IsActive:    true,
				LastPing:    time.Now().Unix(),
				EventChan:   make(chan *models.MT5Event, 1000),
				EnforceChan: make(chan *models.EnforcementMessage, 1000),
			}

			s.mu.Lock()
			s.connections[dllID] = dllConn
			s.mu.Unlock()

			go s.handleConnection(dllConn)
			go s.enforceWriter(dllConn)
		}
	}()

	return nil
}

func (s *DLLService) handleConnection(dllConn *models.DLLConnection) {
	defer func() {
		dllConn.Conn.Close()
		s.mu.Lock()
		delete(s.connections, dllConn.ID)
		s.mu.Unlock()
	}()

	buffer := make([]byte, 8192)
	
	for {
		n, err := dllConn.Conn.Read(buffer)
		if err != nil {
			log.Printf("DLL connection %s read error: %v", dllConn.ID, err)
			break
		}

		events := s.parseBinaryProtocol(buffer[:n])
		for _, event := range events {
			select {
			case s.eventChan <- event:
			default:
				log.Printf("Event buffer full, dropping event from DLL %s", dllConn.ID)
			}
		}

		dllConn.Mu.Lock()
		dllConn.LastPing = time.Now().Unix()
		dllConn.Mu.Unlock()
	}
}

func (s *DLLService) parseBinaryProtocol(data []byte) []*models.MT5Event {
	events := []*models.MT5Event{}
	reader := bytes.NewReader(data)

	for reader.Len() > 0 {
		var msgLen uint32
		if err := binary.Read(reader, binary.LittleEndian, &msgLen); err != nil {
			break
		}

		if msgLen > uint32(reader.Len()) {
			break
		}

		msgData := make([]byte, msgLen)
		if _, err := reader.Read(msgData); err != nil {
			break
		}

		event := &models.MT5Event{}
		if err := event.Deserialize(msgData); err != nil {
			log.Printf("Failed to deserialize MT5Event: %v", err)
			continue
		}
		events = append(events, event)
	}

	return events
}

func (s *DLLService) SendEnforcement(enforcement *models.EnforcementMessage) {
	s.mu.RLock()
	var targetConn *models.DLLConnection
	for _, conn := range s.connections {
		if conn.IsActive {
			targetConn = conn
			break
		}
	}
	s.mu.RUnlock()

	if targetConn != nil {
		select {
		case targetConn.EnforceChan <- enforcement:
		default:
			log.Printf("Enforcement channel full for DLL %s", targetConn.ID)
		}
	}
}

func (s *DLLService) enforceWriter(dllConn *models.DLLConnection) {
	for enforcement := range dllConn.EnforceChan {
		data, err := enforcement.Serialize()
		if err != nil {
			continue
		}

		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, uint32(len(data)))
		buf.Write(data)

		dllConn.Conn.Write(buf.Bytes())
	}
}

func (s *DLLService) GetConnections() []*dto.ConnectionInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conns := make([]*dto.ConnectionInfo, 0, len(s.connections))
	for _, conn := range s.connections {
		conn.Mu.RLock()
		conns = append(conns, &dto.ConnectionInfo{
			ID:       conn.ID,
			Active:   conn.IsActive,
			LastPing: conn.LastPing,
		})
		conn.Mu.RUnlock()
	}

	return conns
}
func (s *DLLService) GetActiveConnectionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, conn := range s.connections {
		if conn.IsActive {
			count++
		}
	}
	return count
}

func (s *DLLService) CheckHealth() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Unix()
	for id, conn := range s.connections {
		conn.Mu.RLock()
		lastPing := conn.LastPing
		conn.Mu.RUnlock()

		if now-lastPing > 60 {
			conn.IsActive = false
			log.Printf("DLL connection %s marked as inactive", id)
		}
	}
}
