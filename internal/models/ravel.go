package models

type RavelConfig struct {
	Id   string          `json:"id"`
	Meta RavelConfigMeta `json:"meta"`
	Spec RavelConfigSpec `json:"spec"`
}

type Labels map[string]string
type Scope map[string]string

type RavelResourceMeta struct {
	Name   string `json:"name"`
	Scope  Scope  `json:"scope,omitempty"`
	Labels Labels `json:"labels,omitempty"`
}

type RavelConfigMeta struct {
	RavelResourceMeta
	Version int64 `json:"version,omitempty"`
}

type RavelSchemaMeta struct {
	RavelResourceMeta
	Version string `json:"version,omitempty"`
}

type RavelConfigSpec struct {
	ConfigurationFormat *RavelSchemaMeta `json:"configurationFormat,omitempty"`
	Def                 map[string]any   `json:"def"`
}
