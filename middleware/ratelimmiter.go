package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/dwarkesh2810/golang-demo/conf"
	"github.com/dwarkesh2810/golang-demo/helpers"
)

var (
	accessCount = make(map[string]int)
	mutex       = &sync.Mutex{}
	blocked     = make(map[string]bool)
	unBlocked   = make(map[string]int64)
	timeOut     = make(map[string]int64)
)

type RateLimmiterResponse struct {
	Message string `json:"message"`
}

func RateLimiter(ctx *context.Context) {
	// Get IP address of the client
	ip := ctx.Input.IP()
	limit, err := strconv.Atoi(conf.ConfigMaps["RATELIMITER"])
	if err != nil {
		log.Print(err)
	}
	blockTime, err := strconv.Atoi(conf.ConfigMaps["BLOCKTIME"])

	if err != nil {
		log.Print(err)
	}

	// Limit requests from an IP address
	mutex.Lock()
	defer mutex.Unlock()

	if timeOut[ip] == 0 {
		timeOut[ip] = time.Now().Add(60 * time.Second).Unix()
	}

	if timeOut[ip] < time.Now().Unix() {
		timeOut[ip] = 0
		accessCount[ip] = 0
		blocked[ip] = false
		unBlocked[ip] = 0
	}

	accessCount[ip]++
	if accessCount[ip] > limit {
		blocked[ip] = true
		if blocked[ip] && unBlocked[ip] > 0 {
			if unBlocked[ip] < time.Now().Unix() {
				accessCount[ip] = 0
				blocked[ip] = false
				unBlocked[ip] = 0
				return
			}
		}

		if blocked[ip] && unBlocked[ip] > 0 {
			unBlocked[ip] = unBlocked[ip] + int64(blockTime)
			timeOut[ip] = unBlocked[ip]
		} else {
			unBlocked[ip] = int64(time.Now().Add(time.Duration(int64(blockTime)) * time.Second).Unix())
			timeOut[ip] = unBlocked[ip]
		}

		remainingSeconds := unBlocked[ip] - time.Now().Unix()

		day, hr, min, sec := helpers.SecondsToDayHourMinAndSeconds(int(remainingSeconds))

		message := fmt.Sprintf("Too many request, Please try again after %d days %d hours %d min %d.", day, hr, min, sec)

		resp := &RateLimmiterResponse{
			Message: message,
		}
		data, _ := json.Marshal(resp)
		ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
		ctx.ResponseWriter.WriteHeader(http.StatusTooManyRequests)
		ctx.ResponseWriter.Write(data)
		return
	}
}
