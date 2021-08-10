from ictsc2021 import Rikka

import logging
from http import client


def main():
    client.HTTPConnection.debuglevel = 1
    logging.basicConfig(level=logging.DEBUG)

    rikka = Rikka(baseurl="https://ss.ictsc.net/api")

    print(f"\x1b[33m\n*** signin\x1b[0m")
    rikka.signin("ictsc", "")

    print(f"\x1b[33m\n*** Create user group\x1b[0m")
    resp = rikka.create_usergroup("team90", "team90", "ictsc2021team90hotstage", False)
    data = resp.json()
    print(data)

    print(f"\x1b[33m\n*** Create user group\x1b[0m")
    resp = rikka.create_usergroup("team99", "team99", "ictsc2021team99hotstage", False)
    data = resp.json()
    print(data)



main()
