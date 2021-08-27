from ictsc2021 import Rikka

import json
import logging
from http import client


def load_setting():
    filename = "./setting.prod.json"
    with open(filename, "r") as f:
        setting_json = f.read()
    return json.loads(setting_json)

def main():
    client.HTTPConnection.debuglevel = 1
    logging.basicConfig(level=logging.DEBUG)

    setting = load_setting()

    rikka = Rikka(baseurl="https://contest.ictsc.net/api")

    print(f"\x1b[33m\n*** signin\x1b[0m")
    rikka.signin("ictsc", "")

    for team in setting["teams"]:
        name = team["name"]
        organization = team["organization"]
        invitation_code = team["invitation_code"]

        print(f"\x1b[33m\n*** Create user group {name}\x1b[0m")
        rikka.create_usergroup(name, organization, invitation_code, False)


main()
