package handballnet

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

func mapTeamMembers(players []player, staffMembers []staffMember, teamID int64) []db.UpsertTeamMemberParams {
	var teamMembers []db.UpsertTeamMemberParams

	mappedPlayers := mapPlayers(players, teamID)

	mappedStaffMembers := mapStaffMembers(staffMembers, teamID)

	teamMembers = append(teamMembers, mappedPlayers...)
	teamMembers = append(teamMembers, mappedStaffMembers...)

	return teamMembers
}

func mapPlayers(players []player, teamID int64) []db.UpsertTeamMemberParams {
	var teamMembers []db.UpsertTeamMemberParams

	for _, player := range players {
		mappedTeamMember := mapPlayer(player, teamID)
		teamMembers = append(teamMembers, mappedTeamMember)
	}

	return teamMembers
}

func mapPlayer(player player, teamID int64) db.UpsertTeamMemberParams {
	stringNumber := strconv.Itoa(player.Number)
	return db.UpsertTeamMemberParams{
		Uid:    getTeamMemberUID(teamID, player.ID, player.Firstname, player.Lastname, stringNumber),
		TeamID: teamID,
		Name:   strings.Join([]string{player.Firstname, player.Lastname}, " "),
		Number: stringNumber,
		Type:   string(model.TeamMemberTypePlayer),
	}
}

func mapStaffMembers(staffMembers []staffMember, teamID int64) []db.UpsertTeamMemberParams {
	var teamMembers []db.UpsertTeamMemberParams

	for _, staffMember := range staffMembers {
		mappedTeamMember := mapStaffMember(staffMember, teamID)
		teamMembers = append(teamMembers, mappedTeamMember)
	}

	return teamMembers
}

func mapStaffMember(staffMember staffMember, teamID int64) db.UpsertTeamMemberParams {
	return db.UpsertTeamMemberParams{
		Uid:    getTeamMemberUID(teamID, staffMember.ID, staffMember.Firstname, staffMember.Lastname, staffMember.Position),
		TeamID: teamID,
		Name:   strings.Join([]string{staffMember.Firstname, staffMember.Lastname}, " "),
		Number: staffMember.Position,
		Type:   string(model.TeamMemberTypeOfficial),
	}
}

func getTeamMemberUID(teamID int64, teamMemberUID string, firstname string, lastname string, number string) string {
	reg := regexp.MustCompile(`(\.-\d+)$`)
	matches := reg.Match([]byte(teamMemberUID))
	if !matches {
		return teamMemberUID
	}

	namespace := reg.ReplaceAllString(teamMemberUID, "$1W")
	preparedFirstname := strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(firstname), ".", ""), " ", "")
	preparedLastname := strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(lastname), ".", ""), " ", "")
	return strings.ToLower(strings.Join([]string{namespace, strconv.Itoa(int(teamID)), number, preparedFirstname, preparedLastname}, "."))
}
