from ictsc2021 import Rikka

import time
from math import pow, sqrt


def main():
    rikka = Rikka(baseurl="http://localhost:8080/api")

    print(f"\x1b[33m\n*** signin\x1b[0m")
    rikka.signin("ictsc", "2ht4BN9q6tjc")

    print(f"\x1b[33m\n*** get self\x1b[0m")
    results = []
    for i in range(256):
        start = time.clock_gettime_ns(time.CLOCK_MONOTONIC)
        rikka.self()
        end = time.clock_gettime_ns(time.CLOCK_MONOTONIC)

        req_ms = (end - start) / 1000000
        results.append(req_ms)
        print(f"{i} - {req_ms}")

    mean = sum(results) / len(results)
    stddev = sqrt(sum([pow(res - mean, 2) for res in results]) / (len(results) - 1))

    print(f"mean = {mean}, stddev = {stddev}")


main()