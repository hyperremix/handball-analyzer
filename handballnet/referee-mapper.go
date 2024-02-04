package handballnet

import (
	"strings"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

func mapReferees(referees []referee) []db.UpsertRefereeParams {
	var mappedReferees []db.UpsertRefereeParams

	for _, referee := range referees {
		mappedReferee := mapReferee(referee)
		mappedReferees = append(mappedReferees, mappedReferee)
	}

	return mappedReferees
}

func mapReferee(referee referee) db.UpsertRefereeParams {
	return db.UpsertRefereeParams{
		Uid:  referee.ID,
		Name: strings.Join([]string{referee.Firstname, referee.Lastname}, " "),
		Type: string(mapRefereeType(referee.Position)),
	}
}

func mapRefereeType(position string) model.RefereeType {
	switch position {
	case "Schiedsrichter (1)":
		return model.RefereeTypeFirstReferee
	case "Schiedsrichter (2)":
		return model.RefereeTypeSecondReferee
	case "Sekretär":
		return model.RefereeTypeSecretary
	case "Zeitnehmer":
		return model.RefereeTypeTimekeeper
	default:
		return model.RefereeTypeFirstReferee
	}
}
