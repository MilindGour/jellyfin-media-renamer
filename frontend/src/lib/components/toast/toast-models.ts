import type { TransitionConfig } from "svelte/transition";
import { cubicOut } from "svelte/easing";
import type { ToastService } from "./toast-service.svelte";

export enum ToastType { INFO, WARNING, ERROR, SUCCESS };

export class Toast {
    id: number = -1;
    message: string;
    type: ToastType = ToastType.INFO;
    title: string = "";
    closeinMS: number = 3000;
    closeIcon: boolean = true;

    constructor(message: string) {
        this.message = message;
    }
}

export class ToastBuilder {
    #toast: Toast;

    constructor(message: string) {
        this.#toast = new Toast(message);
    }
    setTitle(title: string) {
        this.#toast.title = title;
        return this;
    }
    setCloseInMS(millis: number) {
        this.#toast.closeinMS = millis;
        return this;
    }
    setCloseIconVisible(visible: boolean) {
        this.#toast.closeIcon = visible;
        return this;
    }
    disableAutoClose() {
        this.#toast.closeinMS = -1;
        return this;
    }
    setMessage(message: string) {
        this.#toast.message = message;
        return this;
    }
    setType(newType: ToastType) {
        this.#toast.type = newType;
        return this;
    }
    build() {
        return this.#toast;
    }
}

export class ToastFactory {

    static createErrorToast(message: string): Toast {
        return (new ToastBuilder(message))
            .setTitle("Error")
            .setType(ToastType.ERROR)
            .setCloseInMS(5000)
            .build();
    }
    static createWarningToast(message: string): Toast {
        return (new ToastBuilder(message))
            .setTitle("Warning")
            .setType(ToastType.WARNING)
            .setCloseInMS(5000)
            .build();
    }
    static createSuccessToast(message: string, title?: string): Toast {
        return (new ToastBuilder(message))
            .setTitle(title || "Success")
            .setType(ToastType.SUCCESS)
            .setCloseInMS(5000)
            .build();
    }

}

export type ToastCloseReason = "CLOSE_BUTTON" | "CLOSE_TIMER" | "CLOSED_BY_MGR";
export type CloseFunctionHandler = (reason: ToastCloseReason) => void;

export function slideFromRight(node: Element): TransitionConfig {

    return {
        css(_, u) {
            return `translate: ${u * 100}%`;
        },
        duration: 250,
        easing: cubicOut,
    };
}
