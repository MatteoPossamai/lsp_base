import json
import requests

def rpc_call(url, method, args):
    headers = {'content-type': 'application/json'}
    payload = {
        "method": method,
        "params": [args],
        "jsonrpc": "2.0",
        "id": 1,
    }
    response = requests.post(url, data=json.dumps(payload), headers=headers).json()
    return response

url = 'http://localhost:5000/'

emailArgs = {'To': 'demo@example.com','Subject': 'Hello', 'Content': 'Hi!!!'}
print(rpc_call(url, 'email.SendEmail', emailArgs))