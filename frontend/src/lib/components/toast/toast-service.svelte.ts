import { Log } from "$lib/services/logger"
import type { Toast } from "./toast-models";
import type { ToastManagerStore } from "./toast-mgr-store.svelte";

export class ToastService {
    static #instance: ToastService | null = null;
    #mgr = $state<ToastManagerStore | null>(null);

    constructor() {
        if (ToastService.#instance === null) {
            ToastService.#instance = this;
        }
        return ToastService.#instance;
    }

    registerManager(toastManagerStore: ToastManagerStore) {
        this.#mgr = toastManagerStore;
    }

    show(t: Toast) {
        this.#mgr?.addToast(t);
    }
    closeAll() {
        this.#mgr?.removeAllToasts();
    }
}
