package services

import (
	"github.com/NOTMKW/DLLBEL/internal/dto"
	"github.com/NOTMKW/DLLBEL/internal/models"
	"github.com/NOTMKW/DLLBEL/internal/repository"
	"fmt"
	"strconv"
	"time"
)

type RuleService struct {
	repo *repository.RedisRepository
}

func NewRuleService(repo *repository.RedisRepository) *RuleService {
	return &RuleService{repo: repo}
}

func (s *RuleService) CreateRule(req *dto.CreateRuleRequest) (*models.Rule, error) {
	rule := &models.Rule{
		ID:         fmt.Sprintf("rule-%d", time.Now().UnixNano()),
		Name:       req.Name,
		Conditions: req.Conditions,
		Actions:    req.Actions,
		Enabled:    req.Enabled,
		Priority:   req.Priority,
		CreatedAt:  time.Now().UnixNano(),
		UpdatedAt:  time.Now().UnixNano(),
	}

	if err := s.repo.SaveRule(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *RuleService) UpdateRule(id string, req *dto.UpdateRuleRequest) (*models.Rule, error) {
	rule, err := s.repo.GetRule(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		rule.Name = req.Name
	}
	if req.Conditions != nil {
		rule.Conditions = req.Conditions
	}
	if req.Actions != nil {
		rule.Actions = req.Actions
	}
	if req.Enabled != nil {
		rule.Enabled = *req.Enabled
	}
	if req.Priority != nil {
		rule.Priority = *req.Priority
	}
	rule.UpdatedAt = time.Now().UnixNano()

	if err := s.repo.SaveRule(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *RuleService) DeleteRule(id string) error {
	return s.repo.DeleteRule(id)
}

func (s *RuleService) GetAllRules() ([]*models.Rule, error) {
	return s.repo.GetAllRules()
}

func (s *RuleService) EvaluateRule (rule *models.Rule, event *models.MT5Event, state *models.UserState) bool {
	for field, value := range rule.Conditions {
		switch field {
		case "max_volume":
			if maxVol, err := strconv.ParseFloat(value, 64); err == nil {
				if event.Volume > maxVol {
				return false
			}
		}
		case "max_positions":
			if maxPos, err := strconv.Atoi(value); err == nil {
				if state.OpenPositions > maxPos {
					return true
				}
			}
		case "max_day_volume":
			if maxDayVol, err := strconv.ParseFloat(value, 64); err == nil {
				if state.DayVolume > maxDayVol {
					return true
				}
			}
		case "symbol_restricted":
			if event.Symbol == value {
				return true
			}
		}
	}

	return false
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}