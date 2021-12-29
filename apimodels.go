package main

import (
	"time"
)

type VPSInfo struct {
	ID           uint      `json:"id"`
	DisplayName  string    `json:"display_name"`
	CreationTime time.Time `json:"creation_time"`
	RAM          int       `json:"ram"`
	CPU          int       `json:"cpu"`
	Disk         int       `json:"disk"`
	OS           string    `json:"os"`
}

type RequestInfo struct {
	RequestTime    time.Time `json:"request_time"`
	RequestPurpose int       `json:"request_purpose"`
	// TODO may or may not write api models for request data
	// RequestData    string
	Message string `json:"message"`
}

func VPSToVPSInfo(vps VPS) VPSInfo {
	return VPSInfo{
		ID:           vps.ID,
		DisplayName:  vps.DisplayName,
		CreationTime: vps.CreationTime,
		RAM:          vps.RAM,
		CPU:          vps.CPU,
		Disk:         vps.Disk,
		OS:           vps.OS,
	}
}

func RequestToRequestInfo(request Request) RequestInfo {
	return RequestInfo{
		RequestTime:    request.RequestTime,
		RequestPurpose: request.RequestPurpose,
		Message:        request.Message,
	}
}
