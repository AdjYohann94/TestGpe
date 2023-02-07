package domain

import "gpe_project/internal/app/adapter/postgresql/model"

type Quiz struct {
	*model.Quiz
	Responded bool `json:"responded"`
}
