import requests
import json
import time

task_num = 0

begin_t = 0
stat_cnt = 900
cur_cnt = 0
while True:

    is_proxy = False
    srv_url = "http://XXX:3611"



    url = srv_url +'/get_task/'

    if is_proxy:
        r = requests.get(url, proxies=dict(http='http:192.168.CCC:1080'))
    else:
        r = requests.get(url)

    # print(r)
    # r = requests.get(url)

    task_obj = json.loads(r.content)
    print(task_obj)
    code = task_obj["code"]
    msg = task_obj["msg"]

    if code == 2000:
        time.sleep(1)
        print "No task."
        continue

    if code == 0:
        # print "xxx"
        client_id = task_obj["data"]["client_id"]
        task_uuid = task_obj["data"]["task_uuid"]
        # pt_seq = task_obj["data"]["pt_seq"]
        img_url = task_obj["data"]["img_url"]
        # pt_seq = task_obj["data"]["latex_gt"]
        print task_uuid
        print client_id
        print img_url
        # print pt_seq

        recog_result = {}
        recog_result["error"] = 0
        # recog_result["activation_idx"] = [14, 24, 42, 50, 66]
        # recog_result["labels"] = [["A", "B", "C", "C"], ["C","D","F"]]
        # recog_result["latex"] = "\\frac{2}{8}"
        # recog_result["correct"] = True
        recog_result["result_img_url"] = "http://cs.101.com/v0.1/download?dentryId=b6d54f48-731b-4922-bba3-9ef7f3290164"
        recog_result["metarial_id"] = "123465465"

        data = json.dumps({"task_uuid": task_uuid, "result": json.dumps(recog_result)})
        # files = {"upl_file": open(file_name, "rb")}
        headers = {'Content-Type': 'application/json'}

        if is_proxy:
            result = requests.post(url= srv_url + '/upload_result/', headers=headers, data=data, proxies=dict(http='http:192.168.46.116:1080'))
        else:
            result = requests.post(url=srv_url + '/upload_result/', headers=headers, data=data)

        print(result)

    # except Exception, e:
    #     pass
        # print type(e), e.message
