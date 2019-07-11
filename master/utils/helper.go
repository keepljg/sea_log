package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"sea_log/master/conf"
	"strings"
	"time"
)

func GetNodeIp(key string) string {
	return strings.Replace(strings.TrimPrefix(key, conf.JobConf.JobInfo), "/", "", -1)
}

func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, conf.JobConf.JobSave)
}

func ExtractRuningJob(pre ,jobKey string) ([]string, error) {
	info := strings.TrimPrefix(jobKey, pre)
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
