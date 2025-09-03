export type SourcesResponse = {
    sources: Source[]
};

export type Source = {
    name: string,
    path: string,
    id: number,
};

export type SourceDirectoriesResponse = {
    source: Source,
    entries: DirEntry[]
};

export type DirEntry = {
    id: number,
    name: string,
    path: string,
    size: number,

    isDirectory: boolean,
    children?: DirEntry[]
};
