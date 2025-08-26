package services

import (
	"log"
	"time"

	"github.com/NOTMKW/DLLBEL/internal/models"
)

type EventService struct {
	ruleService *RuleService
	userService *UserService
	dllService  *DLLService
	wsService   *WebSocketService
	eventChan   chan *models.MT5Event
	done        chan bool
}

func NewEventService(ruleService *RuleService, userService *UserService, dllService *DLLService, wsService *WebSocketService, bufferSize int) *EventService {
	return &EventService{
		ruleService: ruleService,
		userService: userService,
		dllService:  dllService,
		wsService:   wsService,
		eventChan:   make(chan *models.MT5Event, bufferSize),
		done:        make(chan bool),
	}
}

func (s *EventService) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func() {
			for {
				select {
				case event := <-s.eventChan:
					s.processEvent(event)
				case <-s.done:
					return
				}
			}
		}()
	}
	log.Printf("Event service started with %d workers", workers)
}

func (s *EventService) Stop() {
	close(s.done)
	log.Println("Event service stopped")
}

func (s *EventService) GetEventChannel() chan *models.MT5Event {
	return s.eventChan
}

func (s *EventService) processEvent(event *models.MT5Event) {
	userState := s.userService.GetUserState(event.UserId)
	if userState == nil {
		userState = s.userService.CreateUserState(event.UserId)
	}
	s.userService.UpdateUserStateWithEvent(userState, event)

	rules, err := s.ruleService.GetAllRules()
	if err != nil {
		log.Printf("Failed to get rules: %v", err)
		return
	}

	for _, rule := range rules {
		if rule.Enabled && s.ruleService.EvaluateRule(rule, event, userState) {
			for _, action := range rule.Actions {
				enforcement := &models.EnforcementMessage{
					UserId:    event.UserId,
					Action:    action.Type,
					Reason:    "Rule violation: " + rule.Name,
					Severity:  action.Severity,
					Timestamp: time.Now().Unix(),
				}
				s.dllService.SendEnforcement(enforcement)
				s.wsService.SendEnforcement(enforcement)
				log.Printf("Enforcement action '%s' triggered for user %s due to rule '%s'", action.Type, event.UserId, rule.Name)
			}
		}
	}
}
