package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"context"
)

type LLMServiceImpl struct {
	Repo domain.LLMRepository
}

func NewLLMService(repo domain.LLMRepository) domain.LLMService {
	return &LLMServiceImpl{Repo: repo}
}

func (s *LLMServiceImpl) CreateLLM(llm *models.UserLLM) error {
	return s.Repo.Create(context.Background(), llm)
}

func (s *LLMServiceImpl) GetLLMByID(id uint) (*models.UserLLM, error) {
	llm, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return llm, nil
}

func (s *LLMServiceImpl) GetAllLLMs(limit, offset int) ([]models.UserLLM, error) {
	llms, err := s.Repo.GetAll(context.Background(), limit, offset)
	if err != nil {
		return nil, err
	}
	return llms, nil
}

func (s *LLMServiceImpl) UpdateLLM(id uint, llm *models.UserLLM) error {
	// First check if the record exists
	existing, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	// Update the existing record with new data
	existing.RequestID = llm.RequestID
	existing.ModelName = llm.ModelName
	existing.Type = llm.Type
	existing.IsActive = llm.IsActive

	// Save the updated record
	err = s.Repo.Update(context.Background(), existing)
	if err != nil {
		return err
	}

	return nil
}

func (s *LLMServiceImpl) DeleteLLM(id uint) error {
	// First check if the record exists
	_, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	err = s.Repo.Delete(context.Background(), id)
	if err != nil {
		return err
	}

	return nil
}
