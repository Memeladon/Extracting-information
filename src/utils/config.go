package utils

type Settings struct {
	DBUser string `json:"DBUser"`
	DBPass string `json:"DBPass"`
	DBHost string `json:"DBHost"`
	DBPort string `json:"DBPort"`
	DBName string `json:"DBName"`
	DBUrl  string `json:"DBUrl"`
}
