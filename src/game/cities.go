package game

import (
	"encoding/json"
	"math"
)

type infection struct {
	pathogen   *pathogen
	population int
	untilRound int
}

func newInfection(pa *pathogen, po, r int) *infection {
	return &infection{
		pathogen:   pa,
		population: po,
		untilRound: r,
	}
}

type immunity struct {
	pathogen   *pathogen
	population int
}

func newImmunity(pa *pathogen, po int) *immunity {
	return &immunity{
		pathogen:   pa,
		population: po,
	}
}

type city struct {
	name        string
	latitude    float64
	longitude   float64
	population  int
	connections map[string]bool
	events      []event
	economy     float64
	government  float64
	hygiene     float64
	awareness   float64
	infections  []*infection
	immunities  []*immunity
}

func (c *city) fuzzyEconomy() string {
	return defaultFuzzy(c.economy)
}

func (c *city) fuzzyGovernment() string {
	return defaultFuzzy(c.government)
}

func (c *city) fuzzyHygiene() string {
	return defaultFuzzy(c.hygiene)
}

func (c *city) fuzzyAwareness() string {
	return defaultFuzzy(c.awareness)
}

func (c *city) distanceTo(c2 *city) float64 {
	if c == c2 {
		return 0
	}

	rad := math.Pi / 180
	φ1, φ2, Δλ := c.latitude*rad, c2.latitude*rad, (c2.longitude-c.longitude)*rad
	return math.Acos(math.Sin(φ1)*math.Sin(φ2)+math.Cos(φ1)*math.Cos(φ2)*math.Cos(Δλ)) * 6371
}

func (c *city) infectedWith(p *pathogen) int {
	t := 0
	for _, i := range c.infections {
		if i.pathogen == p {
			t += i.population
		}
	}

	return t
}

func (c *city) infectionsMap() map[*pathogen]bool {
	ps := make(map[*pathogen]bool)
	for _, i := range c.infections {
		ps[i.pathogen] = true
	}

	return ps
}

func (c *city) immuneTo(p *pathogen) int {
	t := 0
	for _, i := range c.immunities {
		if i.pathogen == p {
			t += i.population
		}
	}

	return t
}

func (c *city) immunitiesMap() map[*pathogen]bool {
	ps := make(map[*pathogen]bool)
	for _, i := range c.immunities {
		ps[i.pathogen] = true
	}

	return ps
}

// MarshalJSON is required to implement json.Marshaler.
// See https://golang.org/pkg/encoding/json/.
func (c *city) MarshalJSON() ([]byte, error) {
	cs := make([]string, len(c.connections))
	i := 0
	for c := range c.connections {
		cs[i] = c
		i++
	}

	return json.Marshal(struct {
		Name        string   `json:"name"`
		Latitude    float64  `json:"latitude"`
		Longitude   float64  `json:"longitude"`
		Population  int      `json:"population"`
		Connections []string `json:"connections"`
		Events      []event  `json:"events,omitempty"`
		Economy     string   `json:"economy"`
		Government  string   `json:"government"`
		Hygiene     string   `json:"hygiene"`
		Awareness   string   `json:"awareness"`
	}{
		Name:        c.name,
		Latitude:    c.latitude,
		Longitude:   c.longitude,
		Population:  c.population,
		Connections: cs,
		Events:      c.events,
		Economy:     c.fuzzyEconomy(),
		Government:  c.fuzzyGovernment(),
		Hygiene:     c.fuzzyHygiene(),
		Awareness:   c.fuzzyAwareness(),
	})
}

func newCity(n string, lat, lon float64, p int, e, g, h, a float64, cs ...string) *city {
	csm := make(map[string]bool, len(cs))
	for i := 0; i < len(cs); i++ {
		csm[cs[i]] = true
	}

	return &city{
		name:        n,
		latitude:    lat,
		longitude:   lon,
		population:  p,
		connections: csm,
		events:      make([]event, 0),
		economy:     e,
		government:  g,
		hygiene:     h,
		awareness:   a,
		infections:  make([]*infection, 0),
		immunities:  make([]*immunity, 0),
	}
}

func maximumDistanceBetweenCities() float64 {
	return (&city{}).distanceTo(&city{longitude: 180})
}

func reducedCities() []*city {
	return []*city{
		newCity("Berlin", 52.520861, 13.409419, 3748, 0.82, 0.658, 0.689, 0.849, "Hamburg", "Köln", "München", "Roma"),
		newCity("Hamburg", 53.54845, 9.978514, 1822, 0.892, 0.706, 0.934, 0.684, "Berlin", "Köln", "München"),
		newCity("Köln", 50.941441, 6.958324, 1085, 0.925, 0.776, 0.73, 0.789, "Berlin", "Hamburg"),
		newCity("München", 48.137687, 11.579932, 1451, 0.926, 0.915, 0.949, 0.684, "Berlin", "Hamburg"),
		newCity("Roma", 41.902337, 12.453997, 2872, 0.752, 0.821, 0.653, 0.661, "Berlin"),
	}
}
