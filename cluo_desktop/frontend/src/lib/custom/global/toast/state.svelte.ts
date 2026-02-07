import { getContext, onDestroy, setContext } from "svelte"
import { type ToastLevel, type Toast } from "./type"

export class ToastState {
    toasts = $state<Toast[]>([])
    toastToTimeoutMap = new Map<string, number>()

    constructor() {
        onDestroy(() => {
            for (const timeout of this.toastToTimeoutMap.values()) {
                clearTimeout(timeout)
            }
            this.toastToTimeoutMap.clear()
        })
    }

    add(
        level: ToastLevel,
        title: string,
        message: string,
        durationMs: number = 3000,
    ) {
        const id = crypto.randomUUID()
        this.toasts.push({
            id,
            level,
            title,
            message,
        })
        this.toastToTimeoutMap.set(
            id,
            setTimeout(() => {
                this.remove(id)
            }, durationMs)
        )
    }

    remove(id: string) {
        const timeout = this.toastToTimeoutMap.get(id)
        if (timeout) {
            clearTimeout(timeout)
            this.toastToTimeoutMap.delete(id)
        }
        this.toasts = this.toasts.filter((toast) => toast.id !== id)
    }
}

const TOAST_KEY = Symbol("toast")

export function setToastContext() {
    return setContext(TOAST_KEY, new ToastState())
}

export function getToastContext() {
    return getContext<ReturnType<typeof setToastContext>>(TOAST_KEY)
}
