import type { Source, SourceDirectoriesResponse, SourcesResponse } from "$lib/models/models";
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

    // Static methods
    static http(): API {
        return new API(new HttpService());
    }
}
