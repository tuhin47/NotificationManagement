package models

import (
	"google.golang.org/genai"
	"gorm.io/gorm"
)

type CurlRequest struct {
	gorm.Model
	URL              string              `gorm:"type:text" json:"url"`
	Method           string              `gorm:"type:varchar(10)" json:"method"`
	Headers          string              `gorm:"type:text" json:"headers"`
	Body             string              `gorm:"type:text" json:"body"`
	RawCurl          string              `gorm:"type:text" json:"rawCurl"`
	Reminders        *[]Reminder         `gorm:"foreignKey:RequestID"`
	Models           *[]RequestAIModel   `gorm:"foreignKey:RequestID"`
	AdditionalFields *[]AdditionalFields `gorm:"foreignKey:RequestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"additional_fields"`
}

func (c *CurlRequest) UpdateFromModel(source ModelInterface) {
	if src, ok := source.(*CurlRequest); ok {
		copyFields(c, src)
	}
}

func (c *CurlRequest) GetGenaiSchemaProperties() map[string]*genai.Schema {
	properties := make(map[string]*genai.Schema)
	for _, field := range *c.AdditionalFields {
		var schemaType genai.Type
		switch field.Type {
		case "number":
			schemaType = genai.TypeNumber
		case "boolean":
			schemaType = genai.TypeBoolean
		default:
			schemaType = genai.TypeString
		}
		properties[field.PropertyName] = &genai.Schema{
			Type:        schemaType,
			Description: field.Description,
		}
	}
	return properties
}

type AdditionalFields struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	PropertyName string `gorm:"type:varchar(100)" json:"property_name"`
	Type         string `gorm:"type:varchar(10)" json:"type"` // allowed: number, boolean, text
	Description  string `gorm:"type:text" json:"description,omitempty"`
	RequestID    uint   `json:"request_id"`
}
