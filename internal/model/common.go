package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

const GormDataTypeJSON = "json"

type Trackers struct {
	CreatedAt time.Time      `json:"created_at" gorm:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type JSON map[string]interface{}

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return ErrorCannotUnmarshalJSONB
	}

	return json.Unmarshal(bytes, j)
}

func (j JSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j JSON) GormDataType() string {
	return GormDataTypeJSON
}
