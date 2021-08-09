from ictsc2021 import Rikka

import logging
from http import client


def main():
    client.HTTPConnection.debuglevel = 1
    logging.basicConfig(level=logging.DEBUG)

    rikka = Rikka(baseurl="https://ss.ictsc.net/api")

    ugn = input("user_group_name: ")
    ugo = input("user_group_organization: ")
    ugit = input("user_group_invitation_token: ")
    un = input("user_name: ")
    up = input("user_password: ")

    print(f"\x1b[33m\n*** signin\x1b[0m")
    rikka.signin("ictsc", "2ht4BN9q6tjc")


    print(f"\x1b[33m\n*** Create user group\x1b[0m")
    resp = rikka.create_usergroup(ugn, ugo, ugit, False)
    data = resp.json()

    print(f"\x1b[33m\n*** Create user\x1b[0m")
    ugi = data["data"]["user_group"]["id"]

    rikka.create_user(un, up, ugi, ugit)


main()