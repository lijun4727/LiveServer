#!python3
import requests
import ssl

ssl._create_default_https_context = ssl._create_unverified_context

url = 'http://127.0.0.1:8083/account/login'
params = {"phone":"13651884967","device_id":"tv888"}
r = requests.post(url, json=params).json()
print(r)