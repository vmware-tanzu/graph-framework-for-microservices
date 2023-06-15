package common

type Application struct {
	Name   string `json:"name, omitempty"`
	SemVer string `json:"sem_ver, omitempty"`
}

type ApplicationStatus struct {
	NexusApps map[string]NexusApp `json:"nexusApps, omitempty"`
}

type NexusApp struct {
	OamApp      OamApp `json:"oamApp, omitempty"`
	State       string `json:"state, omitempty"`
	StateReason string `json:"stateReason, omitempty"`
}

type OamApp struct {
	Components map[string]ComponentDefinition `json:"components, omitempty"`
}

type ComponentDefinition struct {
	Name   string `json:"name, omitempty"`
	Sync   string `json:"sync, omitempty"`
	Health string `json:"health, omitempty"`
}
