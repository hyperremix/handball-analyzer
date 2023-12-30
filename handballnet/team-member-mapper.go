package handballnet

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

func mapTeamMembers(players []player, staffMembers []staffMember, teamID uint) ([]model.TeamMember, error) {
	var teamMembers []model.TeamMember

	mappedPlayers, err := mapPlayers(players, teamID)
	if err != nil {
		return []model.TeamMember{}, err
	}

	mappedStaffMembers, err := mapStaffMembers(staffMembers, teamID)
	if err != nil {
		return []model.TeamMember{}, err
	}

	teamMembers = append(teamMembers, mappedPlayers...)
	teamMembers = append(teamMembers, mappedStaffMembers...)

	return teamMembers, nil
}

func mapPlayers(players []player, teamID uint) ([]model.TeamMember, error) {
	var teamMembers []model.TeamMember

	for _, player := range players {
		mappedTeamMember, err := mapPlayer(player, teamID)
		if err != nil {
			return []model.TeamMember{}, err
		}
		teamMembers = append(teamMembers, mappedTeamMember)
	}

	return teamMembers, nil
}

func mapPlayer(player player, teamID uint) (model.TeamMember, error) {
	var teamMember model.TeamMember

	stringNumber := strconv.Itoa(player.Number)

	if err := db.Get().Where(&model.TeamMember{
		UID: getTeamMemberUID(teamID, player.ID, player.Firstname, player.Lastname, stringNumber),
	}).Attrs(model.TeamMember{
		TeamID: teamID,
		Name:   strings.Join([]string{player.Firstname, player.Lastname}, " "),
		Number: stringNumber,
		Type:   model.TeamMemberTypePlayer,
	}).FirstOrCreate(&teamMember).Error; err != nil {
		return model.TeamMember{}, fmt.Errorf("error creating player: %s", err)
	}

	return teamMember, nil
}

func mapStaffMembers(staffMembers []staffMember, teamID uint) ([]model.TeamMember, error) {
	var teamMembers []model.TeamMember

	for _, staffMember := range staffMembers {
		mappedTeamMember, err := mapStaffMember(staffMember, teamID)
		if err != nil {
			return []model.TeamMember{}, err
		}
		teamMembers = append(teamMembers, mappedTeamMember)
	}

	return teamMembers, nil
}

func mapStaffMember(staffMember staffMember, teamID uint) (model.TeamMember, error) {
	var teamMember model.TeamMember

	if err := db.Get().Where(&model.TeamMember{
		UID: getTeamMemberUID(teamID, staffMember.ID, staffMember.Firstname, staffMember.Lastname, staffMember.Position),
	}).Attrs(model.TeamMember{
		TeamID: teamID,
		Name:   strings.Join([]string{staffMember.Firstname, staffMember.Lastname}, " "),
		Number: staffMember.Position,
		Type:   model.TeamMemberTypeOfficial,
	}).FirstOrCreate(&teamMember).Error; err != nil {
		return model.TeamMember{}, fmt.Errorf("error creating staff member: %s", err)
	}

	return teamMember, nil
}

func getTeamMemberUID(teamID uint, teamMemberUID string, firstname string, lastname string, number string) string {
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
