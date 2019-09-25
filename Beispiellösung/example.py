#!/usr/bin/env python3

from bottle import post, request, run, BaseRequest


@post("/")
def index():
    game = request.json
    print(f'round: {game["round"]}, outcome: {game["outcome"]}')
    return {"type": "endRound"}


BaseRequest.MEMFILE_MAX = 10 * 1024 * 1024
run(host="0.0.0.0", port=50123, quiet=True)
