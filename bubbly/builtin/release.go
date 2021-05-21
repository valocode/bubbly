package builtin

type ReleaseStatus string

const (
	ReadyReleaseStatus   ReleaseStatus = "READY"
	PendingReleaseStatus ReleaseStatus = "PENDING"
	BlockedReleaseStatus ReleaseStatus = "BLOCKED"
)

func ReleaseStatusByStages(release Release) ReleaseStatus {
	for _, stage := range release.ReleaseStage {
		status := ReleaseStageStatus(stage)
		if status != ReadyReleaseStatus {
			return status
		}
	}
	return ReadyReleaseStatus
}

func ReleaseStageStatus(stage ReleaseStage) ReleaseStatus {
	for _, criteria := range stage.ReleaseCriteria {
		status := ReleaseCriteriaStatus(criteria)
		if status != ReadyReleaseStatus {
			return status
		}
	}
	return ReadyReleaseStatus
}

func ReleaseCriteriaStatus(criteria ReleaseCriteria) ReleaseStatus {
	if len(criteria.ReleaseEntry) == 0 {
		return PendingReleaseStatus
	}
	// We only care about the first (or latest) release entry
	if !criteria.ReleaseEntry[0].Result {
		return BlockedReleaseStatus
	}
	return ReadyReleaseStatus
}
