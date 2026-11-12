package main

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type GardenDome struct {
	ID   uint        `gorm:"primaryKey"`
	Name string      `gorm:"not null;uniqueIndex"`
	Beds []GardenBed `gorm:"foreignKey:DomeID;references:ID;constraint:OnDelete:RESTRICT"`
}
type GardenBed struct {
	ID     uint   `gorm:"primaryKey"`
	DomeID uint   `gorm:"not null;uniqueIndex:bed_label"`
	Label  string `gorm:"not null;uniqueIndex:bed_label"`
	Active bool
}

func InitGarden(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).AutoMigrate(&GardenDome{}, &GardenBed{})
}
func CreateGardenDome(ctx context.Context, db *gorm.DB, name string, labels []string) (GardenDome, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return GardenDome{}, errors.New("blank dome name")
	}
	out := GardenDome{Name: name, Beds: []GardenBed{}}
	for _, label := range labels {
		label = strings.TrimSpace(label)
		if label == "" {
			return GardenDome{}, errors.New("blank bed label")
		}
		out.Beds = append(out.Beds, GardenBed{Label: label, Active: true})
	}
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return tx.Create(&out).Error })
	if err != nil {
		return GardenDome{}, err
	}
	return out, nil
}
func ListGarden(ctx context.Context, db *gorm.DB) ([]GardenDome, error) {
	out := []GardenDome{}
	err := db.WithContext(ctx).Preload("Beds", func(q *gorm.DB) *gorm.DB { return q.Order("id ASC") }).Order("id ASC").Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}
func SetGardenBed(ctx context.Context, db *gorm.DB, id uint, active bool) error {
	result := db.WithContext(ctx).Model(&GardenBed{}).Where("id = ?", id).Update("active", active)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func GardenQueryCount(domes int, preload bool) int {
	if domes == 0 {
		return 1
	}
	if preload {
		return 2
	}
	return domes + 1
}
func main() {}
