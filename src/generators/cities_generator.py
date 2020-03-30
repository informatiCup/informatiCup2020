#!/usr/bin/env python3

import csv
import random


def random_bool(t):
    return random.random() < t


def parse_float(s):
    return float(s.replace(".", "").replace(",", "."))


def value_from_base(b):
    v = random.gauss(b, 0.1)
    if v < 0:
        v = 0
    if v > 1:
        v = 1
    return v


cities = {}
with open("cities.csv") as csv_file:
    csv_reader = csv.reader(csv_file, delimiter="\t")
    for row in csv_reader:
        if not row:
            continue
        name = row[0].strip()
        base = parse_float(row[4])
        cities[name] = {
            "n": name,
            "la": parse_float(row[1]),
            "lo": parse_float(row[2]),
            "p": int(parse_float(row[3])),
            "e": value_from_base(base),
            "g": value_from_base(base),
            "h": value_from_base(base),
            "a": value_from_base(base),
            "c": set()
        }


def num_connections(c):
    return len(c["c"])


def maybe_connect(c1, c2):
    if c1 == c2:
        return
    c1["c"].add(c2["n"])
    c2["c"].add(c1["n"])


huge_cities = [c for c in cities.values() if c["p"] > 5000]
huge_cities.sort(key=lambda c: c["p"], reverse=True)
large_cities = [c for c in cities.values() if 1000 < c["p"] <= 5000]
large_cities.sort(key=lambda c: c["p"], reverse=True)
medium_cities = [c for c in cities.values() if 100 < c["p"] <= 1000]
medium_cities.sort(key=lambda c: c["p"], reverse=True)
small_cities = [c for c in cities.values() if c["p"] <= 100]
small_cities.sort(key=lambda c: c["p"], reverse=True)

for c in huge_cities:
    while num_connections(c) < 4:
        maybe_connect(c, random.choice(huge_cities))
    while num_connections(c) < 6:
        maybe_connect(c, random.choice(large_cities))
    while num_connections(c) < 8:
        maybe_connect(c, random.choice(medium_cities))
for c in large_cities:
    while num_connections(c) < 3:
        maybe_connect(c, random.choice(large_cities))
    while num_connections(c) < 5:
        maybe_connect(c, random.choice(medium_cities))
for c in medium_cities:
    while num_connections(c) < 2:
        maybe_connect(c, random.choice(medium_cities))

print("// GENERATED CODE. DO NOT EDIT!")
print()
print("package game")
print()
print("func defaultCities() []*city {")
print("\treturn []*city{")
for c in cities.values():
    cs = ", " + ", ".join([f'"{n}"' for n in c["c"]]) if c["c"] else ""
    print(
        f'\t\tnewCity("{c["n"]}", {c["la"]}, {c["lo"]}, {c["p"]}, {round(c["e"], 3)}, {round(c["g"], 3)}, {round(c["h"], 3)}, {round(c["a"], 3)}{cs}),'
    )
print("\t}")
print("}")
