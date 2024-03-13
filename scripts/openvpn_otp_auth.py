# encoding: utf-8
#  Copyright © 2024 MicroOps-cn.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
import base64
import logging
import sys, requests

# change it
AUTH_SERVER = "http://127.0.0.1:8081/api/v1/oauth/token"
CLIENT_ID = "HkdQYrf8NNaAqmTG7S5BThJAAslstKCJ1uNuRymJs7M"
CLIENT_SECRET = "vYRssEjBxXZQf95hdFH9Iy6tgUi75RBRnZdiVR46TFM"

def read_auth_file_to_auth_body(path):
    """
    Read authentication information from the file.
    :type path: str
    :rtype dict[str,str]
    """
    auth = {}
    with open(path) as f:
        username = f.readline().strip()
        passwords = f.readline().strip().split(":")
        if len(passwords)>=3:
            auth.update({"username":username,"code": base64.b64decode(passwords[2]).decode(),"password": base64.b64decode(passwords[1]).decode()})
        else:
            auth.update({"username":username,"password": base64.b64decode(passwords[1]).decode()})
    return auth

def auth_from_oauth2_password(auth_payload):
    """
    Use OAuth2.0 password mode for authentication.
    :type auth_payload: dict[str,str]
    """
    payload = {
        "client_id": CLIENT_ID,
        "client_secret": CLIENT_SECRET,
        "grant_type": "password",
        **auth_payload,
    }
    headers = {
        'User-Agent': 'OpenVPN Auth Client;',
        'Content-Type': 'application/json'
    }
    response = requests.request("POST", AUTH_SERVER, headers=headers, json=payload)

    r = response.json()
    if r.get("success",True) and r.get("access_token",None):
        return
    else:
        logging.error("auth failed: errorCode={}, errorMessage={}".format(r["errorCode"],r["errorMessage"]))
        sys.exit(1)

if __name__ == '__main__':
    if len(sys.argv) > 2:
        sys.stderr.write("Only one argument expected - credential file.")
        exit(1)
    if len(sys.argv) < 2:
        sys.stderr.write("Need one argument - credential file.")
        exit(1)
    body = read_auth_file_to_auth_body(sys.argv[1])
    auth_from_oauth2_password(body)