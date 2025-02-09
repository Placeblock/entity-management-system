package rest

import "github.com/codelix/ems/pkg/models"

type Response struct {
	Data interface{} `json:"data"`
}

type CreateTeamData struct {
	Team   models.Team   `json:"team"`
	Member models.Member `json:"member"`
}
