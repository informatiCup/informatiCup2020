#!/usr/bin/env python3

from bottle import post, request, run, BaseRequest
from math import floor, ceil
from random import seed, shuffle, choice
from os import environ


def get_city_events_of_type(city, type):
    return [e for e in city.get("events") or [] if e["type"] == type]


def get_game_events_of_type(game, type):
    return [e for e in game.get("events") or [] if e["type"] == type]


def get_cities_events_of_type(game, type):
    return [(c, e) for c in game["cities"].values() for e in get_city_events_of_type(c, type)]


def get_cities_ordered_by_population_descendingly(game):
    tuples = [(c, c["population"]) for c in game["cities"].values()]
    tuples.sort(key=lambda ci: ci[1], reverse=True)
    return [ci[0] for ci in tuples]


def get_encountered_pathogens(game, lethalties):
    return [
        e["pathogen"]["name"]
        for e in get_game_events_of_type(game, "pathogenEncountered") if e["pathogen"]["lethality"] in lethalties
    ]


def get_pathogens_with_available_vaccines(game):
    return [e["pathogen"]["name"] for e in get_game_events_of_type(game, "vaccineAvailable")]


def get_pathogens_with_available_medications(game):
    return [e["pathogen"]["name"] for e in get_game_events_of_type(game, "medicationAvailable")]


def put_under_quarantine_heuristic(game):
    # Put under quarantine the city with the highest population and an outbreak of the most lethal
    # pathogen for up to 10 rounds. Ignore outbreaks of pathogens with lethality lower than "o".
    outbreaks = get_cities_events_of_type(game, "outbreak")
    outbreaks = [ce for ce in outbreaks if ce[1]["pathogen"]["lethality"] in ["o", "+", "++"]]
    if not outbreaks:
        return None
    outbreaks.sort(key=lambda ce: (-ce[0]["population"], ce[1]["pathogen"]["lethality"]))
    city = outbreaks[0][0]
    if get_city_events_of_type(city, "quarantine"):
        return None
    rounds = floor((game["points"] - 20) / 10)
    if rounds > 10:
        rounds = 10
    if rounds > 0:
        return {"type": "putUnderQuarantine", "city": city["name"], "rounds": rounds}
    return None


def close_airport_heuristic(game):
    # Close the airport of the city with the highest population and an outbreak
    # of the most lethal pathogen for up to 10 rounds. Ignore outbreaks of pathogens with
    # lethality lower than "o".
    outbreaks = get_cities_events_of_type(game, "outbreak")
    outbreaks = [ce for ce in outbreaks if ce[1]["pathogen"]["lethality"] in ["+", "++", "o"]]
    if not outbreaks:
        return None
    outbreaks.sort(key=lambda ce: (-ce[0]["population"], ce[1]["pathogen"]["lethality"]))
    city = outbreaks[0][0]
    if get_city_events_of_type(city, "airportClosed"):
        return None
    rounds = floor((game["points"] - 15) / 5)
    if rounds > 10:
        rounds = 10
    if rounds > 0:
        return {"type": "closeAirport", "city": city["name"], "rounds": rounds}
    return None


def develop_vaccine_heuristic(game):
    # Develop a vaccine against an encountered pathogen with lethality greater than "o".
    if game["points"] <= 40:
        return None
    encountered_pathogens = get_encountered_pathogens(game, ["+", "++"])
    if not encountered_pathogens:
        return None
    shuffle(encountered_pathogens)
    pathogens_with_vaccines_in_development = [
        e["pathogen"]["name"] for e in get_game_events_of_type(game, "vaccineInDevelopment")
    ]
    pathogens_with_available_vaccines = get_pathogens_with_available_vaccines(game)
    for pathogen in encountered_pathogens:
        if not pathogen in pathogens_with_vaccines_in_development and not pathogen in pathogens_with_available_vaccines:
            return {"type": "developVaccine", "pathogen": pathogen}
    return None


def deploy_vaccine_heuristic(game):
    # Deploy an available vaccine in the city with the highest population.
    pathogens_with_available_vaccines = get_pathogens_with_available_vaccines(game)
    if not pathogens_with_available_vaccines:
        return None
    for city in get_cities_ordered_by_population_descendingly(game):
        if game["points"] < 5:
            continue
        pathogens_with_deployed_vaccines = [
            e["pathogen"]["name"] for e in get_city_events_of_type(city, "vaccineDeployed")
        ]
        for pathogen in pathogens_with_available_vaccines:
            if not pathogen in pathogens_with_deployed_vaccines:
                return {"type": "deployVaccine", "city": city["name"], "pathogen": pathogen}
    return None


def develop_medication_heuristic(game):
    # Develop a medication against an encountered pathogen.
    if game["points"] < 20:
        return None
    encountered_pathogens = get_encountered_pathogens(game, ["--", "-", "o", "+", "++"])
    if not encountered_pathogens:
        return None
    shuffle(encountered_pathogens)
    pathogens_with_medications_in_development = [
        e["pathogen"]["name"] for e in get_game_events_of_type(game, "medicationInDevelopment")
    ]
    pathogens_with_available_medications = get_pathogens_with_available_medications(game)
    for pathogen in encountered_pathogens:
        if not pathogen in pathogens_with_medications_in_development and not pathogen in pathogens_with_available_medications:
            return {"type": "developMedication", "pathogen": pathogen}
    return None


def deploy_medication_heuristic(game):
    # Deploy an available medication in the city with the highest population.
    pathogens_with_available_medications = get_pathogens_with_available_medications(game)
    if not pathogens_with_available_medications:
        return None
    for city in get_cities_ordered_by_population_descendingly(game):
        if game["points"] < 10:
            continue
        pathogens_with_outbreaks = [e["pathogen"]["name"] for e in get_city_events_of_type(city, "outbreak")]
        for pathogen in pathogens_with_available_medications:
            if pathogen in pathogens_with_outbreaks:
                return {"type": "deployMedication", "city": city["name"], "pathogen": pathogen}
    return None


def use_random_property_modifier_heuristic(game):
    if game["points"] < 3:
        return
    type = choice(["applyHygienicMeasures", "callElections", "excertInfluence", "launchCampaign"])
    outbreaks = get_cities_events_of_type(game, "outbreak")
    if not outbreaks:
        return
    city = choice(outbreaks)[0]
    return {"type": type, "city": city["name"]}


def end_round_heuristic(game):
    # End the round.
    return {"type": "endRound"}


def invalid_action_heuristic(game):
    return {"type": "invalid"}


def get_action(game):
    heuristics = []
    heuristics += [put_under_quarantine_heuristic] * 10
    heuristics += [close_airport_heuristic] * 2
    heuristics += [develop_vaccine_heuristic] * 3
    heuristics += [deploy_vaccine_heuristic] * 3
    heuristics += [develop_medication_heuristic] * 1
    heuristics += [deploy_medication_heuristic] * 1
    heuristics += [use_random_property_modifier_heuristic] * 1
    heuristics += [end_round_heuristic] * 2
    heuristics += [invalid_action_heuristic] * 1
    shuffle(heuristics)
    for heuristic in heuristics:
        action = heuristic(game)
        if action is not None:
            return action
    return None


def format_event(event):
    return ", ".join([f'{key}: {value}' for key, value in event.items()])


def print_game(game, verbosity=0):
    if verbosity > 0:
        print(
            f'round: {game["round"]}, points: {game["points"]}, events: {len(game.get("events") or [])}, error: {game.get("error") or "?"}'
        )
    if verbosity > 1:
        for event in game["events"]:
            print(">", format_event(event))
    for city_name, city in game["cities"].items():
        if verbosity > 1:
            properties = ", ".join(
                [f'{k}: {city[k]}' for k in ["population", "awareness", "economy", "government", "hygiene"]]
            )
            print(f'  {city_name} - {properties}, events: {len(city.get("events") or [])}')
        if not "events" in city:
            continue
        if verbosity > 1:
            for event in city["events"]:
                print("  >", format_event(event))


@post("/")
def index():
    verbosity = int(environ.get("VERBOSITY") or 0)
    game = request.json
    if game["outcome"] != "pending":
        if verbosity > 0:
            print(f'outcome: {game["outcome"]}')
        return ""
    print_game(game, verbosity=verbosity)
    action = get_action(game)
    if verbosity > 0:
        print(f'performing: {action}')
    return action


seed(0)
BaseRequest.MEMFILE_MAX = 10 * 1024 * 1024
run(port=50123, quiet=True)
