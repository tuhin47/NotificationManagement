package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"NotificationManagement/utils/errutil"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AIModelControllerImpl struct {
	AiManager      domain.AIServiceManager
	AIModelService domain.AIModelService
}

func NewAIModelController(manager domain.AIServiceManager, AIModelService domain.AIModelService) domain.AIModelController {
	return &AIModelControllerImpl{AiManager: manager, AIModelService: AIModelService}
}

func (dc *AIModelControllerImpl) CreateAIModel(c echo.Context) error {
	var req types.AIModelRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	model, err := req.ToModel()
	if err != nil {
		return err
	}

	service, err := dc.AiManager.GetService(req.Type)
	if err != nil {
		return err
	}

	switch s := service.(type) {
	case domain.DeepseekService:
		deepseekModel, ok := model.(*models.DeepseekModel)
		if !ok {
			return errutil.NewAppError(errutil.ErrFailedToCastModel, errutil.ErrFailedToCastDeepseekModel)
		}
		err = s.CreateModel(deepseekModel)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, types.FromDeepseekModel(deepseekModel))
	case domain.GeminiService:
		geminiModel, ok := model.(*models.GeminiModel)
		if !ok {
			return errutil.NewAppError(errutil.ErrFailedToCastModel, errutil.ErrFailedToCastGeminiModel)
		}
		err = s.CreateModel(geminiModel)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, types.FromGeminiModel(geminiModel))
	default:
		return errutil.NewAppError(errutil.ErrUnsupportedAIModelType, errutil.ErrUnsupportedAIModelTypeMsg)
	}
}

func (dc *AIModelControllerImpl) GetAIModelByID(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}
	modelByID, err := dc.AIModelService.GetModelByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, modelByID)
}

func (dc *AIModelControllerImpl) GetAllAIModels(c echo.Context) error {
	limit, offset := utils.ParseLimitAndOffset(c)

	// TODO
	deepseekService, err := dc.AiManager.GetService("deepseek")
	if err != nil {
		return err
	}
	deepseekModels, err := deepseekService.(domain.DeepseekService).GetAllModels(limit, offset)
	if err != nil {
		return err
	}

	geminiService, err := dc.AiManager.GetService("gemini")
	if err != nil {
		return err
	}
	geminiModels, err := geminiService.(domain.GeminiService).GetAllModels(limit, offset)
	if err != nil {
		return err
	}

	var responses []interface{}
	for _, model := range deepseekModels {
		responses = append(responses, types.FromDeepseekModel(&model))
	}
	for _, model := range geminiModels {
		responses = append(responses, types.FromGeminiModel(&model))
	}

	return c.JSON(http.StatusOK, responses)
}

func (dc *AIModelControllerImpl) UpdateAIModel(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	var req types.AIModelRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	model, err := req.ToModel()
	if err != nil {
		return err
	}

	service, err := dc.AiManager.GetService(req.Type)
	if err != nil {
		return err
	}

	switch s := service.(type) {
	case domain.DeepseekService:
		deepseekModel, ok := model.(*models.DeepseekModel)
		if !ok {
			return errutil.NewAppError(errutil.ErrFailedToCastModel, errutil.ErrFailedToCastDeepseekModel)
		}
		err = s.UpdateModel(id, deepseekModel)
		if err != nil {
			return err
		}
		updatedModel, err := s.GetModelByID(id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, types.FromDeepseekModel(updatedModel))
	case domain.GeminiService:
		geminiModel, ok := model.(*models.GeminiModel)
		if !ok {
			return errutil.NewAppError(errutil.ErrFailedToCastModel, errutil.ErrFailedToCastGeminiModel)
		}
		err = s.UpdateModel(id, geminiModel)
		if err != nil {
			return err
		}
		updatedModel, err := s.GetModelByID(id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, types.FromGeminiModel(updatedModel))
	default:
		return errutil.NewAppError(errutil.ErrUnsupportedAIModelType, errutil.ErrUnsupportedAIModelTypeMsg)
	}
}

func (dc *AIModelControllerImpl) DeleteAIModel(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	return dc.AIModelService.DeleteModel(id)
}
