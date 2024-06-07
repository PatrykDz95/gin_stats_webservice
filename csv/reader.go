package csv

import (
	"encoding/csv"
	"gin/database"
	"gin/entities"
	"log"
	"os"
	"strconv"
)

func WriteCsvIntoDB() {
	// Open the CSV file
	file, err := os.Open("player_stats.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Skip the header
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var playerStatsSlice []entities.PlayerStats

	// Insert all records into the slice
	for _, record := range records {
		age, err := strconv.Atoi(record[3])
		if err != nil {
			log.Fatal(err)
		}
		playerStats := entities.PlayerStats{
			Name:    record[1],
			Surname: record[2],
			Age:     age,
		}
		playerStatsSlice = append(playerStatsSlice, playerStats)
	}

	// Insert the slice into the database in a single batch
	created := database.DB.CreateInBatches(&playerStatsSlice, len(playerStatsSlice))
	if created.Error != nil {
		log.Fatal(created.Error)
	}

}
