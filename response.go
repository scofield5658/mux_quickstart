package main

import (
	"encoding/json"
	"time"
)

type TimeStamp int64

type FileUploadResponse struct {
	Filename     string    `json:"filename"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	ModifiedTime TimeStamp `json:"modified_time"`
}

func (d TimeStamp) MarshalJSON() ([]byte, error) {
	rs := time.Unix(int64(d), 0).Format("2006-01-02 15:04:05")
	js, er := json.Marshal(rs)
	return js, er
}
func (d *TimeStamp) UnmarshalJSON(data []byte) error {
	var rs string
	e := json.Unmarshal(data, &rs)
	if e != nil {
		return e
	}
	t, er := time.Parse("2006-01-02 15:04:05", rs)
	if er != nil {
		return er
	}
	*d = TimeStamp(t.Unix())
	return nil
}

type ErrorResponse struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}
