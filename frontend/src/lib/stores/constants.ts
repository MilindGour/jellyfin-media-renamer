export const Constants = {
  API_BASE_URL: '/backend',
  WS_BASE_URL: "/ws",

  API_GET_CONFIG: "api/config",
  API_GET_SOURCE: "api/sources",
  API_GET_SOURCE_DIRS: "api/sources/:id",
  API_GET_DESTINATIONS: "api/destinations",
  API_POST_IDENTIFY_MEDIA_NAMES: "api/media/identify-names",
  API_POST_IDENTIFY_MEDIA_INFO: "api/media/identify-info",
  API_POST_IDENTIFY_CONFIRM_MEDIA_INFO: "api/media/confirm-info",
  API_POST_MEDIA_RENAMES: "api/media/rename",
  API_POST_MEDIA_RENAMES_CONFIRM: "api/media/rename-confirm",

  API_GET_NEW_MEDIA_SEARCH_RESULTS: "api/new-media/search",
  API_POST_ADD_NEW_MEDIA: "api/new-media/download"
};

export function GetApiBaseUrl(): string {
  return Constants.API_BASE_URL;
}

export function getWSBaseUrl(): string {
  return Constants.WS_BASE_URL;
}

