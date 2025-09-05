import { dev } from '$app/environment';

export const Constants = {
    API_BASE_URL_DEV: "http://localhost:7749",
    API_BASE_URL: "https://jmr.miracular.in",

    API_GET_SOURCE: "api/sources",
    API_GET_SOURCE_DIRS: "api/sources/:id",
};

export function GetApiBaseUrl(): string {
    return dev ? Constants.API_BASE_URL_DEV : Constants.API_BASE_URL;
}

