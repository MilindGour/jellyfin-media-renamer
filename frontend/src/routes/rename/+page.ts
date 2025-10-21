import type { DestConfig } from "$lib/models";
import { getApiUrl } from "$lib/services/network";
import { Constants } from "$lib/stores/constants";
import type { PageLoad } from "../$types";

export const load: PageLoad = async ({ fetch }) => {
  const rest = await fetch(getApiUrl(Constants.API_GET_DESTINATIONS));
  const config = await rest.json() as DestConfig[];

  return {
    destinations: config
  }
}
