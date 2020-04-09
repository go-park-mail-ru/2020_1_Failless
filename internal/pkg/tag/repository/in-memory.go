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
		{"хочувБАР", 1},
		{"хочувБАР", 2},
		{"хочувКИНО", 3},
		{"хочувТЕАТР", 4},
		{"хочувКЛУБ", 5},
		{"хочунаКОНЦЕРТ", 6},
		{"хочуГУЛЯТЬ", 7},
		{"хочунаКАТОК", 8},
		{"хочунаВЫСТАВКУ", 9},
		{"хочуСПАТЬ", 11},
		{"хочунаСАЛЮТ", 12},
		{"хочувСПОРТ", 13},
		{"хочувМУЗЕЙ", 14},
		{"хочунаЛЕКЦИЮ", 15},
		{"хочуБОТАТЬ", 16},
		{"хочувПАРК", 17},
		{"хочувГОСТИ", 18},
	}

	return tags, nil
}
