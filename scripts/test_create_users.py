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

    for team in setting["teams"]:
        name = team["team_id"]
        user_group_id = team["user_group_id"]
        invitation_code = team["invitation_code"]

        print(f"\x1b[33m\n*** [TEST] create user {name}\x1b[0m")
        rikka.create_user(name, "ictsc2021", user_group_id, invitation_code)


main()