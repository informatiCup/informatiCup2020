package game

import "encoding/json"

type pathogen struct {
	name        string
	infectivity float64
	mobility    float64
	duration    int
	lethality   float64
}

func (p *pathogen) fuzzyInfectivity() string {
	return defaultFuzzy(p.infectivity)
}

func (p *pathogen) fuzzyDuration() string {
	switch {
	case p.duration <= 1:
		return "--"
	case p.duration <= 3:
		return "-"
	case p.duration <= 7:
		return "o"
	case p.duration <= 10:
		return "+"
	default:
		return "++"
	}
}

func (p *pathogen) fuzzyLethality() string {
	return defaultFuzzy(p.lethality)
}

func (p *pathogen) fuzzyMobility() string {
	return defaultFuzzy(p.mobility)
}

// MarshalJSON is required to implement json.Marshaler.
// See https://golang.org/pkg/encoding/json/.
func (p *pathogen) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name        string `json:"name"`
		Infectivity string `json:"infectivity"`
		Mobility    string `json:"mobility"`
		Duration    string `json:"duration"`
		Lethality   string `json:"lethality"`
	}{
		Name:        p.name,
		Infectivity: p.fuzzyInfectivity(),
		Mobility:    p.fuzzyMobility(),
		Duration:    p.fuzzyDuration(),
		Lethality:   p.fuzzyLethality(),
	})
}

func newPathogen(n string, i, l, m float64, d int) *pathogen {
	return &pathogen{
		name:        n,
		infectivity: i,
		lethality:   l,
		mobility:    m,
		duration:    d,
	}
}
