import type { Source, SourceDirectory } from "$lib/models/models";
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
    setSourceDirectories(s: SourceDirectory[]) {
        this.#selectSourceDirectories = s;
    }
}
