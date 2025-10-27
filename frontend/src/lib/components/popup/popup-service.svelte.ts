import ConfirmPopup from './confirm-popup.svelte';
import TVEpisodePopup from './tv-episode-popup.svelte';
import CopyStatusPopup from './copy-status-popup.svelte';
import ConfirmRenameImpactPopup from './confirm-rename-impact-popup.svelte';
import { Log } from "$lib/services/logger";
import { Popup, PopupStore, PopupType } from "./popup-store.svelte";
import { WebSocketService } from '$lib/services/network';
import type { ConfirmMediaRequestItem, DestConfig } from '$lib/models';

export class PopupService {
  static #instance: PopupService | null = null;
  #popupStore: PopupStore | null = null;
  #log = new Log("PopupService");

  constructor() {
    if (PopupService.#instance === null) {
      PopupService.#instance = this;
    }
    return PopupService.#instance;
  }

  registerPopupStore(store: PopupStore) {
    this.#popupStore = store;
  }

  showConfirmation(confirmText: string, title: string = ""): Promise<boolean> {
    if (this.#popupStore) {
      return this.#popupStore.addPopup(new Popup(
        PopupType.ConfirmYesNo,
        ConfirmPopup,
        {
          bodyText: confirmText,
          ...(title?.length > 0 ? { title } : null)
        }
      ))
    } else {
      this.#log.error("Cannot show popup. PopupStore is not initialized");
      return Promise.reject("Cannot show popup. PopupStore is not initialized");
    }
  }

  showTVEpisodeEdit(season?: number, episode?: number): Promise<{ season: number, episode: number } | null> {
    if (this.#popupStore) {
      return this.#popupStore.addPopup(new Popup(
        PopupType.TVEpisodeEdit,
        TVEpisodePopup,
        {
          title: "Edit Episode Info",
          season,
          episode,
        }
      ));
    } else {
      this.#log.error("Cannot show popup. PopupStore is not initialized");
      return Promise.reject("Cannot show popup. PopupStore is not initialized");
    }
  }

  showFileTransferStatusPopup(): Promise<boolean> {
    if (this.#popupStore) {
      return this.#popupStore.addPopup(new Popup(
        PopupType.FileTransferProgress,
        CopyStatusPopup,
        {
          ws: new WebSocketService()
        }
      ));
    } else {
      this.#log.error("Cannot show popup. PopupStore is not initialized");
      return Promise.reject("Cannot show popup. PopupStore is not initialized");
    }
  }

  showConfirmRenameImpactPopup(selections: ConfirmMediaRequestItem[], allDestinations: DestConfig[]): Promise<boolean> {
    if (this.#popupStore) {
      return this.#popupStore.addPopup(new Popup(
        PopupType.ConfirmRenameImpact,
        ConfirmRenameImpactPopup,
        {
          selections,
          allDestinations
        }
      ));
    } else {
      this.#log.error("Cannot show popup. PopupStore is not initialized");
      return Promise.reject("Cannot show popup. PopupStore is not initialized");
    }
  }
}
