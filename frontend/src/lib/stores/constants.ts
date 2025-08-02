import { dev } from '$app/environment';

export const Constants = {
  API_BASE_URL_DEV: "http://localhost:7749",
  API_BASE_URL: "https://jmr.miracular.in",

  API_GET_CONFIG_SOURCE: "api/config/source",
  API_GET_CONFIG_SOURCE_BY_ID: "api/config/source/:id",
};

export function GetApiBaseUrl(): string {
  return dev ? Constants.API_BASE_URL_DEV : Constants.API_BASE_URL;
}

