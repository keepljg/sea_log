package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"sea_log/master/conf"
	"strings"
	"time"
)

func GetNodeIp(key string) string {
	return strings.TrimPrefix(key, conf.JobConf.JobInfo)
}

func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, conf.JobConf.JobSave)
}

func ExtractRuningJob(jobKey string) ([]string, error) {
	info := strings.TrimPrefix(jobKey, conf.JobConf.JobLock)
	res := strings.Split(info, "/")
	if len(res) != 2 {
		return nil, errors.New("Runing Job Extract Failed")
	}
	return res, nil
}

func ExtractJobSave(ip, name string) string {
	return fmt.Sprintf("%s%s/%s", conf.JobConf.JobSave, ip, name)
}

func RangeInt(start, end int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(end-start+1) + start
}

func Success() (int, interface{}) {
	return http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	}
}

func SuccessWithDate(data interface{}) (int, interface{}) {
	return http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	}
}
