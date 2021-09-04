import { writable } from "svelte/store";


function createFilterHeadOnly() {
    const { set, subscribe } = writable<boolean>(false);
    let state: boolean;
    subscribe((b) => (state = b));
    return {
        set,
        subscribe,
        isEnabled: () => state,
        toggled: () => { set(!state) }
    };
};
export const filterHeadOnly = createFilterHeadOnly();

export type ProjectSelectMap = Record<string, boolean>
function createFilterProjects() {
    const { set, subscribe } = writable<ProjectSelectMap>({});
    let projects: ProjectSelectMap;
    subscribe((value) => (projects = value));
    return {
        set,
        subscribe,
        selectedProjects: (): string[] => Object.entries(projects).filter((p) => p[1]).map((p) => p[0]),
        toggle: (name: string) => {
            projects[name] = !projects[name]
            set(projects)
        }
    };
};
export const filterProjects = createFilterProjects();


export enum SortReleaseFields {
    name = "Name",
    version = "Version",
}

export const SortFields = Object.values(SortReleaseFields)

function createSortReleases() {
    const { set, subscribe } = writable<SortReleaseFields>(null);
    let field: SortReleaseFields;
    subscribe((value) => (field = value));
    return {
        set,
        subscribe,
        sortByField: (): string => field
    };
};
export const sortReleases = createSortReleases();
