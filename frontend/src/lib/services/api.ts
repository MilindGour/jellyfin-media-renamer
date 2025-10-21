import type { RenameMediaResponseItem, SourceDirectoriesResponse, SourceDirectory, SourceDirWithInfo, Source, ConfirmMediaRequestItem, DestConfig } from "$lib/models";
import { Constants } from "$lib/stores/constants";
import { Log } from "./logger";
import { HttpService } from "./network";

const log = new Log("API");

export class API {
  constructor(private http: HttpService) { }

  async getSourceDirectoriesAsync(s: Source): Promise<SourceDirectoriesResponse | null> {
    try {
      const apiUrl = Constants.API_GET_SOURCE_DIRS.replace(":id", s.id.toString());
      const res = await this.http.getJSON<SourceDirectoriesResponse>(apiUrl);
      return res;

    } catch (err) {
      log.error("getSourceDirectoriesAsync err:", err);
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

  /**
  * Returns all the rename previews for the selected media entries.
  */
  async getMediaSelectionsForRenames(input: SourceDirWithInfo[]): Promise<RenameMediaResponseItem[] | null> {
    try {
      const apiUrl = Constants.API_POST_MEDIA_RENAMES;
      const res = await this.http.postJSON<RenameMediaResponseItem[]>(apiUrl, input);
      return res;

    } catch {
      return null;
    }
  }

  async confirmMediaRenames(mediaSelections: RenameMediaResponseItem[], mediaDestinations: DestConfig[]): Promise<any[] | null> {
    try {
      const input: ConfirmMediaRequestItem[] = [];
      for (let i = 0; i < mediaSelections.length; i++) {
        input.push({
          ...mediaSelections[i],
          destination: mediaDestinations[i]
        })
      }
      log.info("TODO: confirm and sync request:", input);
      const apiUrl = Constants.API_POST_MEDIA_RENAMES_CONFIRM;
      const res = await this.http.postJSON<any[]>(apiUrl, input);
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
