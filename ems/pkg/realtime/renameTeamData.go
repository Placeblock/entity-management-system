package realtime

type RenameTeamData struct {
	TeamID  uint   `json:"team_id"`
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}
