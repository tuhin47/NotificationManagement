package models

import (
	"gorm.io/gorm"
)

type RequestAIModel struct {
	gorm.Model
	RequestID uint     `gorm:"index:idx_request_ai_model,unique"`
	IsActive  bool     `gorm:"default:true"`
	AiModelID uint     `gorm:"index:idx_request_ai_model,unique"`
	AiModel   *AIModel `gorm:"foreignKey:AiModelID"`
	//Parameters   JSON   `gorm:"type:jsonb"`
}

func (u *RequestAIModel) UpdateFromModel(source ModelInterface) {
	if src, ok := source.(*RequestAIModel); ok {
		copyFields(u, src)
	}
}
