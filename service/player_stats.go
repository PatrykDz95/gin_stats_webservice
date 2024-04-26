package service

import (
	"gin/models"
)

func Get() ([]*models.PlayerStats, error) {
	var stats []*models.PlayerStats
	//fill stats with data
	stats = append(stats, &models.PlayerStats{
		Points:   100,
		Assists:  50,
		Rebounds: 30,
	})
	return stats, nil
}
