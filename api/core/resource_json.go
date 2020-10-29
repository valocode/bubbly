package core

type ResourceJSON struct {
	Kind     string `json:"kind,string,omitempty"`
	Name     string `json:"name,string,omitempty"`
	Resource string `json:"resource,string,omitempty"`
}
