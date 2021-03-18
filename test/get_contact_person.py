#!python3
import requests
import ssl

ssl._create_default_https_context = ssl._create_unverified_context

url = 'http://127.0.0.1:8083/getContactPersons'
params = {
    "phone": "13651884967",
    "token": 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6IjEzNjUxODg0OTY3dHY4ODgiLCJpcCI6IiIsImV4cCI6MTYxNDU2OTg3Nn0.QOfS6GUFFl8aQaam75tLwBbvKz88Wp5KY1jCGka96Cg'
}
r = requests.post(url, json=params).json()
print(r)
