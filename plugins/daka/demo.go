package daka

import (
	"encoding/json"
	"fmt"
	go_mybots "github.com/3343780376/go-mybots"
	"github.com/robfig/cron"
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
type date struct {
	Xm string `json:"xm"`
	Xh string `json:"xh"`
}

var bot = go_mybots.Bots{Address: "127.0.0.1", Port: 5700, Admin: 3343780376}

func Cr() {
	c := cron.New()
	spec := "0 0 1 * * ?"
	err := c.AddFunc(spec, func() {
		IS := false
		if !Do() {
			if Do() {
				IS = true
			}
		} else {
			IS = true
		}
		if IS {
			_, _ = bot.SendPrivateMsg(3343780376, "打卡成功\nhttp://47.110.228.1/log/"+time.Now().Format("2006-01-02")+".log", false)
		} else {
			_, _ = bot.SendPrivateMsg(3343780376, "打卡失败", false)
		}
	})
	if err != nil {
		log.Println(err)
	}
	c.Start()
	select {}
}

func Do() bool {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("打卡失败")
		}
	}()

	var Date []date
	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	file, err := os.Open(dir + "/plugins/daka/daka.json")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Date)
	if err != nil {
		log.Panic(err)
	}
	for _, v := range Date {
		err := commit(v.Xh, v.Xm)
		if err != nil {
			panic(err)
			return false

		}
	}
	return true
}
func commit(xh, xm string) error {

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
	if err != nil {
		return err
	}
	if d.Code == 0 {
		WriteFile("打早卡成功")
	} else if d.Code == 1 {
		WriteFile("今日早卡已打卡")
	}
	d, err = commitData(response, client, 2)
	if err != nil {
		return err
	}
	if d.Code == 0 {
		WriteFile("打午卡成功")
	} else if d.Code == 1 {
		WriteFile("今日午卡已打卡")
	}
	d, err = commitData(response, client, 3)
	if err != nil {
		return err
	}
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
	file, err := os.OpenFile(path+"/plugins/logs/"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Panic(err)
	}
	_, err = file.WriteString(time.Now().Format("2006-01-02 15:04:05") + " " + string2 + "\n\n")
	if err != nil {
		panic("log  " + string2 + "  error")
	}
}
