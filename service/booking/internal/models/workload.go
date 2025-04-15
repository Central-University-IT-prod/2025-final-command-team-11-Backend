package models

import (
	"time"
)

type WorkloadItem struct {
	Time   time.Time
	IsFree bool
}

type Workload []WorkloadItem

type FloorWorkloadItem struct {
	Entity BookingEntity
	IsFree bool
}

type FloorWorkload []FloorWorkloadItem
