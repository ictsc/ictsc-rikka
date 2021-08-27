from ictsc2021 import Rikka

import logging
import time
from math import pow, sqrt

import random, string

def randomname(n):
   randlst = [random.choice(string.ascii_letters + string.digits) for i in range(n)]
   return ''.join(randlst)

def main():
    logging.basicConfig(level=logging.DEBUG)

    rikka = Rikka(baseurl="https://contest.ictsc.net/api")
    rikka2 = Rikka(baseurl="https://contest.ictsc.net/api")

    print(f"\x1b[33m\n*** signin\x1b[0m")
    rikka.signin("team92", "ictsc2021")
    rikka2.signin("ictsc", "")

    print(f"\x1b[33m\n*** send answer\x1b[0m")
    results = []
    problem_id = "0060002b-8898-4035-bf3a-f313426cd013"
    for i in range(256):
        start = time.clock_gettime_ns(time.CLOCK_MONOTONIC)
        answer = rikka.send_answer(problem_id, randomname(1024)).json()
        answer_id = answer["data"]["answer"]["id"]
        rikka2.point(problem_id, answer_id, 100)
        end = time.clock_gettime_ns(time.CLOCK_MONOTONIC)

        req_ms = (end - start) / 1000000
        results.append(req_ms)
        print(f"{i} - {req_ms}")

    mean = sum(results) / len(results)
    stddev = sqrt(sum([pow(res - mean, 2) for res in results]) / (len(results) - 1))

    print(f"mean = {mean}, stddev = {stddev}")


main()
