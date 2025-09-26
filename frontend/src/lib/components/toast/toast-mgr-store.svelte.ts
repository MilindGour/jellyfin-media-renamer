import { Log } from "$lib/services/logger";
import type { Toast } from "./toast-models";

export class ToastManagerStore {
    #toastId = $state<number>(99);
    id: string;
    toasts: Array<Toast> = $state<Array<Toast>>([]);

    constructor(id: string) {
        this.id = id;
    }
    addToast(t: Toast) {
        t.id = ++this.#toastId;
        this.toasts.push(t);
    }
    removeToast(t: Toast) {

        this.toasts = this.toasts.filter(x => x.id !== t.id);
    }
    removeAllToasts() {
        this.toasts = [];
    }
}
