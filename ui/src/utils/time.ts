
export const friendlyTimeElapsed = (timeFrom: Date, timeTo: Date): string => {
    const elapsed = new Date(timeTo).getTime() - new Date(timeFrom).getTime();
    return friendlyTime(elapsed)
};

export const friendlyTime = (elapsed: number): string => {
    const msPerMinute = 60 * 1000;
    const msPerHour = msPerMinute * 60;
    const msPerDay = msPerHour * 24;
    const msPerMonth = msPerDay * 30;
    const msPerYear = msPerDay * 365;
    if (elapsed < msPerMinute) {
        return Math.round(elapsed / 1000) + ' seconds';
    } else if (elapsed < msPerHour) {
        return Math.round(elapsed / msPerMinute) + ' minutes';
    } else if (elapsed < msPerDay) {
        return Math.round(elapsed / msPerHour) + ' hours';
    } else if (elapsed < msPerMonth) {
        return Math.round(elapsed / msPerDay) + ' days';
    } else if (elapsed < msPerYear) {
        return Math.round(elapsed / msPerMonth) + ' months';
    } else {
        return Math.round(elapsed / msPerYear) + ' years';
    }
}
