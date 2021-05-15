package builtin

func ReleaseStatusByStages(release Release) string {
	for _, stage := range release.ReleaseStage {
		status := ReleaseStageStatus(stage)
		if status != "READY" {
			return status
		}
	}
	return "READY"
}

func ReleaseStageStatus(stage ReleaseStage) string {
	for _, criteria := range stage.ReleaseCriteria {
		status := ReleaseCriteriaStatus(criteria)
		if status != "READY" {
			return status
		}
	}
	return "READY"
}

func ReleaseCriteriaStatus(criteria ReleaseCriteria) string {
	if len(criteria.ReleaseEntry) == 0 {
		return "PENDING"
	}
	// We only care about the first (or latest) release entry
	if !criteria.ReleaseEntry[0].Result {
		return "BLOCKED"
	}
	return "READY"
}
