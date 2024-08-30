package controllers

import "time"

func UpdateStockTimestamp() *time.Time {
	now := time.Now().Format("2006-02-01T15:04:05.999999")
	nowT,_ := time.Parse("2006-02-01T15:04:05.999999",now)
	nowT = nowT.Add(-(8*time.Hour))
	return &nowT
}