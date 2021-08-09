from http.cookiejar import CookieJar

import requests


class Rikka:
    def __init__(self, baseurl):
        self._jar = CookieJar()
        self._baseurl = baseurl

    def _get_url(self, endpoint):
        return self._baseurl + endpoint

    def _get(self, endpoint, **kwargs):
        url = self._get_url(endpoint)
        return requests.get(url, cookies=self._jar, **kwargs)

    def _post(self, endpoint, **kwargs):
        url = self._get_url(endpoint)
        return requests.post(url, cookies=self._jar, **kwargs)

    def signin(self, username, password):
        resp = self._post("/auth/signin", json={
            "name": username,
            "password": password,
        })
        self._jar = resp.cookies

    def self(self):
        return self._get("/auth/self")

    def create_usergroup(self, name, organization, invitation_code, is_full_access):
        return self._post("/usergroups", json={
            "name": name,
            "organization": organization,
            "invitation_code": invitation_code,
            "is_full_access": is_full_access,
        })
