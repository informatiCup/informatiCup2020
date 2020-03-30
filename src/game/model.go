package game

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

type model struct {
	game *Game
}

func (m *model) handleApplyHygienicMeasuresAction(a applyHygienicMeasuresAction) error {
	g := m.game

	// Preconditions
	c, f := g.cities[a.City]
	if !f {
		return errors.New("city does not exist")
	}
	p := 3
	if g.points < p {
		return m.insufficientPointsError(p)
	}

	c.hygiene += randomFloat(.1, .35)
	if c.hygiene > 1 {
		c.hygiene = 1
	}
	g.points -= p
	c.events = append(c.events, newHygienicMeasuresAppliedEvent(g.round))

	return nil
}

func (m *model) handleCallElectionsAction(a callElectionsAction) error {
	g := m.game

	// Preconditions
	c, f := g.cities[a.City]
	if !f {
		return errors.New("city does not exist")
	}
	p := 3
	if g.points < p {
		return m.insufficientPointsError(p)
	}

	c.government = randomFloat(.1, 1)
	g.points -= p
	c.events = append(c.events, newElectionsCalledEvent(g.round))

	return nil
}

func (m *model) handleCloseAirportAction(a closeAirportAction) error {
	g := m.game

	// Preconditions
	c, f := g.cities[a.City]
	if !f {
		return errors.New("city does not exist")
	}
	if a.Rounds < 1 {
		return errors.New("number of rounds is less than 1")
	}
	p := 5*a.Rounds + 15
	if g.points < p {
		return m.insufficientPointsError(p)
	}
	for _, e := range c.events {
		if _, ok := e.(*airportClosedEvent); ok {
			return errors.New("airport has already been closed")
		}
	}

	c.events = append(c.events, newAirportClosedEvent(g.round, g.round+a.Rounds))
	g.points -= p

	return nil
}

func (m *model) handleCloseConnectionAction(a closeConnectionAction) error {
	g := m.game

	// Preconditions
	c, f := g.cities[a.FromCity]
	if !f {
		return errors.New("city (from) does not exist")
	}
	c2, f := g.cities[a.ToCity]
	if !f {
		return errors.New("city (to) does not exist")
	}
	if a.Rounds < 1 {
		return errors.New("number of rounds is less than 1")
	}
	p := 3*a.Rounds + 3
	if g.points < p {
		return m.insufficientPointsError(p)
	}
	for _, e := range c.events {
		if v, ok := e.(*connectionClosedEvent); ok {
			if v.City == c2.name {
				return errors.New("connection has already been closed")
			}
		}
	}

	c.events = append(c.events, newConnectionClosedEvent(c2.name, g.round, g.round+a.Rounds))
	g.points -= p

	return nil
}

func (m *model) handleDevelopMedicationAction(a developMedicationAction) error {
	g := m.game

	// Preconditions
	p := 20
	if g.points < p {
		return m.insufficientPointsError(p)
	}
	ec := false
	for _, e := range g.events {
		switch v := e.(type) {
		case *medicationAvailableEvent:
			if v.Pathogen.name == a.Pathogen {
				return errors.New("medication is already available")
			}
		case *medicationInDevelopmentEvent:
			if v.Pathogen.name == a.Pathogen {
				return errors.New("medication is already in development")
			}
		case *pathogenEncounteredEvent:
			if v.Pathogen.name == a.Pathogen {
				ec = true
			}
		}
	}
	if !ec {
		return errors.New("pathogen has not been encountered")
	}

	g.events = append(g.events, newMedicationInDevelopmentEvent(g.pathogens[a.Pathogen], g.round, g.round+3))
	g.points -= p

	return nil
}

func (m *model) handleDevelopVaccineAction(a developVaccineAction) error {
	g := m.game

	// Preconditions
	p := 40
	if g.points < p {
		return m.insufficientPointsError(p)
	}
	ec := false
	for _, e := range g.events {
		switch v := e.(type) {
		case *vaccineAvailableEvent:
			if v.Pathogen.name == a.Pathogen {
				return errors.New("vaccine is already available")
			}
		case *vaccineInDevelopmentEvent:
			if v.Pathogen.name == a.Pathogen {
				return errors.New("vaccine is already in development")
			}
		case *pathogenEncounteredEvent:
			if v.Pathogen.name == a.Pathogen {
				ec = true
			}
		}
	}
	if !ec {
		return errors.New("pathogen has not been encountered")
	}

	g.points -= p
	g.events = append(g.events, newVaccineInDevelopmentEvent(g.pathogens[a.Pathogen], g.round, g.round+6))

	return nil
}

func (m *model) handleDeployMedicationAction(a deployMedicationAction) error {
	g := m.game

	// Preconditions
	c, f := g.cities[a.City]
	if !f {
		return errors.New("city does not exist")
	}
	p := 10
	if g.points < p {
		return m.insufficientPointsError(p)
	}
	va := false
	for _, e := range g.events {
		if v, ok := e.(*medicationAvailableEvent); ok {
			if v.Pathogen.name == a.Pathogen {
				va = true
				break
			}
		}
	}
	if !va {
		return errors.New("medication is not available")
	}

	r := randomFloat(.4, .6)
	pa := g.pathogens[a.Pathogen]
	for _, i := range c.infections {
		if pa.name != a.Pathogen {
			continue
		}

		n := int(float64(i.population) * r)
		i.population -= n
		c.immunities = append(c.immunities, newImmunity(pa, n))
	}
	g.points -= p
	c.events = append(c.events, newMedicationDeployedEvent(pa, g.round))

	for _, cn := range g.cityNames {
		c := g.cities[cn]
		es := make([]event, 0)
		for _, e := range c.events {
			switch v := e.(type) {
			case *outbreakEvent:
				m.handleOutbreakEvent(&es, v, c)
			default:
				es = append(es, e)
			}
		}
		c.events = es
	}

	return nil
}

func (m *model) handleDeployVaccineAction(a deployVaccineAction) error {
	g := m.game

	// Preconditions
	c, f := g.cities[a.City]
	if !f {
		return errors.New("city does not exist")
	}
	p := 5
	if g.points < p {
		return m.insufficientPointsError(p)
	}
	va := false
	for _, e := range g.events {
		if v, ok := e.(*vaccineAvailableEvent); ok {
			if v.Pathogen.name == a.Pathogen {
				va = true
				break
			}
		}
	}
	if !va {
		return errors.New("vaccine is not available")
	}

	pt := g.pathogens[a.Pathogen]
	c.immunities = append(c.immunities, newImmunity(pt, c.population-c.infectedWith(pt)-c.immuneTo(pt)))
	g.points -= p
	c.events = append(c.events, newVaccineDeployedEvent(g.pathogens[a.Pathogen], g.round))

	return nil
}

func (m *model) handleEndRoundAction(a endRoundAction) error {
	g := m.game

	g.round++
	g.points += 20

	for _, cn := range g.cityNames {
		c := g.cities[cn]
		for _, cn2 := range g.cityNames {
			c2 := g.cities[cn2]
			if c == c2 {
				continue
			}
			m.spreadToConnected(c, c2)
			m.spreadToAny(c, c2)
		}
		m.spreadWithin(c)
	}

	m.handleEvents()
	m.createEvents()

	switch {
	case g.totalPopulation() < g.initialTotalPopulation/2:
		g.outcome = gameOutcomeLoss
	case g.totalInfected() == 0:
		g.outcome = gameOutcomeWin
	}

	m.validate()

	return nil
}

func (m *model) handleExertInfluenceAction(a exertInfluenceAction) error {
	g := m.game

	// Preconditions
	c, f := g.cities[a.City]
	if !f {
		return errors.New("city does not exist")
	}
	p := 3
	if g.points < p {
		return m.insufficientPointsError(p)
	}

	c.economy = randomFloat(.1, 1)
	g.points -= p
	c.events = append(c.events, newInfluenceExertedEvent(g.round))

	return nil
}

func (m *model) handleLaunchCampaignAction(a launchCampaignAction) error {
	g := m.game

	// Preconditions
	c, f := g.cities[a.City]
	if !f {
		return errors.New("city does not exist")
	}
	p := 3
	if g.points < p {
		return m.insufficientPointsError(p)
	}

	c.awareness += randomFloat(.1, .35)
	if c.awareness > 1 {
		c.awareness = 1
	}
	g.points -= p
	c.events = append(c.events, newCampaignLaunchedEvent(g.round))

	return nil
}

func (m *model) handlePutUnderQuarantineAction(a putUnderQuarantineAction) error {
	g := m.game

	// Preconditions
	c, f := g.cities[a.City]
	if !f {
		return errors.New("city does not exist")
	}
	if a.Rounds < 1 {
		return errors.New("number of rounds is less than 1")
	}
	p := 10*a.Rounds + 20
	if g.points < p {
		return m.insufficientPointsError(p)
	}
	for _, e := range c.events {
		if _, ok := e.(*quarantineEvent); ok {
			return errors.New("city is already under quarantine")
		}
	}

	g.points -= p
	c.events = append(c.events, newQuarantineEvent(g.round, g.round+a.Rounds))

	return nil
}

func (m *model) handleAction(a action) error {
	switch v := a.(type) {
	case applyHygienicMeasuresAction:
		return m.handleApplyHygienicMeasuresAction(v)
	case callElectionsAction:
		return m.handleCallElectionsAction(v)
	case closeAirportAction:
		return m.handleCloseAirportAction(v)
	case closeConnectionAction:
		return m.handleCloseConnectionAction(v)
	case developMedicationAction:
		return m.handleDevelopMedicationAction(v)
	case developVaccineAction:
		return m.handleDevelopVaccineAction(v)
	case deployMedicationAction:
		return m.handleDeployMedicationAction(v)
	case deployVaccineAction:
		return m.handleDeployVaccineAction(v)
	case endRoundAction:
		return m.handleEndRoundAction(v)
	case exertInfluenceAction:
		return m.handleExertInfluenceAction(v)
	case launchCampaignAction:
		return m.handleLaunchCampaignAction(v)
	case putUnderQuarantineAction:
		return m.handlePutUnderQuarantineAction(v)
	default:
		panic("unhandled action type")
	}
}

func (m *model) createOutbreakEvent(c *city, p *pathogen) {
	g := m.game

	n := int(float64(c.population-c.infectedWith(p)-c.immuneTo(p)) * m.probabilityInhabitantIsInfectedDuringBreakout(c, p))
	if n <= 0 {
		return
	}
	c.infections = append(c.infections, newInfection(p, n, g.round+p.duration))

	pv := float64(n) / float64(c.population)
	if pv > 0 {
		c.events = append(c.events, newOutbreakEvent(p, pv, g.round))

		ec := false
		for _, e := range m.game.events {
			if v, ok := e.(*pathogenEncounteredEvent); ok {
				if v.Pathogen.name == p.name {
					ec = true
					break
				}
			}
		}
		if !ec {
			g.events = append(g.events, newPathogenEncounteredEvent(p, g.round))
		}
	}
}

func (m *model) createAntiVaccinationismEvent(c *city) {
	c.events = append(c.events, newAntiVaccinationismEvent(m.game.round))
}

func (m *model) createBioterrorismEvent(c *city) {
	g := m.game

	p := g.randomPathogen()
	c.events = append(c.events, newBioterrorismEvent(g.round, p))
	m.createOutbreakEvent(c, p)
}

func (m *model) createUprisingEvent(c *city) {
	c.events = append(c.events, newUprisingEvent(m.game.round, rand.Intn(c.population/10+1)))
}

func (m *model) createEconomicCrisisEvent() {
	g := m.game

	g.events = append(g.events, newEconomicCrisisEvent(g.round))
}

func (m *model) createLargeScalePanicEvent() {
	g := m.game

	g.events = append(g.events, newLargeScalePanicEvent(g.round))
}

func (m *model) createEvents() {
	g := m.game

	for _, cn := range g.cityNames {
		c := g.cities[cn]
		if randomBool(m.probabilityAntiVaccinationismOccurs(c)) {
			m.createAntiVaccinationismEvent(c)
		}
		if randomBool(m.probabilityBioterrorismOccurs(c)) {
			m.createBioterrorismEvent(c)
		}
		if randomBool(m.probabilityUprisingOccurs(c)) {
			m.createUprisingEvent(c)
		}
	}

	if randomBool(m.probabilityEconomicCrisisOccurs()) {
		m.createEconomicCrisisEvent()
	}
	if randomBool(m.probabilityLargeScalePanicOccurs()) {
		m.createLargeScalePanicEvent()
	}
}

func (m *model) handleAirportClosedEvent(es *[]event, e *airportClosedEvent, c *city) {
	if m.game.round < e.UntilRound {
		*es = append(*es, e)
	}
}

func (m *model) handleConnectionClosedEvent(es *[]event, e *connectionClosedEvent, c *city) {
	if m.game.round < e.UntilRound {
		*es = append(*es, e)
	}
}

func (m *model) handleOutbreakEvent(es *[]event, e *outbreakEvent, c *city) {
	ti := c.infectedWith(e.Pathogen)
	if ti == 0 {
		return
	}

	pv := float64(ti) / float64(c.population)
	if pv < 0.05 {
		// Prevalence threshold has been reached.
		ni := make([]*infection, 0)
		for _, i := range c.infections {
			if i.pathogen != e.Pathogen {
				ni = append(ni, i)
			}
		}
		c.infections = ni
		return
	}

	e.Prevalence = pv
	*es = append(*es, e)
}

func (m *model) handleQuarantineEvent(es *[]event, e *quarantineEvent, c *city) {
	if m.game.round < e.UntilRound {
		*es = append(*es, e)
	}
}

func (m *model) handleUprisingEvent(es *[]event, e *uprisingEvent, c *city) {
	p := c.population
	if p == 0 {
		return
	}

	if p < e.Participants {
		e.Participants = p
	}
	r := float64(e.Participants) / float64(p)
	c.government -= randomFloat(0, r/2)
	if c.government < 0 {
		c.government = 0
	}
	*es = append(*es, e)
}

func (m *model) handleVaccineInDevelopmentEvent(es *[]event, e *vaccineInDevelopmentEvent) {
	g := m.game

	if g.round < e.UntilRound {
		*es = append(*es, e)
		return
	}

	*es = append(*es, newVaccineAvailableEvent(e.Pathogen, g.round))
}

func (m *model) handleMedicationInDevelopmentEvent(es *[]event, e *medicationInDevelopmentEvent) {
	g := m.game

	if g.round < e.UntilRound {
		*es = append(*es, e)
		return
	}

	*es = append(*es, newMedicationAvailableEvent(e.Pathogen, g.round))
}

func (m *model) handleEconomicCrisisEvent(es *[]event, e *economicCrisisEvent) {
	g := m.game

	for _, cn := range g.cityNames {
		c := g.cities[cn]
		if randomBool(.3) {
			c.economy -= randomFloat(0.05, 0.1)
			if c.economy < 0 {
				c.economy = 0
			}
		}

	}
	*es = append(*es, e)
}

func (m *model) handleLargeScalePanicEvent(es *[]event, e *largeScalePanicEvent) {
	g := m.game

	for _, cn := range g.cityNames {
		c := g.cities[cn]
		if randomBool(.3) {
			c.government -= randomFloat(0.05, 0.1)
			if c.government < 0 {
				c.government = 0
			}
		}

	}
	*es = append(*es, e)
}

func (m *model) handleEvents() {
	g := m.game

	for _, cn := range g.cityNames {
		c := g.cities[cn]
		es := make([]event, 0)
		for _, e := range c.events {
			switch v := e.(type) {
			case *airportClosedEvent:
				m.handleAirportClosedEvent(&es, v, c)
			case *connectionClosedEvent:
				m.handleConnectionClosedEvent(&es, v, c)
			case *outbreakEvent:
				m.handleOutbreakEvent(&es, v, c)
			case *quarantineEvent:
				m.handleQuarantineEvent(&es, v, c)
			case *uprisingEvent:
				m.handleUprisingEvent(&es, v, c)
			default:
				es = append(es, e)
			}
		}
		c.events = es
	}

	es := make([]event, 0)
	for _, e := range g.events {
		switch v := e.(type) {
		case *vaccineInDevelopmentEvent:
			m.handleVaccineInDevelopmentEvent(&es, v)
		case *medicationInDevelopmentEvent:
			m.handleMedicationInDevelopmentEvent(&es, v)
		case *economicCrisisEvent:
			m.handleEconomicCrisisEvent(&es, v)
		case *largeScalePanicEvent:
			m.handleLargeScalePanicEvent(&es, v)
		default:
			es = append(es, e)
		}
	}
	g.events = es
}

func (m *model) probabilityPathogenSpreadsToAny(c *city, c2 *city, p *pathogen) float64 {
	g := m.game

	for _, e := range c.events {
		if _, ok := e.(*quarantineEvent); ok {
			return 0
		}
	}

	for _, e := range c2.events {
		if _, ok := e.(*quarantineEvent); ok {
			return 0
		}
		if _, ok := e.(*outbreakEvent); ok {
			return 0
		}
	}

	return p.mobility * (1 - c.economy/20 - c2.economy/20 - math.Pow(g.distances[c][c2], .25)/1.125)
}

func (m *model) probabilityPathogenSpreadsToConnected(c *city, c2 *city, p *pathogen) float64 {
	if !c.connections[c2.name] {
		return 0
	}

	for _, e := range c.events {
		if _, ok := e.(*airportClosedEvent); ok {
			return 0
		}
		if v, ok := e.(*connectionClosedEvent); ok {
			if v.City == c2.name {
				return 0
			}
		}
		if _, ok := e.(*quarantineEvent); ok {
			return 0
		}
	}

	for _, e := range c2.events {
		if _, ok := e.(*outbreakEvent); ok {
			return 0
		}
		if _, ok := e.(*airportClosedEvent); ok {
			return 0
		}
		if _, ok := e.(*quarantineEvent); ok {
			return 0
		}
	}

	return p.infectivity * (1 - c.economy/20 - c2.economy/10 - c.government/20 - c2.government/10 - c2.hygiene/10 - c2.awareness/10)
}

func (m *model) probabilityInhabitantIsInfectedDuringBreakout(c *city, p *pathogen) float64 {
	return p.infectivity * (1 - c.economy/10 - c.government/10 - c.hygiene/5)
}

func (m *model) probabilityInhabitantInfectsOther(c *city, p *pathogen) float64 {
	return p.infectivity * (1 - c.awareness/10 - c.government/20 - c.hygiene/10)
}

func (m *model) probabilityPathogenKillsInhabitant(c *city, p *pathogen) float64 {
	return p.lethality * (1 - c.economy/6 - c.government/9 - c.hygiene/3)
}

func (m *model) probabilityAntiVaccinationismOccurs(c *city) float64 {
	for _, e := range c.events {
		if _, ok := e.(*antiVaccinationismEvent); ok {
			return 0
		}
	}

	return .0005 * (1 - c.awareness)
}

func (m *model) probabilityBioterrorismOccurs(c *city) float64 {
	for _, e := range c.events {
		if _, ok := e.(*bioterrorismEvent); ok {
			return 0
		}
		if _, ok := e.(*outbreakEvent); ok {
			return 0
		}
	}

	return .0005 * (1 - c.government)
}

func (m *model) probabilityUprisingOccurs(c *city) float64 {
	g := m.game

	var qe *quarantineEvent
	for _, e := range c.events {
		if _, ok := e.(*uprisingEvent); ok {
			return 0
		}
		if v, ok := e.(*quarantineEvent); ok {
			qe = v
		}
	}
	if qe == nil {
		return 0
	}

	rs := g.round - qe.SinceRound
	return .05 * float64(rs)
}

func (m *model) probabilityEconomicCrisisOccurs() float64 {
	g := m.game

	for _, e := range g.events {
		if _, ok := e.(*economicCrisisEvent); ok {
			return 0
		}
	}

	return .05 * float64(g.totalInfected()) / float64(g.totalPopulation())
}

func (m *model) probabilityLargeScalePanicOccurs() float64 {
	g := m.game

	for _, e := range g.events {
		if _, ok := e.(*largeScalePanicEvent); ok {
			return 0
		}
	}

	return .05 * float64(g.totalInfected()) / float64(g.totalPopulation())
}

func (m *model) spreadToAny(c *city, c2 *city) {
	for _, e := range c.events {
		if oe, ok := e.(*outbreakEvent); ok {
			if randomBool(m.probabilityPathogenSpreadsToAny(c, c2, oe.Pathogen)) {
				m.createOutbreakEvent(c2, oe.Pathogen)
			}
		}
	}
}

func (m *model) spreadToConnected(c *city, c2 *city) {
	for _, e := range c.events {
		if oe, ok := e.(*outbreakEvent); ok {
			if randomBool(m.probabilityPathogenSpreadsToConnected(c, c2, oe.Pathogen)) {
				m.createOutbreakEvent(c2, oe.Pathogen)
			}
		}
	}
}

func (m *model) spreadWithin(c *city) {
	g := m.game

	ni := make([]*infection, 0)
	for _, i := range c.infections {
		p := i.pathogen
		if i.untilRound == g.round {
			d := int(float64(i.population) * m.probabilityPathogenKillsInhabitant(c, p))
			c.immunities = append(c.immunities, newImmunity(p, i.population-d))
			c.population -= d
			continue
		}

		ni = append(ni, i)
	}
	c.infections = ni

	for p := range c.infectionsMap() {
		n := int(float64(c.population-c.infectedWith(p)-c.immuneTo(p)) * m.probabilityInhabitantInfectsOther(c, p))
		if n > 0 {
			c.infections = append(c.infections, newInfection(p, n, g.round+p.duration))
		}
	}
}

func (m *model) insufficientPointsError(r int) error {
	return fmt.Errorf("%d points are required to perform this action, but only %d points are available", r, m.game.points)
}

func (m *model) initialize() {
	g := m.game

	cs := g.randomCities(2, 3)
	ps := make([]*pathogen, len(cs))
	for {
		l := false
		for i := range ps {
			ps[i] = g.randomPathogen()
		}
		for _, p := range ps {
			if p.lethality > .5 {
				l = true
			}
		}
		if l {
			break
		}
	}
	for i := range cs {
		m.createOutbreakEvent(cs[i], ps[i])
	}
}

func (m *model) validate() {
	g := m.game

	for _, c := range g.cities {
		im := c.infectionsMap()
		if len(im) > 1 {
			panic("multiple pathogens among infected")
		}
		for p := range im {
			if c.population-c.infectedWith(p)-c.immuneTo(p) < 0 {
				panic("sum of infected and immune exceeds total")
			}
		}
	}
}

func newModel(g *Game) *model {
	m := &model{game: g}
	m.initialize()

	return m
}
