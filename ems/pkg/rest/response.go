package rest

import "github.com/Placeblock/nostalgicraft-ems/pkg/models"

type Response struct {
	Data interface{} `json:"data"`
}

type CreateTeamData struct {
	Team   models.Team   `json:"team"`
	Member models.Member `json:"member"`
}
