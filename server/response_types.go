package server

type Error struct {
	ErrMessage string `json:"error"`
}
type Status struct {
	Status string `json:"status"`
}
type Data struct {
	Data interface{} `json:"data"`
}

type VersionHeaders struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}
