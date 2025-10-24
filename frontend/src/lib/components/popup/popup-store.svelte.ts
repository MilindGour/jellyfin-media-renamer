export class PopupStore {
  popups: Array<Popup> = $state<Popup[]>([]);
  promiseMap: { [s: symbol]: { resolve: (v?: any) => void, reject: (reason: any) => void } } = {};

  addPopup(p: Popup): Promise<any> {
    p.id = Symbol("popup");
    const { promise, resolve, reject } = this.#createDeferredPromise();
    this.promiseMap[p.id] = { resolve, reject };
    this.popups.push(p);
    return promise;
  }
  removePopup(p: Popup, result: any) {
    this.promiseMap[p.id!].resolve(result);
    this.popups = this.popups.filter(x => x !== p);

    delete this.promiseMap[p.id!];
  }

  #createDeferredPromise(): { promise: Promise<any>, resolve: (v?: any) => void, reject: (reason: any) => void } {
    let resolve!: (val?: any) => void;
    let reject!: (reason?: any) => void;
    const promise = new Promise((res, rej) => {
      resolve = res;
      reject = rej;
    })
    return { promise, resolve, reject };
  }
}

export enum PopupType { ConfirmYesNo, TVEpisodeEdit, FileTransferProgress, ConfirmRenameImpact, Unknown };

export class Popup {
  id?: symbol;
  type: PopupType = PopupType.Unknown;
  data: any = null;
  component: any = null;

  constructor(type: PopupType, component: any, data: any) {
    this.type = type;
    this.component = component;
    this.data = data;
  }
}
