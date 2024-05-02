package models

import "time"

type Job struct {
	StartedAt       time.Time `reindex:"started_at" json:"started_at" validate:"required"`
	EndedAt         time.Time `reindex:"ended_at" json:"ended_at,omitempty" validate:"gtfield=StartedAt"`
	Name            string    `reindex:"name" json:"name" validate:"required"`
	Type            string    `reindex:"type" json:"type" validate:"required"`
	Position        string    `reindex:"position" json:"position" validate:"required"`
	DismissalReason string    `reindex:"dissmisal_reason" json:"dismissal_reason,omitempty"`
}

type PersonPost struct {
	FirstName  string    `reindex:"first_name" json:"first_name" validate:"required"`
	LastName   string    `reindex:"last_name" json:"last_name" validate:"required"`
	Username   string    `reindex:"username,hash" json:"username" validate:"required"`
	Birthdate  time.Time `reindex:"birthdate" json:"birthdate,omitempty"`
	Profession string    `reindex:"profession" json:"profession" validate:"required"`
	Sort       int       `reindex:"sort" json:"sort" validate:"required"`
	Jobs       []Job     `reindex:"jobs" json:"jobs" validate:"required"`
}

/*
Структура документа с двойной вложенностью, как в задании

Разделение на две структуры, потому что я не приудмал,
как лучше гарантировать, что пользователь в теле json не передаст _id или updated_at

И в доках echo есть что-то похожее:
https://echo.labstack.com/docs/binding#security
*/
type Person struct {
	ID 			int64 		`reindex:"id,,pk" json:"_id"`
	PersonPost
	UpdatedAt 	time.Time 	`reindex:"updated_at" json:"updated_at"`
}