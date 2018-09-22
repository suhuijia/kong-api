#coding=utf-8

import os
import base64
import json


url = "http://xxx:8000/v1/face/"
img_path = '/home/nd/work/2.jpg'

def sym_recog():
    import urllib.request
    import urllib.parse
    f = open(img_path, 'rb')
    img = base64.b64encode(f.read())
    img_string = img.decode('utf-8')

    params = {
        # 'img_url': "http://img.my.csdn.net/uploads/201212/25/1356422284_1112.jpg",
        'img_url': img_string,
        # 'top_n': 8
        # 'appid': "1257232266"
    }

    # params = urllib.parse.urlencode(params).encode(encoding='utf-8')
    # with open("test.txt", 'wb') as wf:
    #     wf.write(params)
    # headers = {'Content-Type': 'application/x-www-form-urlencoded'}
    headers = {'Content-Type': 'application/json'}

    request = urllib.request.Request(url=url, headers=headers, data=bytes(json.dumps(params), 'utf8'))  #bytes(json.dumps(params), 'utf8')

    response = urllib.request.urlopen(request)
    content = response.read().decode()
    if content:
        print(content)
        print(type(content))

sym_recog()
