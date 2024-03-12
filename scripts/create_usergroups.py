import json
import logging
import os
from http import client

from ictsc2021 import Rikka


def load_setting():
    filename = "./setting.prod.json"
    with open(filename, "r") as f:
        setting_json = f.read()
    return json.loads(setting_json)


def main():
    client.HTTPConnection.debuglevel = 1
    logging.basicConfig(level=logging.DEBUG)

    setting = load_setting()

    rikka = Rikka(baseurl=str(os.environ.get("baseurl")))

    print("\x1b[33m\n*** signin\x1b[0m")  # Add a placeholder to the f-string
    rikka.signin(os.environ.get("username"), os.environ.get("password"))

    for team in setting["teams"]:
        name = team["name"]
        organization = team["organization"]
        invitation_code = team["invitation_code"]
        bastion_user = ""
        bastion_password = ""
        bastion_host = ""
        bastion_port = 0
        team_id = team["team_name"]

        print(f"\x1b[33m\n*** Create user group {name}\x1b[0m")
        rikka.create_usergroup(
            name,
            organization,
            invitation_code,
            False,
            bastion_user,
            bastion_password,
            bastion_host,
            bastion_port,
            team_id,
        )


main()
