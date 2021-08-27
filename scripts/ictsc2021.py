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

    def _patch(self, endpoint, **kwargs):
        url = self._get_url(endpoint)
        return requests.patch(url, cookies=self._jar, **kwargs)

    def signin(self, username, password):
        resp = self._post("/auth/signin", json={
            "name": username,
            "password": password,
        })
        self._jar = resp.cookies
        return resp

    def self(self):
        return self._get("/auth/self")

    def create_user(self, name, password, usergroupid, invitation_code):
        return self._post("/users", json={
            "name": name,
            "password": password,
            "user_group_id": usergroupid,
            "invitation_code": invitation_code,
        })

    def create_usergroup(self, name, organization, invitation_code, is_full_access):
        return self._post("/usergroups", json={
            "name": name,
            "organization": organization,
            "invitation_code": invitation_code,
            "is_full_access": is_full_access,
        })

    def send_answer(self, problem_id, content):
        return self._post("/problems/%s/answers" % problem_id, json={
            "body": content
        })

    def point(self, problem_id, answer_id, point):
        return self._patch("/problems/%s/answers/%s" % (problem_id, answer_id), json={
            "point": point,
        })
