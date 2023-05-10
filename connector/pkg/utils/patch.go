package utils

import (
	"encoding/json"
)

type Nexus struct {
	SourceGeneration int64 `json:"sourceGeneration"`
	RemoteGeneration int64 `json:"remoteGeneration"`
}

type Status struct {
	Nexus Nexus `json:"nexus"`
}

type PatchOp struct {
	Status Status `json:"status"`
}

func CreatePatch(sourceGenerationID, remoteGenerationID int64) ([]byte, error) {
	payload := PatchOp{
		Status: Status{
			Nexus: Nexus{
				SourceGeneration: sourceGenerationID,
				RemoteGeneration: remoteGenerationID,
			},
		},
	}
	return json.Marshal(payload)
}
