package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const (
	gameOutcomePending = "pending"
	gameOutcomeWin     = "win"
	gameOutcomeLoss    = "loss"
)

// Game represents a game.
type Game struct {
	outcome string
	round   int
	points  int
	cities  map[string]*city
	events  []event
	error   string

	model         *model
	cityNames     []string
	pathogenNames []string
	pathogens     map[string]*pathogen
	distances     map[*city]map[*city]float64

	initialTotalPopulation int
}

func (g *Game) randomCity() *city {
	return g.cities[g.cityNames[rand.Intn(len(g.cityNames))]]
}

func (g *Game) randomCities(min, max int) []*city {
	if min < 1 || max > len(g.cities) || min > max {
		panic(fmt.Sprintf("invalid arguments: %d, %d", min, max))
	}

	ns := g.cityNames
	rand.Shuffle(len(ns), func(i, j int) {
		ns[i], ns[j] = ns[j], ns[i]
	})
	ns = ns[0 : rand.Intn(max-min+1)+min]
	cs := make([]*city, len(ns))
	for i, n := range ns {
		cs[i] = g.cities[n]
	}

	return cs
}

func (g *Game) randomPathogen() *pathogen {
	return g.pathogens[g.pathogenNames[rand.Intn(len(g.pathogenNames))]]
}

func (g *Game) totalPopulation() int {
	t := 0
	for _, c := range g.cities {
		t += c.population
	}

	return t
}

func (g *Game) totalInfected() int {
	t := 0
	for _, c := range g.cities {
		for _, i := range c.infections {
			t += i.population
		}
	}

	return t
}

func (g *Game) initialize(cs []*city, ps []*pathogen) {
	g.cities = make(map[string]*city, len(cs))
	g.cityNames = make([]string, len(cs))
	g.distances = make(map[*city]map[*city]float64, 0)
	md := maximumDistanceBetweenCities()
	for i, c := range cs {
		if g.cities[c.name] != nil {
			panic(fmt.Sprintf("cs contains two cities with name '%s'", c.name))
		}
		g.cities[c.name] = c
		g.cityNames[i] = c.name
		g.distances[c] = make(map[*city]float64, len(cs))
		for _, c2 := range cs {
			nd := c.distanceTo(c2) / md
			if nd > 1 {
				panic(fmt.Sprintf("normalized distance between '%s' and '%s' is greater than 1", c.name, c2.name))
			}
			g.distances[c][c2] = nd
		}
	}

	g.pathogens = make(map[string]*pathogen, len(ps))
	g.pathogenNames = make([]string, len(ps))
	for i, p := range ps {
		if g.pathogens[p.name] != nil {
			panic(fmt.Sprintf("ps contains two pathogens with name '%s'", p.name))
		}
		g.pathogens[p.name] = p
		g.pathogenNames[i] = p.name
	}

	g.initialTotalPopulation = g.totalPopulation()
	g.events = make([]event, 0)

	g.validate()
}

func (g *Game) validate() {
	for n, c := range g.cities {
		acm := make(map[string]bool, 0)
		for cn := range c.connections {
			if n == cn {
				panic(fmt.Sprintf("'%s' is connected to itself", n))
			}
			if g.cities[cn] == nil {
				panic(fmt.Sprintf("'%s' is connected to unknown city '%s'", n, cn))
			}
			if acm[cn] {
				panic(fmt.Sprintf("'%s' is connected to '%s' twice", n, cn))
			}
			acm[cn] = true
		}
	}
}

// MarshalJSON is required to implement json.Marshaler.
// See https://golang.org/pkg/encoding/json/.
func (g *Game) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Outcome string           `json:"outcome"`
		Round   int              `json:"round"`
		Points  int              `json:"points"`
		Cities  map[string]*city `json:"cities"`
		Events  []event          `json:"events"`
		Error   string           `json:"error,omitempty"`
	}{
		Outcome: g.outcome,
		Round:   g.round,
		Points:  g.points,
		Cities:  g.cities,
		Events:  g.events,
		Error:   g.error,
	})
}

// Run runs this game against an endpoint located by url
// and writes states and actions to w. Timeout is set to
// t milliseconds, 0 for unlimited.
func (g *Game) Run(url string, t int, w io.Writer) error {
	le := json.NewEncoder(w)

	post := func() (*http.Response, error) {
		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(g); err != nil {
			return nil, err
		}

		c := &http.Client{
			Timeout: time.Duration(t) * time.Millisecond,
		}
		return c.Post(url, "application/json", b)
	}

	if err := le.Encode(g); err != nil {
		return err
	}

	for g.outcome == gameOutcomePending {
		r, err := post()
		if err != nil {
			return err
		}

		a, err := decodeAction(r.Body)
		if err != nil {
			g.error = err.Error()
			continue
		}

		if err := r.Body.Close(); err != nil {
			return err
		}

		if err := le.Encode(a); err != nil {
			return err
		}

		if err := g.model.handleAction(a); err != nil {
			g.error = err.Error()
			continue
		}
		g.error = ""

		if err := le.Encode(g); err != nil {
			return err
		}
	}

	post()
	return nil
}

func defaultGame() *Game {
	return &Game{
		outcome: gameOutcomePending,
		round:   1,
		points:  40,
	}
}

// New returns a Game that is initialized
// with default values, cities and pathogens.
func New() *Game {
	g := defaultGame()
	g.initialize(defaultCities(), defaultPathogens())
	g.model = newModel(g)

	return g
}
