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

export class HttpService {
    async getJSON<T>(url: string, queryParams: URLSearchParams | null = null) {
        const apiUrl = getApiUrl(url, queryParams);
        const response = await fetch(apiUrl, {
            headers: {
                'Content-Type': 'application/json'
            },
            method: "GET",
            redirect: "error",
        });
        return response.json() as T;
    }
    async postJSON<T>(url: string, body: any, queryParams: URLSearchParams | null = null) {
        const apiUrl = getApiUrl(url, queryParams);
        const response = await fetch(apiUrl, {
            headers: {
                'Content-Type': 'application/json'
            },
            body,
            method: "POST",
            redirect: "error",
        });
        return response.json() as T;
    }
}
