package models

import (
	"github.com/go-playground/validator/v10"
)

type PlayerStats struct {
	ID        int    `gorm:"primary_key" json:"id"`
	Name      string `json:"name" validate:"required,min=3,max=20"`
	Surname   string `json:"surname" validate:"required,min=3,max=20"`
	Age       int    `json:"age" validate:"required,min=18,max=50"`
	Position  string `json:"position"`
	Team      string `json:"team"`
	Height    int    `json:"height"`
	Weight    int    `json:"weight"`
	Points    int    `json:"points"`
	Assists   int    `json:"assists"`
	Rebounds  int    `json:"rebounds"`
	Steals    int    `json:"steals"`
	Blocks    int    `json:"blocks"`
	Turnovers int    `json:"turnovers"`
	Minutes   int    `json:"minutes"`
	CreatedOn string `json:"created_on"`
}

func (p *PlayerStats) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}
