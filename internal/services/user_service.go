package services

import (
	"sync"
	"time"

	"github.com/NOTMKW/DLLBEL/internal/dto"
	"github.com/NOTMKW/DLLBEL/internal/models"
	"github.com/NOTMKW/DLLBEL/internal/repository"
)

type UserService struct {
	repo   *repository.RedisRepository
	states map[string]*models.UserState
	mu     sync.RWMutex
}

func NewUserService(repo *repository.RedisRepository) *UserService {
	return &UserService{
		repo:   repo,
		states: make(map[string]*models.UserState),
	}
}

func (s *UserService) GetUserState(userID string) *models.UserState {
	s.mu.RLock()
	state, exists := s.states[userID]
	s.mu.RUnlock()
	if exists {
		return state
	}

	if redisState, err := s.repo.GetUserState(userID); err == nil {
		s.mu.Lock()
		s.states[userID] = redisState
		s.mu.Unlock()
		return redisState
	}

	return nil
}

func (s *UserService) CreateUserState(userID string) *models.UserState {
	newState := &models.UserState{
		UserID:       userID,
		CustomData:   make(map[string]string),
		LastActivity: time.Now().Unix(),
	}
	s.mu.Lock()
	s.states[userID] = newState
	s.mu.Unlock()

	go s.repo.SaveUserState(newState)

	return newState
}

func (s *UserService) UpdateUserState(userID string, req *dto.UpdateUserStateRequest) *models.UserState {
	state := s.GetUserState(userID)
	state.Mu.Lock()
	defer state.Mu.Unlock()
	if req.Balance != nil {
		state.Balance = *req.Balance
	}
	if req.Equity != nil {
		state.Equity = *req.Equity
	}

	if req.OpenPositions != nil {
		state.OpenPositions = *req.OpenPositions
	}

	if req.DayVolume != nil {
		state.DayVolume = *req.DayVolume
	}

	if req.RiskLevel != nil {
		state.RiskLevel = *req.RiskLevel
	}

	if req.ViolationCount != nil {
		state.ViolationCount = *req.ViolationCount
	}

	if req.CustomData != nil {
		for k, v := range req.CustomData {
			state.CustomData[k] = v
		}
	}

	state.LastActivity = time.Now().Unix()

	go s.repo.SaveUserState(state)

	return state
}

func (s *UserService) UpdateUserStateWithEvent(state *models.UserState, event *models.MT5Event) {
	state.Mu.Lock()
	defer state.Mu.Unlock()

	state.LastActivity = time.Now().Unix()

	switch event.EventType {
	case "ORDER_OPEN":
		state.DayVolume += event.Volume
		state.OpenPositions += 1
	case "ORDER_CLOSE":
		state.OpenPositions -= 1
	case "BALANCE_UPDATE":
		state.Balance = event.Price
	case "EQUITY_UPDATE":
		state.Equity = event.Price
	}

	state.LastActivity = time.Now().Unix()

	go s.repo.SaveUserState(state)
}

func (s *UserService) GetUserCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.states)
}

func (s *UserService) SyncAllStates() {
	s.mu.RLock()
	states := make([]*models.UserState, 0, len(s.states))
	for _, state := range s.states {
		states = append(states, state)
	}
	s.mu.RUnlock()

	for _, state := range states {
		s.repo.SaveUserState(state)
	}
}
