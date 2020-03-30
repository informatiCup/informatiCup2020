package game

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

const (
	actionTypeApplyHygienicMeasures = "applyHygienicMeasures"
	actionTypeCallElections         = "callElections"
	actionTypeCloseAirport          = "closeAirport"
	actionTypeCloseConnection       = "closeConnection"
	actionTypeDevelopMedication     = "developMedication"
	actionTypeDevelopVaccine        = "developVaccine"
	actionTypeDeployMedication      = "deployMedication"
	actionTypeDeployVaccine         = "deployVaccine"
	actionTypeEndRound              = "endRound"
	actionTypeExertInfluence        = "exertInfluence"
	actionTypeLaunchCampaign        = "launchCampaign"
	actionTypePutUnderQuarantine    = "putUnderQuarantine"
)

type action interface{}

type baseAction struct {
	Type string `json:"type"`
}

type applyHygienicMeasuresAction struct {
	City string `json:"city"`
}

type callElectionsAction struct {
	City string `json:"city"`
}

type closeAirportAction struct {
	City   string `json:"city"`
	Rounds int    `json:"rounds"`
}

type closeConnectionAction struct {
	FromCity string `json:"fromCity"`
	ToCity   string `json:"toCity"`
	Rounds   int    `json:"rounds"`
}

type developMedicationAction struct {
	Pathogen string `json:"pathogen"`
}

type developVaccineAction struct {
	Pathogen string `json:"pathogen"`
}

type deployMedicationAction struct {
	Pathogen string `json:"pathogen"`
	City     string `json:"city"`
}

type deployVaccineAction struct {
	Pathogen string `json:"pathogen"`
	City     string `json:"city"`
}

type endRoundAction struct{}

type exertInfluenceAction struct {
	City string `json:"city"`
}

type launchCampaignAction struct {
	City string `json:"city"`
}

type putUnderQuarantineAction struct {
	City   string `json:"city"`
	Rounds int    `json:"rounds"`
}

func decodeAction(r io.Reader) (action, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var ba baseAction
	if err := json.Unmarshal(bs, &ba); err != nil {
		return nil, errors.New("failed to decode action")
	}

	if ba.Type == "" {
		return nil, errors.New("action type is missing")
	}

	d := json.NewDecoder(bytes.NewReader(bs))
	switch ba.Type {
	case actionTypeApplyHygienicMeasures:
		a := applyHygienicMeasuresAction{}
		d.Decode(&a)
		return a, nil
	case actionTypeCallElections:
		a := callElectionsAction{}
		d.Decode(&a)
		return a, nil
	case actionTypeCloseAirport:
		a := closeAirportAction{}
		d.Decode(&a)
		return a, nil
	case actionTypeCloseConnection:
		a := closeConnectionAction{}
		d.Decode(&a)
		return a, nil
	case actionTypeDevelopMedication:
		a := developMedicationAction{}
		d.Decode(&a)
		return a, nil
	case actionTypeDevelopVaccine:
		a := developVaccineAction{}
		d.Decode(&a)
		return a, nil
	case actionTypeDeployMedication:
		a := deployMedicationAction{}
		d.Decode(&a)
		return a, nil
	case actionTypeDeployVaccine:
		a := deployVaccineAction{}
		d.Decode(&a)
		return a, nil
	case actionTypeEndRound:
		a := endRoundAction{}
		d.Decode(&a)
		return a, nil
	case actionTypeExertInfluence:
		a := exertInfluenceAction{}
		d.Decode(&a)
		return a, nil
	case actionTypeLaunchCampaign:
		a := launchCampaignAction{}
		d.Decode(&a)
		return a, nil
	case actionTypePutUnderQuarantine:
		a := putUnderQuarantineAction{}
		d.Decode(&a)
		return a, nil
	default:
		return nil, fmt.Errorf("action type '%s' is unknown", ba.Type)
	}
}
