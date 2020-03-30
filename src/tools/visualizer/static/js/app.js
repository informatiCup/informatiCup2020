Vue.component("icv-games", {
  template: "#icv-games-template",
  methods: {
    async onDocumentKeyDown(event) {
      if (this.busy) {
        return
      }
      switch (event.key) {
        case "ArrowLeft":
          this.shiftIndex(-1)
          break
        case "ArrowRight":
          this.shiftIndex(1)
          break
      }
    },
    async queryGame(name) {
      let response = await fetch(`/game/${name}`)
      return await response.json()
    },
    shiftIndex(delta) {
      let indexShifted = false
      let atStart = false
      let atEnd = true
      for (let game of this.games) {
        let nextIndex = game.index + delta
        let nextState = game.states[nextIndex]
        if (nextState && nextState.round === this.round) {
          game.index += delta
          indexShifted = true
        }
        if (nextIndex === -1) {
          atStart = true
        }
        if (nextIndex < game.states.length) {
          atEnd = false
        }
      }
      if (!indexShifted && !atStart && !atEnd) {
        this.round += delta
        this.shiftIndex(delta)
      }
    },
    onPreviousStateClick() {
      this.shiftIndex(-1)
    },
    onNextStateClick() {
      this.shiftIndex(1)
    }
  },
  async created() {
    for (let i = 0; i < 4; i++) {
      this.games.push({
        name: "" + i,
        busy: true,
        index: 0,
        states: []
      })
    }
  },
  async mounted() {
    document.addEventListener("keydown", this.onDocumentKeyDown)

    for (let game of this.games) {
      this.status = `Querying game '${game.name}'...`
      game.states = await this.queryGame(game.name)
    }
    this.busy = false
  },
  data() {
    return {
      busy: true,
      status: "Initializing...",
      games: [],
      round: 1
    }
  }
})

Vue.component("icv-game", {
  template: "#icv-game-template",
  props: {
    states: { type: Array, default: () => [] },
    index: { type: Number, default: 0 },
    busy: { type: Boolean, default: false }
  },
  watch: {
    states() {
      this.initialPopulation = this.population
      this.initialPopulations = this.cities.reduce((r, c) => {
        r[c.name] = c.population
        return r
      }, {})
    }
  },
  computed: {
    state() {
      return this.states[this.index] || this.states[this.states.length - 1] || {}
    },
    round() {
      return this.state.round
    },
    cities() {
      return Object.values(this.state.cities || {})
    },
    population() {
      return this.cities.reduce((p, c) => p + c.population, 0)
    },
    cityIcons() {
      let icons = []
      for (let city of this.cities) {
        let iconAdded = false
        let addIcon = properties => {
          let cartesian = this.latitudeLongitudeToCartesian(city.latitude, city.longitude)
          let style = Object.assign(properties.style, { left: `${cartesian.x}px`, top: `${cartesian.y}px` })
          icons.push(Object.assign({ style: style }, properties))
          iconAdded = true
        }

        // High priority icons
        for (let event of city.events || []) {
          let round = this.state.round
          if (event.round && event.round !== round) {
            continue
          }
          switch (event.type) {
            case "airportClosed":
            case "bioTerrorism":
            case "connectionClosed":
            case "campaignLaunched":
            case "electionsCalled":
            case "influenceExerted":
            case "hygienicMeasuresApplied":
            case "quarantine":
              addIcon({
                name: event.type,
                style: { width: "32px" },
                classes: ["city", "icon", "tint"],
                title: this.cityIconTitle(city, event)
              })
              break
          }
        }
        if (iconAdded) {
          continue
        }

        // Low priority icons
        for (let event of city.events || []) {
          let round = this.state.round
          if (event.round && event.round !== round) {
            continue
          }
          switch (event.type) {
            case "outbreak":
              addIcon({
                name: `outbreak${event.pathogen.lethality}`,
                style: { width: `${this.outbreakIconWidth(city, event)}px` },
                classes: ["city", "icon"],
                title: this.cityIconTitle(city, event)
              })
              break
          }
        }
      }

      return icons
    },
    gameIcons() {
      let icons = []
      for (let event of this.state.events || []) {
        switch (event.type) {
          case "economicCrisis":
          case "largeScalePanic":
          case "medicationAvailable":
          case "medicationInDevelopment":
          case "pathogenEncountered":
          case "vaccineAvailable":
          case "vaccineInDevelopment":
            icons.push({
              name: event.type,
              text: event.round ? `${event.round}` : `${event.sinceRound}-${event.untilRound || ""}`,
              title: JSON.stringify(event),
              style: { width: "32px" }
            })
            break
        }
      }
      return icons
    },
    containerStyle() {
      return {
        "padding-bottom": `${(this.mapHeight / this.mapWidth) * 100}%`
      }
    }
  },
  methods: {
    latitudeLongitudeToCartesian(latitude, longitude) {
      // Miller projection
      let mX = longitude => longitude
      let mY = latitude => 1.25 * Math.log(Math.tan(Math.PI / 4 + 0.4 * latitude))

      let radian = Math.PI / 180
      return {
        x:
          ((mX(longitude * radian) - mX(this.longitudeLeft * radian)) * this.containerWidth) /
          (mX(this.longitudeRight * radian) - mX(this.longitudeLeft * radian)),
        y:
          ((mY(this.latitudeTop * radian) - mY(latitude * radian)) * this.containerHeight) /
          (mY(this.latitudeTop * radian) - mY(this.latitudeBottom * radian))
      }
    },
    cityIconTitle(city, event) {
      return city.name + ": " + JSON.stringify(event)
    },
    outbreakIconWidth(city, event) {
      let infected = event.prevalence * city.population
      let width = infected / 1000
      if (width < 6) {
        return 6
      }
      if (width > 16) {
        return 16
      }
      return width
    },
    queryContainerDimensions() {
      let element = this.$refs.container
      this.containerWidth = element.clientWidth
      this.containerHeight = element.clientHeight
    }
  },
  async mounted() {
    window.addEventListener("resize", this.queryContainerDimensions)
    this.queryContainerDimensions()
  },
  data() {
    return {
      alignment: false,
      mapWidth: 1670,
      mapHeight: 939,
      latitudeTop: 88,
      latitudeBottom: -64,
      longitudeLeft: -168,
      longitudeRight: 191,
      containerWidth: 0,
      containerHeight: 0,
      initialPopulation: 0,
      initialPopulations: {}
    }
  }
})

new Vue({
  el: "#app"
})
