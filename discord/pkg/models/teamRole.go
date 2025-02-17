package models

import "github.com/Placeblock/nostalgicraft-ems/pkg/models"

type TeamRole struct {
	TeamID uint
	Team   models.Team
	RoleID string
}
