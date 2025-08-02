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
  basePath: string;
  directoryEntries: DirEntry[]
}
