package realtime

import "github.com/Placeblock/nostalgicraft-ems/pkg/models"

type CreateTeamData struct {
	Team   models.Team   `json:"team"`
	Member models.Member `json:"member"`
}
