export type ConfigSource = {
    id: string,
    name: string,
    path: string
};

export type DirEntry = {
    id: number;
    isDirectory: boolean;
    name: string;
    path: string;
    size: number;
    children: DirEntry[] | null;
}

export type ConfigSourceByID = {
    source: ConfigSource
    entries: DirEntry[]
}

export type ConfigSourcesResponse = {
    sources: ConfigSource[]
}
export type ConfigSourcesByIDResponse = {
    source: ConfigSource
    entries: DirEntry[]
}
