package repository

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/tag"
)

type tagRepository struct {
}

func NewTagRepository() tag.Repository {
	return &tagRepository{}
}

func (tr *tagRepository) GetAllTags() ([]models.Tag, error) {

	tags := []models.Tag{
		{Name: "хочувБАР", TagId: 1},
		{Name: "хочувБАР",TagId: 2},
		{Name: "хочувКИНО",TagId: 3},
		{Name: "хочувТЕАТР",TagId: 4},
		{Name: "хочувКЛУБ",TagId: 5},
		{Name: "хочунаКОНЦЕРТ",TagId: 6},
		{Name: "хочуГУЛЯТЬ",TagId: 7},
		{Name: "хочунаКАТОК",TagId: 8},
		{Name: "хочунаВЫСТАВКУ",TagId: 9},
		{Name: "хочуСПАТЬ",TagId: 11},
		{Name: "хочунаСАЛЮТ",TagId: 12},
		{Name: "хочувСПОРТ",TagId: 13},
		{Name: "хочувМУЗЕЙ",TagId: 14},
		{Name: "хочунаЛЕКЦИЮ",TagId: 15},
		{Name: "хочуБОТАТЬ",TagId: 16},
		{Name: "хочувПАРК",TagId: 17},
		{Name: "хочувГОСТИ",TagId: 18},
	}

	return tags, nil
}
