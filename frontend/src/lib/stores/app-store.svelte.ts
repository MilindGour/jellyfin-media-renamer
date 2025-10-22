import type {
  SourceDirWithInfo,
  SourceDirectory,
  SourceDirectoriesResponse,
  RenameMediaResponseItem,
  Config,
  Source,
  DestConfig
} from "$lib/models";
import type { API } from "$lib/services/api";
import { Log } from "$lib/services/logger";

const log = new Log("app-store");

// Main store class for the whole application
export class JmrApplicationStore {
  static instance?: JmrApplicationStore;

  /* Variables */
  config = $state<Config | null>(null);
  #source = $state<Source | null>(null);
  #selectSourceDirectories = $state<SourceDirectory[]>([]);

  sourceDirectories = $state<SourceDirectoriesResponse | null>(null);
  sourceDirsWithMediaInfo = $state<SourceDirWithInfo[]>([]);
  mediaSelectionForRenames = $state<RenameMediaResponseItem[]>([]);
  mediaDestinationSelections = $state<DestConfig[]>([]);

  /* Constructor */
  constructor(private api: API) {
    if (!JmrApplicationStore.instance) {
      JmrApplicationStore.instance = this;
    }
    return JmrApplicationStore.instance;
  }

  async setConfig(config: Config) {
    this.config = config;
  }
  /* First page related methods */
  async setSource(s: Source) {
    this.#source = s;
    if (this.#source === null) {
      this.sourceDirectories = null;
    }
    this.sourceDirectories = await this.api.getSourceDirectoriesAsync(this.#source);
  }
  async setSourceDirectories(s: SourceDirectory[]) {
    this.#selectSourceDirectories = s;
    if (this.#selectSourceDirectories.length === 0) {
      this.sourceDirsWithMediaInfo = [];
    }
    const sourceWithNames = await this.api.identifyMediaNames(this.#selectSourceDirectories);
    if (sourceWithNames) {
      this.sourceDirsWithMediaInfo = sourceWithNames;
    } else {
      this.sourceDirsWithMediaInfo = [];
    }
  }

  /* Second page related methods */
  async searchMediaInfoProvider() {
    if (this.sourceDirsWithMediaInfo.length > 0) {
      const result = await this.api.identifyMediaInfos(this.sourceDirsWithMediaInfo);
      if (result && result.length === this.sourceDirsWithMediaInfo.length) {
        // both the equal in length, can be assigned safely.
        this.sourceDirsWithMediaInfo = result;
      } else {
        log.error("Cannot identify media info from given names. Unknown error occured.");
      }
    }
  }
  async getMediaSelectionForRenames() {
    if (this.sourceDirsWithMediaInfo.length > 0) {
      const result = await this.api.getMediaSelectionsForRenames(this.sourceDirsWithMediaInfo);
      if (result && result.length === this.sourceDirsWithMediaInfo.length) {
        this.mediaSelectionForRenames = result;
        this.mediaDestinationSelections = new Array(result.length).fill(null);
      } else {
        log.error("Cannot get media renames. Unknown error occured.");
        this.mediaSelectionForRenames = [];
        this.mediaDestinationSelections = [];
      }
    }
  }

  async confirmMediaRequest() {
    if (this.mediaSelectionForRenames.length > 0 && this.mediaDestinationSelections.length > 0 && this.mediaSelectionForRenames.length === this.mediaDestinationSelections.length) {
      // get preview of renaming, and possibly start the renaming process
      const result = await this.api.confirmMediaRenames(this.mediaSelectionForRenames, this.mediaDestinationSelections);
      log.info("Media Confirm Response:", result);
    } else {
      log.error("Media selections and media destinations must have a non-zero equal length");
    }
  }
}
