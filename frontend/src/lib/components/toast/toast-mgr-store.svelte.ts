import { Log } from "$lib/services/logger";
import type { Toast } from "./toast-models";

export class ToastManagerStore {
    #toastId = $state<number>(99);
    #log = new Log("ToastManagerStore");
    id: string;
    toasts: Array<Toast> = $state<Array<Toast>>([]);

    constructor(id: string) {
        this.id = id;
    }
    addToast(t: Toast) {
        t.id = ++this.#toastId;
        this.toasts.push(t);
        this.#log.info("Adding toast", t);
    }
    removeToast(t: Toast) {

        this.toasts = this.toasts.filter(x => x.id !== t.id);
    }
    removeAllToasts() {
        this.toasts = [];
    }
}
