import type { ConfigSource } from "$lib/models/config-models";
import { getApiUrl } from "$lib/services/network";
import { Constants } from "$lib/stores/constants";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch }) => {
  const rest = await fetch(getApiUrl(Constants.API_GET_CONFIG_SOURCE));
  const configSources = await rest.json() as ConfigSource[];

  return {
    configSources
  }
}
