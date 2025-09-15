export class Log {
    #ctx: string = "JMR";

    constructor(context: string) {
        this.#ctx = context;
    }

    info(...message: any[]) {
        const prefix = this.#composePrefix("INFO");
        console.log(prefix, ...message);
    }
    warn(...message: any[]) {
        const prefix = this.#composePrefix("WARN");
        console.warn(prefix, ...message);
    }
    error(...message: any[]) {
        const prefix = this.#composePrefix("ERROR");
        console.error(prefix, ...message);
    }

    #composePrefix(level: LogLevel): string {
        return `[${this.#ctx}][${level}]`;
    }
}

type LogLevel = "INFO" | "WARN" | "ERROR";
