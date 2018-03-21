/*
 * Copyright (c) 2018. 
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/21.

// RunLog HICD运行日志
type RunLog struct {
	// Timestamp 日志时间戳, 精确到s
	Timestamp int64 `json:"timestamp"`
	// Message 日志内容
	Message string `json:"message"`
}
