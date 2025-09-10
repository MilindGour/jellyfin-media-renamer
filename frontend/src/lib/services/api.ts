import type { Source, SourceDirectoriesResponse, SourceDirectory, SourceDirWithInfo } from "$lib/models/models";
import { Constants } from "$lib/stores/constants";
import { HttpService } from "./network";

export class API {
    constructor(private http: HttpService) { }

    async getSourceDirectoriesAsync(s: Source): Promise<SourceDirectoriesResponse | null> {
        try {
            const apiUrl = Constants.API_GET_SOURCE_DIRS.replace(":id", s.id.toString());
            const res = await this.http.getJSON<SourceDirectoriesResponse>(apiUrl);
            return res;

        } catch {
            return null;
        }
    }

    /**
    * Identifies and returns media names using the dirty directory / filenames.
    * @param input The selected source directories.
    * @returns Directory info along with guessed media names.
    */
    async identifyMediaNames(input: SourceDirectory[]): Promise<SourceDirWithInfo[] | null> {
        try {
            const apiUrl = Constants.API_POST_IDENTIFY_MEDIA_NAMES;
            const res = await this.http.postJSON<SourceDirWithInfo[]>(apiUrl, input);
            return res;

        } catch {
            return null;
        }
    }

    /**
    * Identifies and returns media info (zero or more) after searching on media information provider.
    * Must contain media names for searching.
    * @param input The selected source directories with media names.
    * @returns Same directory info after appending media info from the provider.
    */
    async identifyMediaInfos(input: SourceDirWithInfo[]): Promise<SourceDirWithInfo[] | null> {
        try {
            const apiUrl = Constants.API_POST_IDENTIFY_MEDIA_INFO;
            const res = await this.http.postJSON<SourceDirWithInfo[]>(apiUrl, input);
            return res;

        } catch {
            return null;
        }
    }
    async confirmMediaInfo(input: SourceDirWithInfo[]) {
        throw Error("Not implemented");
    }

    // Static methods
    static http(): API {
        return new API(new HttpService());
    }
}
