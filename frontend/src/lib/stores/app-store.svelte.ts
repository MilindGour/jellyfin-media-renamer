import type { SourceDirWithInfo, Source, SourceDirectory, SourceDirectoriesResponse } from "$lib/models/models";
import type { API } from "$lib/services/api";
import { Log } from "$lib/services/logger";

const log = new Log("app-store");
// Main store class for the whole application
export class JmrApplicationStore {
    static instance?: JmrApplicationStore;

    /* Variables */
    #source = $state<Source | null>(null);
    #selectSourceDirectories = $state<SourceDirectory[]>([]);

    sourceDirectories = $state<SourceDirectoriesResponse | null>(null);
    sourceDirsWithMediaInfo = $state<SourceDirWithInfo[]>([]);

    /* Constructor */
    constructor(private api: API) {
        if (!JmrApplicationStore.instance) {
            JmrApplicationStore.instance = this;
        }
        return JmrApplicationStore.instance;
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
    async confirmMediaInfos() {
        // TODO: Confirm the final media IDs before moving on to next page.
        throw Error("Not implemented");
    }
}
