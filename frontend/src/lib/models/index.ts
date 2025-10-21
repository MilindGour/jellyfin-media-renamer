export type SourcesResponse = {
  sources: Source[]
};

export type Config = {
  version: string;
  port: string;
  allowedExtensions: AllowedExtensions;
  source: Source[];
}

export type AllowedExtensions = {
  subtitle: string[];
  media: string[];
}

export type Source = {
  name: string;
  path: string;
  id: number;
}

export type SourceDirectoriesResponse = {
  source: Source,
  entries: DirEntry[]
};

export type DestConfig = {
  name: string;
  path: string;
  id: number;
  type: MediaType;
  mount_point: string;
  total_size_kb: number;
  free_size_kb: number;
  used_size_kb: number;
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
  identifiedMediaId?: string,
  identifiedMediaInfos?: MediaInfo[],
};

export type RenameEntry = {
  subtitle?: DirEntry,
  media: DirEntry,
  season?: number,
  episode?: number
};

export type RenameMediaResponseItem = {
  info: MediaInfo,
  type: MediaType,
  entry: DirEntry,
  selected: RenameEntry[],
  ignored: DirEntry[],
};

export type Filetype = "MEDIA" | "SUBTITLE" | "UNKNOWN";

export type ConfirmMediaRequestItem = RenameMediaResponseItem & {
  destination: DestConfig;
};
