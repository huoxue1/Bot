package daka

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type data struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Do() bool {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("打卡失败")
		}
	}()
	var data = map[string]string{
		"19104978":     "苟江山",
		"19106360":     "周玲",
		"19105101":     "遇溪涓",
		"19104950":     "陈峰",
		"19107611":     "王干",
		"19104916":     "张灿",
		"19104977":     "陈月皓",
		"19104958":     "陈伟",
		"19104965":     "周杨琪",
		"19101242":     "张志成",
		"19104671":     "郭立扬",
		"19208932":     "李宗杰",
		"201813015120": "李明宸",
		"19208581":     "潘鹏程",
		"201817025138": "杨新",
		"19104824":     "巫雨",
		"19104966":     "白义枭",
		"19106543":     "付焱青",
		"19106564":     "蒲延慧",
		"202542020058": "苏骏",
		"201830055117": "黎智超",
		"19104668":     "吴仲鑫"}
	hour := time.Now().Hour()
	fmt.Println(hour)
	for k, v := range data {
		err := commit(k, v)
		if err != nil {
			panic(err)
			return false

		}
	}
	return true
}
func commit(xh, xm string) error {
	//proxy := func(_ *http.Request) (*url.URL, error) {
	//	return url.Parse("http://127.0.0.1:8888")
	//}

	//transport := &http.Transport{Proxy: proxy}
	client := http.Client{}
	values := url.Values{}
	values.Set("xh", xh)
	values.Set("xm", xm)
	request, err := http.NewRequest(http.MethodPost, "http://xxcj.scnucas.com/xxcj/Login_ck.php", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("User-agent", "micromessenger")
	request.Header.Set(`Content-Type`, `application/x-www-form-urlencoded`)
	response, err := client.Do(request)
	WriteFile(xm + "  " + xh)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	d, err := commitData(response, client, 1)

	if d.Code == 0 {
		WriteFile("打早卡成功")
	} else if d.Code == 1 {
		WriteFile("今日早卡已打卡")
	}
	d, err = commitData(response, client, 2)

	if d.Code == 0 {
		WriteFile("打午卡成功")
	} else if d.Code == 1 {
		WriteFile("今日午卡已打卡")
	}
	d, err = commitData(response, client, 3)

	if d.Code == 0 {
		WriteFile("打晚卡成功")
	} else if d.Code == 1 {
		WriteFile("今日晚卡已打卡")
	}
	return err
}

func commitData(response *http.Response, client http.Client, num int) (data, error) {
	v := url.Values{}
	v.Set("post_date", fmt.Sprintf("day_tj1=A&day_tj2=A&day_tj2_zzxq=&day_tj4=A&day_tj4_1=&day_tj5=B&day_tj7=1&szdq_no=&szdq_vl=&jianyi=&lx=%v", num))
	newRequest, err := http.NewRequest(http.MethodPost, "http://xxcj.scnucas.com/xxcj/fx_action.php", strings.NewReader(v.Encode()))
	if err != nil {
		return data{}, err
	}
	newRequest.Header.Set("User-agent", "micromessenger")
	newRequest.Header.Set(`Content-Type`, `application/x-www-form-urlencoded`)
	for _, v := range response.Cookies() {
		newRequest.AddCookie(v)
	}
	do, err := client.Do(newRequest)
	if err != nil {
		return data{}, err
	}
	defer do.Body.Close()
	all, err := ioutil.ReadAll(do.Body)
	var res data
	err = json.Unmarshal(all, &res)
	if err != nil {
		fmt.Println(err)
		return data{}, err
	}
	return res, err
}

func WriteFile(string2 string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println(string2)
	path, err := os.Getwd()
	file, err := os.OpenFile(path+"/logs/"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	_, err = file.WriteString(time.Now().Format("2006-01-02 15:04:05") + " " + string2 + "\n\n")
	if err != nil {
		panic("log  " + string2 + "  error")
	}
}
