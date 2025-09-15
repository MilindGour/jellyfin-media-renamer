class LoaderService {
    static instance?: LoaderService;

    serviceCounter = $state<number>(0);
    visible = $derived(this.serviceCounter > 0);

    constructor() {
        if (!LoaderService.instance) {
            LoaderService.instance = this;
        }
        return LoaderService.instance;
    }
    addAPICounter() {
        this.serviceCounter = this.serviceCounter + 1;
    }
    subtractAPICounter() {
        if (this.serviceCounter > 0) {
            this.serviceCounter = this.serviceCounter - 1;
        }
    }
}

export const loaderService = new LoaderService();
