package builtin

type ReleaseStatus string

const (
	ReadyReleaseStatus   ReleaseStatus = "READY"
	PendingReleaseStatus ReleaseStatus = "PENDING"
	BlockedReleaseStatus ReleaseStatus = "BLOCKED"
)

func ReleaseStatusByStages(release Release) ReleaseStatus {
	var status ReleaseStatus = ReadyReleaseStatus
	for _, stage := range release.ReleaseStage {
		stageStatus := ReleaseStageStatus(stage)
		switch stageStatus {
		case BlockedReleaseStatus:
			return BlockedReleaseStatus
		case PendingReleaseStatus:
			status = PendingReleaseStatus
		}
	}
	return status
}

func ReleaseStageStatus(stage ReleaseStage) ReleaseStatus {
	var status ReleaseStatus = ReadyReleaseStatus
	for _, criteria := range stage.ReleaseCriteria {
		criteriaStatus := ReleaseCriteriaStatus(criteria)
		switch criteriaStatus {
		case BlockedReleaseStatus:
			return BlockedReleaseStatus
		case PendingReleaseStatus:
			status = PendingReleaseStatus
		}
	}
	return status
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
