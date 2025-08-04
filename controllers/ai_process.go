package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/services"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"NotificationManagement/utils/errutil"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AIRequestControllerImpl struct {
	AIModelService domain.AIModelService
	D              services.DeepseekProcessServiceImpl
	G              services.GeminiProcessServiceImpl
}

func NewAIRequestController(A domain.AIModelService, D services.DeepseekProcessServiceImpl, G services.GeminiProcessServiceImpl) domain.AIRequestController {
	return &AIRequestControllerImpl{AIModelService: A, D: D, G: G}
}

func (a *AIRequestControllerImpl) MakeAIRequestHandler(c echo.Context) error {
	var req types.MakeAIRequestPayload
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}
	service, err := a.GetServiceManagerById(req.ModelID)
	if err != nil {
		return err
	}
	aiResponse, err := service.MakeAIRequest(&req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, aiResponse)
}

func (a *AIRequestControllerImpl) CreateAIModel(c echo.Context) error {

	var req types.AIModelRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}
	model, err := req.ToModel()
	if err != nil {
		return err
	}
	service, err := a.GetServiceManagerByType(model.GetType())
	if err != nil {
		return err
	}
	err = service.CreateModel(model)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, model)

}

func (a *AIRequestControllerImpl) GetAIModelByID(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}
	service, err := a.GetServiceManagerById(id)
	if err != nil {
		return err
	}
	modelById, err := service.GetModelById(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, modelById)
}

func (a *AIRequestControllerImpl) GetAllAIModels(c echo.Context) error {
	limit, offset := utils.ParseLimitAndOffset(c)

	deepseekService, err := a.GetServiceManagerByType("deepseek")
	if err != nil {
		return err
	}
	deepseekModels, err := deepseekService.GetAllModels(limit, offset)
	if err != nil {
		return err
	}

	geminiService, err := a.GetServiceManagerByType("gemini")
	if err != nil {
		return err
	}
	geminiModels, err := geminiService.GetAllModels(limit, offset)
	if err != nil {
		return err
	}

	var responses []interface{}
	for _, model := range deepseekModels.([]models.DeepseekModel) {
		//deepseekModel := any(model).(*models.DeepseekModel)
		responses = append(responses, types.FromDeepseekModel(&model))
	}
	for _, model := range geminiModels.([]models.GeminiModel) {
		//geminiModel := any(model).(*models.GeminiModel)
		responses = append(responses, types.FromGeminiModel(&model))
	}

	return c.JSON(http.StatusOK, responses)
}

func (a *AIRequestControllerImpl) UpdateAIModel(c echo.Context) error {
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
	s, err := a.GetServiceManagerByType(req.Type)
	if err != nil {
		return err
	}
	updateModel, err := s.UpdateModel(id, model)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updateModel)
}

func (a *AIRequestControllerImpl) DeleteAIModel(c echo.Context) error {
	id, err := utils.ParseIDFromContext(c)
	if err != nil {
		return err
	}

	return a.AIModelService.DeleteModel(id)
}

func (a *AIRequestControllerImpl) GetServiceManagerById(id uint) (domain.AIProcessService[domain.AIService[any], any], error) {
	model, err := a.AIModelService.GetModelById(id)
	if err != nil {
		return nil, err
	}
	return a.GetServiceManagerByType(model.Type)
}

func (a *AIRequestControllerImpl) GetServiceManagerByType(modelType string) (domain.AIProcessService[domain.AIService[any], any], error) {
	switch modelType {
	case "deepseek":
		return a.D, nil
	case "gemini":
		return a.G, nil
	}
	return nil, errutil.NewAppError(errutil.ErrFeatureNotAvailable, errutil.ErrInvalidFeature)
}
