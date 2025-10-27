import type { HttpService } from "$lib/services/network";
import { Constants } from "./constants";

export class NewMediaStore {
  items = $state<NewMediaSearchItem[]>([]);
  #http: HttpService;

  constructor(http: HttpService) {
    this.#http = http;
  }

  async getSearchResultsForTerm(term: string) {
    try {
      this.items = await this.#http.getJSON<NewMediaSearchItem[]>(
        Constants.API_GET_NEW_MEDIA_SEARCH_RESULTS,
        new URLSearchParams([["term", term]])
      );
    } catch {
      this.items = [];
    }
  }

  async addItemToDownloadQueue(item: NewMediaSearchItem) {
    try {
      return await this.#http.postJSON<boolean>(Constants.API_POST_ADD_NEW_MEDIA, item);
    } catch {
      return false;
    }
  }
}

export type NewMediaSearchItem = {
  id: string,
  name: string,
  info_hash: string,
  leechers: string,
  seeders: string,
  num_files: string,
  size: string,
  username: string,
  added: string,
  status: string,
  category: string,
  imdb: string
};
