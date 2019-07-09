package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UnmarshalJobInfo(info []byte) (nodeInfo NodeInfo, err error) {
	err = json.Unmarshal(info, &nodeInfo)
	return
}

func UnPackJob(value []byte) (*Jobs, error) {
	var (
		err error
		job Jobs
	)
	if err = json.Unmarshal(value, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

func PackJob(jobs Jobs) ([]byte, error) {
	return json.Marshal(&jobs)
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
