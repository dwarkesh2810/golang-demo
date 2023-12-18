package middleware

import (
	"sync"
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

// func RateLimiter(ctx *context.Context) {
// 	// Get IP address of the client
// 	ip := ctx.Input.IP()

// 	// Limit requests from an IP address
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	if timeOut[ip] == 0 {
// 		timeOut[ip] = time.Now().Add(60 * time.Second).Unix()
// 	}

// 	if timeOut[ip] < time.Now().Unix() {
// 		timeOut[ip] = 0
// 		accessCount[ip] = 0
// 		blocked[ip] = false
// 		unBlocked[ip] = 0
// 	}

// 	accessCount[ip]++
// 	if accessCount[ip] > conf.EnvConfig.RateLimiter {
// 		blocked[ip] = true
// 		if blocked[ip] && unBlocked[ip] > 0 {
// 			if unBlocked[ip] < time.Now().Unix() {
// 				accessCount[ip] = 0
// 				blocked[ip] = false
// 				unBlocked[ip] = 0
// 				return
// 			}
// 		}

// 		if blocked[ip] && unBlocked[ip] > 0 {
// 			unBlocked[ip] = unBlocked[ip] + conf.EnvConfig.BlockTime
// 			timeOut[ip] = unBlocked[ip]
// 		} else {
// 			log.Print(conf.EnvConfig.BlockTime)
// 			unBlocked[ip] = int64(time.Now().Add(time.Duration(conf.EnvConfig.BlockTime) * time.Second).Unix())
// 			timeOut[ip] = unBlocked[ip]
// 		}

// 		remainingSeconds := unBlocked[ip] - time.Now().Unix()

// 		day, hr, min, sec := helpers.SecondsToDayHourMinAndSeconds(int(remainingSeconds))

// 		message := fmt.Sprintf("Too many request, Please try again after %d days %d hours %d min %d.", day, hr, min, sec)

// 		resp := &RateLimmiterResponse{
// 			Message: message,
// 		}
// 		data, _ := json.Marshal(resp)
// 		ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
// 		ctx.ResponseWriter.WriteHeader(http.StatusTooManyRequests)
// 		ctx.ResponseWriter.Write(data)
// 		return
// 	}
// }
