#!/usr/bin/env python3

import glob
import re
import json
import sys

root = sys.argv[1] if len(sys.argv) > 1 else "."
paths = sorted(glob.glob(f"{root}/**/*.json"))
for path in paths:
    match = re.findall("(\d+)-games/(\d+).json", path)
    if not match:
        continue
    group, seed = int(match[0][0]), int(match[0][1])
    with open(path, "r") as log_file:
        lines = log_file.readlines()
    last_state = json.loads(lines[-1])
    rounds, outcome = last_state["round"], last_state["outcome"]
    print(seed, "\t", group, "\t", outcome, "\t", rounds)
