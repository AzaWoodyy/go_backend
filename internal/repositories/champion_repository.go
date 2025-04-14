package repositories

import (
	"log"

	"github.com/AzaWoodyy/go_backend/internal/models"
	"gorm.io/gorm"
)

type ChampionRepository struct {
	db *gorm.DB
}

func NewChampionRepository(db *gorm.DB) *ChampionRepository {
	return &ChampionRepository{db: db}
}

func (r *ChampionRepository) SaveChampions(champions []models.Champion) error {
	for _, champion := range champions {
		log.Printf("Processing champion: %s", champion.RiotID)

		// First, save or update the champion without any associations
		var savedChampion models.Champion
		if err := r.db.Where(models.Champion{RiotID: champion.RiotID}).FirstOrCreate(&savedChampion, models.Champion{
			RiotID: champion.RiotID,
			Key:    champion.Key,
			Name:   champion.Name,
			Title:  champion.Title,
			Blurb:  champion.Blurb,
		}).Error; err != nil {
			log.Printf("Error saving champion %s: %v", champion.RiotID, err)
			continue
		}
		log.Printf("Successfully saved champion: %s", champion.RiotID)

		// Handle tags independently
		for _, tag := range champion.Tags {
			log.Printf("Processing tag: %s for champion: %s", tag.Key, champion.RiotID)

			// Try to create or find the tag
			var existingTag models.Tag
			if err := r.db.Where(models.Tag{Key: tag.Key}).FirstOrCreate(&existingTag).Error; err != nil {
				log.Printf("Warning: Could not create/find tag %s: %v", tag.Key, err)
				continue
			}

			// Try to create the association
			if err := r.db.Model(&savedChampion).Association("Tags").Append(&existingTag); err != nil {
				log.Printf("Warning: Could not associate tag %s with champion %s: %v", tag.Key, champion.RiotID, err)
				continue
			}
			log.Printf("Successfully associated tag %s with champion %s", tag.Key, champion.RiotID)
		}

		// Handle versions independently
		for _, version := range champion.Versions {
			log.Printf("Processing version: %s for champion: %s", version.Key, champion.RiotID)

			// First, ensure the version exists
			var existingVersion models.Version
			result := r.db.Where(models.Version{Key: version.Key}).First(&existingVersion)

			if result.Error != nil {
				if result.Error == gorm.ErrRecordNotFound {
					// Version doesn't exist, create it
					if err := r.db.Create(&version).Error; err != nil {
						log.Printf("Warning: Could not create version %s: %v", version.Key, err)
						continue
					}
					existingVersion = version
				} else {
					log.Printf("Warning: Could not find version %s: %v", version.Key, result.Error)
					continue
				}
			}

			// Now associate the version with the champion
			if err := r.db.Model(&savedChampion).Association("Versions").Append(&existingVersion); err != nil {
				log.Printf("Warning: Could not associate version %s with champion %s: %v", version.Key, champion.RiotID, err)
				continue
			}
			log.Printf("Successfully associated version %s with champion %s", version.Key, champion.RiotID)
		}
	}
	return nil
}

func (r *ChampionRepository) GetChampions() ([]models.Champion, error) {
	var champions []models.Champion
	err := r.db.Preload("Tags").Preload("Versions").Find(&champions).Error
	return champions, err
}
