package models

// Crack password structure
type Target struct {
	IP       string
	Port     int
	Protocol string
}

type Service struct {
	IP       string
	Port     int
	Protocol string
	Username string
	Password string
}

type CrackResult struct {
	Service Service
	Result  bool
}
