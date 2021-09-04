// calculates the status of a given test_run and returns the status along with

import type { TestRun } from '$model/schema';
import { Status } from '$types/release';

// the number of failed test cases
export function calculateRunStatus(run: TestRun): [status: Status, failCount: number] {
	let status = Status.Passing;
	let failCount = 0;

	for (let c of run.tests) {
		if (c.result === false) {
			failCount++;
			status = Status.Failing;
		}
	}
	return [status, failCount];
}
