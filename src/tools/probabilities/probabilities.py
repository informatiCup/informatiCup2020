#!/usr/bin/env python3

import math
import random
import inspect


def inhabitant_is_infected_during_breakout(*, infectivity, economy, government, hygiene):
    return infectivity * (1 - economy / 20 - government / 20 - hygiene / 10)


def inhabitant_infects_other(*, infectivity, awareness, government, hygiene):
    return infectivity * (1 - awareness / 10 - government / 20 - hygiene / 10)


def pathogen_kills_inhabitant(*, lethality, economy, government, hygiene):
    return lethality * (1 - economy / 6 - government / 9 - hygiene / 3)


def pathogen_spreads_to_any(*, mobility, from_economy, to_economy, distance):
    return mobility * (1 - from_economy / 20 - to_economy / 20 - math.pow(distance, 1 / 4) / 1.125)


def pathogen_spreads_to_connected(
    *, infectivity, from_economy, to_economy, from_government, to_government, to_hygiene, to_awareness
):
    return infectivity * (
        1 - from_economy / 10 - to_economy / 5 - from_government / 10 - to_government / 5 - to_hygiene / 5 -
        to_awareness / 5
    )


def p(fn, **kargs):
    args = ", ".join([f'{k} = {v}' for k, v in kargs.items()])
    print(f'{fn.__name__}({args}) = {round(fn(**kargs), 6)}')


def s(fn, n, d={}):
    for _ in range(n):
        kargs = {n: round(random.random(), 3) for n in inspect.signature(fn).parameters.keys()}
        kargs.update(d)
        p(fn, **kargs)


s(pathogen_spreads_to_connected, 20, {"infectivity": 0.5})
