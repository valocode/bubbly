package builtin

func ReleaseStatusByStages(release Release) string {

	for _, stage := range release.ReleaseStage {
		for _, criteria := range stage.ReleaseCriteria {
			if len(criteria.ReleaseEntry) == 0 {
				return "PENDING"
			}
			for _, entry := range criteria.ReleaseEntry {
				if !entry.Result {
					return "BLOCKED"
				}
			}
		}
	}
	return "READY"
}
