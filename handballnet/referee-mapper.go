package handballnet

import (
	"fmt"
	"strings"

	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
)

func mapReferees(referees []referee) ([]model.Referee, error) {
	var mappedReferees []model.Referee

	for _, referee := range referees {
		mappedReferee, err := mapReferee(referee)
		if err != nil {
			return mappedReferees, err
		}
		mappedReferees = append(mappedReferees, mappedReferee)
	}

	return mappedReferees, nil
}

func mapReferee(referee referee) (model.Referee, error) {
	var mappedReferee model.Referee

	if err := db.Get().Where(&model.Referee{UID: referee.ID}).Attrs(model.Referee{
		Name: strings.Join([]string{referee.Firstname, referee.Lastname}, " "),
		Type: mapRefereeType(referee.Position),
	}).FirstOrCreate(&mappedReferee).Error; err != nil {
		return model.Referee{}, fmt.Errorf("error creating referee: %s", err)
	}

	return mappedReferee, nil
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
