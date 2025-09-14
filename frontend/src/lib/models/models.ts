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

export type MediaType = "MOVIE" | "TV";

export type SourceDirectory = {
    entry: DirEntry,
    type: MediaType | null,
    selected: boolean
};

export type MediaInfo = {
    name: string,
    description: string,
    yearOfRelease: number,
    thumbnailUrl: string,
    mediaId: string
}

export type SourceDirWithInfo = {
    sourceDirectory: SourceDirectory,
    identifiedMediaName?: string,
    identifiedMediaYear?: number,
    identifiedMediaInfos?: MediaInfo[],
};
