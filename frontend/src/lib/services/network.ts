import { loaderService } from "$lib/components/loader/loader-store.svelte";
import { Constants, GetApiBaseUrl } from "$lib/stores/constants";

function getBaseUrl(): string {
  return GetApiBaseUrl();
}

export function getApiUrl(url: string, searchParams: URLSearchParams | null = null): string {
  const _url = new URL(url, getBaseUrl());
  if (searchParams !== null) {
    for (const [spName, spValue] of searchParams.entries()) {
      _url.searchParams.append(spName, spValue);
    }
  }

  return _url.href;
}

export class HttpService {

  async getJSON<T>(url: string, queryParams: URLSearchParams | null = null) {
    loaderService.addAPICounter();
    const apiUrl = getApiUrl(url, queryParams);

    try {
      const response = await fetch(apiUrl, {
        headers: {
          'Content-Type': 'application/json'
        },
        method: "GET",
        redirect: "error",
      });
      const result = await response.json() as T;
      return result;

    } catch (error) {
      throw new Error("[JMR Network] Error", { cause: error });
    } finally {
      loaderService.subtractAPICounter();
    }
  }
  async postJSON<T>(url: string, body: any, queryParams: URLSearchParams | null = null) {
    loaderService.addAPICounter();
    const apiUrl = getApiUrl(url, queryParams);

    try {
      const response = await fetch(apiUrl, {
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(body),
        method: "POST",
        redirect: "error",
      });
      const result = await response.json() as T;
      return result;

    } catch (error) {
      throw new Error("[JMR Network] Error", { cause: error });
    } finally {
      loaderService.subtractAPICounter();
    }
  }
}

export class WebSocketService {
  static #instance: WebSocketService;
  #ws?: WebSocket;

  constructor() {
    if (!WebSocketService.#instance) {
      WebSocketService.#instance = this;
    }
    return WebSocketService.#instance;
  }
  connect() {
    if (this.#ws) {
      this.#ws.close(0, "Initializing another instance");
    }
    this.#ws = new WebSocket(this.getWSURL());
  }
  getWSURL(): string {
    const wsBaseURL = GetApiBaseUrl().replace(/^https?:\/\//, '');
    return `ws://${wsBaseURL}/${Constants.API_WS}`;
  }
}
