package common

import "errors"

var (
	ErrLockBusy  = errors.New("锁被占用")
	ErrNoIpFound = errors.New("机器没有物理网卡")
)
