import type { ConfigSource, ConfigSourceByID } from "$lib/models/config-models";
import { JNetworkClient } from "$lib/services/network";
import { Constants } from "./constants";

export class JmrStore {
  static #instance: null | JmrStore = null;
  #networkClient: JNetworkClient;

  // all private store variables
  #configSource: ConfigSource | null = $state(null);

  // all public store variables
  configSourceDetails: Promise<ConfigSourceByID> | null = $derived.by(() => {
    if (this.#configSource !== null) {
      const csId = this.#configSource.id;
      const configSourceDetailApi = Constants.API_GET_CONFIG_SOURCE_BY_ID.replace(":id", csId);
      return this.#networkClient.httpGetJSON<ConfigSourceByID>(configSourceDetailApi);
    }
    return null;
  });

  constructor(netClient: JNetworkClient) {
    if (JmrStore.#instance === null) {
      JmrStore.#instance = this;
    }
    this.#networkClient = netClient;
    return JmrStore.#instance;
  }

  setConfigSource(cs: ConfigSource | null) {
    this.#configSource = cs;
  }
}
