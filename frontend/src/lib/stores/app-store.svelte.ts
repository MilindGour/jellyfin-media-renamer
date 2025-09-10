import type { SourceDirWithInfo, Source, SourceDirectory } from "$lib/models/models";
import type { API } from "$lib/services/api";

// Main store class for the whole application
export class JmrApplicationStore {
    static instance?: JmrApplicationStore;

    /* Variables */
    #source = $state<Source | null>(null);
    #selectSourceDirectories = $state<SourceDirectory[]>([]);

    sourceDirectories = $derived.by(async () => {
        if (this.#source === null) {
            return null;
        }
        const srcRes = await this.api.getSourceDirectoriesAsync(this.#source);
        return srcRes;
    });

    sourceDirsWithMediaInfo = $state<SourceDirWithInfo[]>([]);
    sourceDirsWithMediaNames = $derived.by(async () => {
        if (this.#selectSourceDirectories.length === 0) {
            return null;
        }
        const sourceWithNames = await this.api.identifyMediaNames(this.#selectSourceDirectories);
        return sourceWithNames;
    });

    /* Constructor */
    constructor(private api: API) {
        if (!JmrApplicationStore.instance) {
            JmrApplicationStore.instance = this;
        }
        return JmrApplicationStore.instance;
    }

    /* First page related methods */
    setSource(s: Source) {
        this.#source = s;
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
                console.error("Cannot identify media info from given names. Unknown error occured.");
            }
        }
    }
    async confirmMediaInfos() {
        // TODO: Confirm the final media IDs before moving on to next page.
        throw Error("Not implemented");
    }
}
