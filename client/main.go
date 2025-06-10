package main

import (
	"akile_client/model"
	"bytes"
	"compress/gzip"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/gorilla/websocket"
	"github.com/henrylee2cn/goutil/calendar/cron"
)

func main() {
	LoadConfig()

	go func() {
		c := cron.New()
		c.AddFunc("* * * * * *", func() {
			TrackNetworkSpeed()
		})
		c.Start()
	}()

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := cfg.Url
	log.Printf("connecting to %s", u)

	c, _, err := websocket.DefaultDialer.Dial(cfg.Url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	c.WriteMessage(websocket.TextMessage, []byte(cfg.AuthSecret))

	done := make(chan struct{})

	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("auth_secret验证失败")
		log.Println("read:", err)
		return
	}
	if string(message) == "auth success" {
		log.Println("auth_secret验证成功")
		log.Println("正在上报数据...")
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			var D struct {
				Host      *model.Host
				State     *model.HostState
				TimeStamp int64
			}
			D.Host = GetHost()
			D.State = GetState()
			D.TimeStamp = t.Unix()
			//gzip压缩json
			dataBytes, err := json.Marshal(D)
			if err != nil {
				log.Println("json.Marshal error:", err)
				return
			}

			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			if _, err := gz.Write(dataBytes); err != nil {
				log.Println("gzip.Write error:", err)
				return
			}

			if err := gz.Close(); err != nil {
				log.Println("gzip.Close error:", err)
				return
			}

			err = c.WriteMessage(websocket.TextMessage, buf.Bytes())
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
