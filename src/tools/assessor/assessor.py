#!/usr/bin/env python3

import json
import subprocess
import sys
import os
import shutil

if len(sys.argv) != 2:
    print("Usage: assessor.py <settings file path>")
    sys.exit(1)

# Read and print settings
print("Settings")
with open(sys.argv[1], "r") as f:
    settings = json.load(f)
print("\n".join([f"- {k}: {v}" for k, v in settings.items()]))
print()

# Read seeds
seeds = []
with open(settings["seedsPath"], "r") as f:
    for string_seed in f:
        string_seed = string_seed.rstrip("\n")
        try:
            seed = int(string_seed)
            seeds += [seed]
        except ValueError:
            print(f"Discarding invalid seed: '{string_seed}'")

# Determine results
results = []
wins = 0
win_rounds = 0
losses = 0
loss_rounds = 0
if not os.path.exists(settings["gamesPath"]):
    os.mkdir(settings["gamesPath"])
for seed in seeds:
    try:
        print(f"{seed: <25}", end="", flush=True)
        try:
            output_path = f'{settings["gamesPath"]}/{seed}.json'
            if os.path.exists(output_path):
                print("skipped")
                continue
            process = subprocess.run(
                [
                    settings["binaryPath"], "-s",
                    str(seed), "-u", settings["endpointURL"], "-t",
                    str(settings["timeout"]), "-o", output_path
                ]
            )
            lines = tuple(open(output_path, "r"))
        except Exception as exception:
            print("failed")
            os.remove(output_path)
            print(f"Failed to run command line tool: {exception}")
            break
        state = json.loads(lines[-1])
        rounds = state["round"]
        outcome = state["outcome"]
        result = {"rounds": rounds, "outcome": outcome, "seed": seed}
        results += [result]
        if outcome == "win":
            wins += 1
            win_rounds += rounds
        if outcome == "loss":
            losses += 1
            loss_rounds += rounds
        print(f'{outcome} ({rounds})')
    except KeyboardInterrupt:
        print("interrupted")
        os.remove(output_path)
        break
print(f"{wins} ({win_rounds}) - {losses} ({loss_rounds})")

# Store results
with open(settings["resultsPath"], "w") as f:
    json.dump(results, f)
