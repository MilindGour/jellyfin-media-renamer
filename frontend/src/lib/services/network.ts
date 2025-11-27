import { v4 as uuidv4 } from 'uuid';
import { loaderService } from "$lib/components/loader";
import { GetApiBaseUrl, getWSBaseUrl } from "$lib/stores/constants";
import { Log } from "./logger";

function getBaseUrl(): string {
  return GetApiBaseUrl();
}

export function getApiUrl(url: string, searchParams: URLSearchParams | null = null): string {
  const baseUrl = getBaseUrl();

  let finalUrl = `${baseUrl}/${url}`;
  if (searchParams !== null) {
    finalUrl += '?' + searchParams.toString();
  }
  return finalUrl;
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

export type WSCallbackFunction<T> = (message: T) => void;

export class WebSocketService {
  static #instance: WebSocketService;
  #ws?: WebSocket;
  #log?: Log;
  #uuid?: string;
  #listeners: { [key: string]: Array<WSCallbackFunction<any>> } = {};

  constructor() {
    if (!WebSocketService.#instance) {
      WebSocketService.#instance = this;
      this.#log = new Log("WebSocket");
      this.#uuid = uuidv4();
    }
    return WebSocketService.#instance;
  }

  connect() {
    this.disconnect();
    const wsurl = this.getWSURL();
    this.#log?.info("Connecting to websocket:", wsurl);
    this.#ws = new WebSocket(wsurl);

    this.#ws.onopen = this.handleOnOpen;
    this.#ws.onclose = this.handleOnClose;
    this.#ws.onmessage = this.handleOnMessage;
  }

  disconnect() {
    this.#ws?.close();
  }

  getWSURL(): string {
    const wsBaseURL = getWSBaseUrl();
    return `${wsBaseURL}/${this.#uuid}`;
  }

  handleOnOpen = (event: Event) => {
    this.#log!.info("connection opened with event:", event);
  };

  handleOnClose = (event: CloseEvent) => {
    this.#log!.info("connection closed with event:", event);
  };

  handleOnMessage = (event: MessageEvent<string>) => {
    this.#log!.info("message received:", event);
    const parsedData: WSMessage = JSON.parse(event.data);

    if (parsedData.type in this.#listeners) {
      this.#listeners[parsedData.type].forEach(cb => cb(parsedData.data));
    }
  };

  addListener<T>(messageType: "progress", callback: WSCallbackFunction<T>) {
    if (!(messageType in this.#listeners)) {
      this.#listeners[messageType] = [];
    }
    this.#listeners[messageType].push(callback);
  }

  removeListener(messageType: "progress", callback: WSCallbackFunction<any>): boolean {
    if (messageType in callback) {
      this.#listeners[messageType] = this.#listeners[messageType].filter(x => x !== callback);
      return true;
    }
    return false;
  }
}

export type WSMessage = {
  type: string,
  data: string
};

export type PathPair = {
  old_path: string,
  new_path: string,
};

export type FileTransferData = {
  bytes_transferred: number,
  error: any,
  files: PathPair,
  percent_complete: number,
  raw_string: string,
  time_remaining: string,
  transfer_speed: string,
  total_bytes: number
};
export type ProgressData = Array<FileTransferData>;
