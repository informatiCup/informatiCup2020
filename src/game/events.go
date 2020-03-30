package game

type event interface{}

type outbreakEvent struct {
	Type       string    `json:"type"`
	Pathogen   *pathogen `json:"pathogen"`
	Prevalence float64   `json:"prevalence"`
	SinceRound int       `json:"sinceRound"`
}

func newOutbreakEvent(p *pathogen, pv float64, r int) *outbreakEvent {
	return &outbreakEvent{
		Type:       "outbreak",
		Pathogen:   p,
		Prevalence: pv,
		SinceRound: r,
	}
}

type airportClosedEvent struct {
	Type       string `json:"type"`
	SinceRound int    `json:"sinceRound"`
	UntilRound int    `json:"untilRound"`
}

func newAirportClosedEvent(sr, ur int) *airportClosedEvent {
	return &airportClosedEvent{
		Type:       "airportClosed",
		SinceRound: sr,
		UntilRound: ur,
	}
}

type connectionClosedEvent struct {
	Type       string `json:"type"`
	City       string `json:"city"`
	SinceRound int    `json:"sinceRound"`
	UntilRound int    `json:"untilRound"`
}

func newConnectionClosedEvent(c string, sr, ur int) *connectionClosedEvent {
	return &connectionClosedEvent{
		Type:       "connectionClosed",
		City:       c,
		SinceRound: sr,
		UntilRound: ur,
	}
}

type quarantineEvent struct {
	Type       string `json:"type"`
	SinceRound int    `json:"sinceRound"`
	UntilRound int    `json:"untilRound"`
}

func newQuarantineEvent(sr, ur int) *quarantineEvent {
	return &quarantineEvent{
		Type:       "quarantine",
		SinceRound: sr,
		UntilRound: ur,
	}
}

type pathogenEncounteredEvent struct {
	Type     string    `json:"type"`
	Pathogen *pathogen `json:"pathogen"`
	Round    int       `json:"round"`
}

func newPathogenEncounteredEvent(p *pathogen, r int) *pathogenEncounteredEvent {
	return &pathogenEncounteredEvent{
		Type:     "pathogenEncountered",
		Pathogen: p,
		Round:    r,
	}
}

type vaccineInDevelopmentEvent struct {
	Type       string    `json:"type"`
	Pathogen   *pathogen `json:"pathogen"`
	SinceRound int       `json:"sinceRound"`
	UntilRound int       `json:"untilRound"`
}

func newVaccineInDevelopmentEvent(p *pathogen, sr, ur int) *vaccineInDevelopmentEvent {
	return &vaccineInDevelopmentEvent{
		Type:       "vaccineInDevelopment",
		Pathogen:   p,
		SinceRound: sr,
		UntilRound: ur,
	}
}

type vaccineAvailableEvent struct {
	Type       string    `json:"type"`
	Pathogen   *pathogen `json:"pathogen"`
	SinceRound int       `json:"sinceRound"`
}

func newVaccineAvailableEvent(p *pathogen, sr int) *vaccineAvailableEvent {
	return &vaccineAvailableEvent{
		Type:       "vaccineAvailable",
		Pathogen:   p,
		SinceRound: sr,
	}
}

type medicationInDevelopmentEvent struct {
	Type       string    `json:"type"`
	Pathogen   *pathogen `json:"pathogen"`
	SinceRound int       `json:"sinceRound"`
	UntilRound int       `json:"untilRound"`
}

func newMedicationInDevelopmentEvent(p *pathogen, sr, ur int) *medicationInDevelopmentEvent {
	return &medicationInDevelopmentEvent{
		Type:       "medicationInDevelopment",
		Pathogen:   p,
		SinceRound: sr,
		UntilRound: ur,
	}
}

type medicationAvailableEvent struct {
	Type       string    `json:"type"`
	Pathogen   *pathogen `json:"pathogen"`
	SinceRound int       `json:"sinceRound"`
}

func newMedicationAvailableEvent(p *pathogen, sr int) *medicationAvailableEvent {
	return &medicationAvailableEvent{
		Type:       "medicationAvailable",
		Pathogen:   p,
		SinceRound: sr,
	}
}

type medicationDeployedEvent struct {
	Type     string    `json:"type"`
	Pathogen *pathogen `json:"pathogen"`
	Round    int       `json:"round"`
}

func newMedicationDeployedEvent(p *pathogen, r int) *medicationDeployedEvent {
	return &medicationDeployedEvent{
		Type:     "medicationDeployed",
		Pathogen: p,
		Round:    r,
	}
}

type vaccineDeployedEvent struct {
	Type     string    `json:"type"`
	Pathogen *pathogen `json:"pathogen"`
	Round    int       `json:"round"`
}

func newVaccineDeployedEvent(p *pathogen, r int) *vaccineDeployedEvent {
	return &vaccineDeployedEvent{
		Type:     "vaccineDeployed",
		Pathogen: p,
		Round:    r,
	}
}

type antiVaccinationismEvent struct {
	Type       string `json:"type"`
	SinceRound int    `json:"sinceRound"`
}

func newAntiVaccinationismEvent(r int) *antiVaccinationismEvent {
	return &antiVaccinationismEvent{
		Type:       "antiVaccinationism",
		SinceRound: r,
	}
}

type bioterrorismEvent struct {
	Type     string    `json:"type"`
	Round    int       `json:"round"`
	Pathogen *pathogen `json:"pathogen"`
}

func newBioterrorismEvent(r int, p *pathogen) *bioterrorismEvent {
	return &bioterrorismEvent{
		Type:     "bioTerrorism",
		Round:    r,
		Pathogen: p,
	}
}

type uprisingEvent struct {
	Type         string `json:"type"`
	SinceRound   int    `json:"sinceRound"`
	Participants int    `json:"participants"`
}

func newUprisingEvent(r int, p int) *uprisingEvent {
	return &uprisingEvent{
		Type:         "uprising",
		SinceRound:   r,
		Participants: p,
	}
}

type economicCrisisEvent struct {
	Type       string `json:"type"`
	SinceRound int    `json:"sinceRound"`
}

func newEconomicCrisisEvent(r int) *economicCrisisEvent {
	return &economicCrisisEvent{
		Type:       "economicCrisis",
		SinceRound: r,
	}
}

type largeScalePanicEvent struct {
	Type       string `json:"type"`
	SinceRound int    `json:"sinceRound"`
}

func newLargeScalePanicEvent(r int) *largeScalePanicEvent {
	return &largeScalePanicEvent{
		Type:       "largeScalePanic",
		SinceRound: r,
	}
}

type hygienicMeasuresAppliedEvent struct {
	Type  string `json:"type"`
	Round int    `json:"round"`
}

func newHygienicMeasuresAppliedEvent(r int) *hygienicMeasuresAppliedEvent {
	return &hygienicMeasuresAppliedEvent{
		Type:  "hygienicMeasuresApplied",
		Round: r,
	}
}

type campaignLaunchedEvent struct {
	Type  string `json:"type"`
	Round int    `json:"round"`
}

func newCampaignLaunchedEvent(r int) *campaignLaunchedEvent {
	return &campaignLaunchedEvent{
		Type:  "campaignLaunched",
		Round: r,
	}
}

type influenceExertedEvent struct {
	Type  string `json:"type"`
	Round int    `json:"round"`
}

func newInfluenceExertedEvent(r int) *influenceExertedEvent {
	return &influenceExertedEvent{
		Type:  "influenceExerted",
		Round: r,
	}
}

type electionsCalledEvent struct {
	Type  string `json:"type"`
	Round int    `json:"round"`
}

func newElectionsCalledEvent(r int) *electionsCalledEvent {
	return &electionsCalledEvent{
		Type:  "electionsCalled",
		Round: r,
	}
}
