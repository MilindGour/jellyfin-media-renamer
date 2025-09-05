import type { SourcesResponse } from "$lib/models/models";
import { getApiUrl } from "$lib/services/network";
import { Constants } from "$lib/stores/constants";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch }) => {
    const rest = await fetch(getApiUrl(Constants.API_GET_SOURCE));
    const sourcesResponse = await rest.json() as SourcesResponse;

    return {
        sourcesResponse
    }
}
