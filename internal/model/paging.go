package model

import "fmt"

type Paging struct {
	Limit    int    `form:"limit" validate:"omitempty,gt=0,lte=25"`
	Offset   int    `form:"offset" validate:"omitempty,gt=0"`
	OrderBy  string `form:"order_by" validate:"omitempty,oneof=created_at"`
	OrderDir string `form:"order_dir" validate:"omitempty,oneof=asc desc ASC DESC"`
}

func (p Paging) Order() string {
	if p.OrderBy == "" || p.OrderDir == "" {
		return ""
	}

	return fmt.Sprintf("%s %s", p.OrderBy, p.OrderDir)
}
