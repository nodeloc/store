package services

import (
	"strings"
	"testing"
)

func splitNonEmpty(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

func TestBatchLooksLikeCode(t *testing.T) {
	cases := []struct {
		name string
		text string
		want bool
	}{
		{"python source", `
import requests
import json
import time

def change_step(user, pwd, step):
    url = "https://api.example.com/login"
    headers = {
        "User-Agent": "Mozilla/5.0",
    }
    data = {"user": user, "password": pwd}
    r = requests.post(url, data=data, headers=headers)
    if r.status_code == 200:
        print("ok")
    return r.json()

class Runner:
    pass
`, true},
		{"markdown tutorial with code", `
# 使用说明
1. 安装 python3
2. 运行以下命令：
pip install requests
python main.py --user 123 --step 20000
// 注意替换账号
for i in range(10):
    print(i)
`, true},
		{"plain card keys", `
ABCD-EFGH-IJKL-0001
ABCD-EFGH-IJKL-0002
ABCD-EFGH-IJKL-0003
ABCD-EFGH-IJKL-0004
ABCD-EFGH-IJKL-0005
ABCD-EFGH-IJKL-0006
`, false},
		{"account----password pairs", `
user1@example.com----P@ssw0rd!1
user2@example.com----P@ssw0rd!2
user3@example.com----P@ssw0rd!3
user4@example.com----P@ssw0rd!4
user5@example.com----P@ssw0rd!5
`, false},
		{"netdisk links", `
https://pan.baidu.com/s/1AbCdEf 提取码 x1y2
https://pan.baidu.com/s/1GhIjKl 提取码 a3b4
https://www.aliyundrive.com/s/MnOpQr
https://www.aliyundrive.com/s/StUvWx
https://example.com/download?token=abc123&id=9
`, false},
		{"too few lines to judge", "def a():\n    pass\n", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := batchLooksLikeCode(splitNonEmpty(tc.text)); got != tc.want {
				t.Fatalf("batchLooksLikeCode = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParseCardLine(t *testing.T) {
	no, pwd := parseCardLine("acc----pw----x")
	if no != "acc" || pwd != "pw----x" {
		t.Fatalf("got %q %q", no, pwd)
	}
	no, pwd = parseCardLine("ONLYKEY")
	if no != "ONLYKEY" || pwd != "" {
		t.Fatalf("got %q %q", no, pwd)
	}
}
