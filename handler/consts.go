package handler

const (
	APPLY = iota
	PLAN
	DESTROY
)

const (
	CREATED = iota
	QUEUED
	RUNNING
	COMPLETE
)

const (
	AWS = iota
	AZURE
	GCP
)