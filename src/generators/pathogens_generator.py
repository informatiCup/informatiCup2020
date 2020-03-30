#!/usr/bin/env python3

import csv
import random


def parse_float(s):
    return float(s.replace(".", "").replace(",", "."))


pathogens = []
with open("pathogens.csv") as csv_file:
    csv_reader = csv.reader(csv_file, delimiter="\t")
    for row in csv_reader:
        if not row:
            continue
        pathogens += [
            {
                "n": row[0],
                "i": parse_float(row[1]),
                "l": parse_float(row[2]),
                "m": parse_float(row[3]),
                "d": int(parse_float(row[4]))
            }
        ]

print("// GENERATED CODE. DO NOT EDIT!")
print()
print("package game")
print()
print("func defaultPathogens() []*pathogen {")
print("\treturn []*pathogen{")
for p in pathogens:
    print(f'\t\tnewPathogen("{p["n"]}", {round(p["i"], 3)}, {round(p["l"], 3)}, {round(p["m"], 3)}, {p["d"]}),')
print("\t}")
print("}")
