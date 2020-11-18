package core

import "strings"

// ResourceJSON represents a JSON format of a resource
// the resource is stored as a string for storage in BuntDB
type ResourceJSON struct {
	Kind     string `json:"kind,string,omitempty"`
	Name     string `json:"name,string,omitempty"`
	Resource string `json:"resource,string,omitempty"`
	Namespace string `json:"namespace,string,omitempty"`
}


// GetID returns a string of the resource's ID
func (r *ResourceJSON) GetID() string {
	return strings.Join([]string{r.Namespace, r.Kind, r.Name}, "/")
}

