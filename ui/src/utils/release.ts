
import type { Release } from '$schema/schema_gen';
import { ReleaseStatus } from '$schema/schema_gen';

export const releaseStatusColor = (release: Release): string => {
	switch (release.status) {
		case ReleaseStatus.ready:
			return 'green';
		case ReleaseStatus.blocked:
			return 'red';
		case ReleaseStatus.pending:
			return 'yellow';
		default:
			return 'gray';
	}
};

export const releaseStatusText = (release: Release): string => {
	switch (release.status) {
		case ReleaseStatus.ready:
			return 'Ready';
		case ReleaseStatus.blocked:
			return 'Blocked';
		case ReleaseStatus.pending:
			return 'Pending';
		default:
			return 'Status Unknown';
	}
};
