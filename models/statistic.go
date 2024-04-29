package models

type PlayerStats struct {
	ID        int    `gorm:"primary_key" json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Age       int    `json:"age"`
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
