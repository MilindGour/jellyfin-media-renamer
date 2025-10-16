import { dev } from '$app/environment';

export const Constants = {
  API_BASE_URL_DEV: "http://localhost:7749",
  API_BASE_URL: "https://jmr.miracular.in",

  API_GET_CONFIG: "api/config",
  API_GET_SOURCE: "api/sources",
  API_GET_SOURCE_DIRS: "api/sources/:id",
  API_POST_IDENTIFY_MEDIA_NAMES: "api/media/identify-names",
  API_POST_IDENTIFY_MEDIA_INFO: "api/media/identify-info",
  API_POST_IDENTIFY_CONFIRM_MEDIA_INFO: "api/media/confirm-info",
  API_POST_MEDIA_RENAMES: "api/media/rename",
  API_POST_MEDIA_RENAMES_CONFIRM: "api/media/rename-confirm",
  API_WS: "api/ws",
};

export function GetApiBaseUrl(): string {
  return dev ? Constants.API_BASE_URL_DEV : Constants.API_BASE_URL;
}

