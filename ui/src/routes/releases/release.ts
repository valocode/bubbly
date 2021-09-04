import type { Release } from "$schema/schema_gen";
import { ReleasePolicyViolationSeverity } from "$schema/schema_gen";

export const worstSeverity = (release: Release): ReleasePolicyViolationSeverity => {
    let max = null;
    for (let i = 0; i < release.violations.length; i++) {
        const severity = release.violations[i].severity;
        switch (severity) {
            case ReleasePolicyViolationSeverity.blocking:
                return severity;
            case ReleasePolicyViolationSeverity.warning:
                if (max === ReleasePolicyViolationSeverity.suggestion) {
                    max = severity;
                }
        }
    }
    return max;
};