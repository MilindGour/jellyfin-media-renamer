import { GetApiBaseUrl } from "$lib/stores/constants";

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

export class JNetworkClient {
  httpGetJSON<T>(url: string, queryParams: URLSearchParams | null = null) {
    const apiUrl = getApiUrl(url, queryParams);
    return fetch(apiUrl, {
      headers: {
        'Content-Type': 'application/json'
      }
    })
      .then(response => response.json() as T);
  }
}
