import type { DropdownModel } from "./dropdown-model.svelte";

export class DropdownService {
  private static instance: DropdownService | null = null;
  private static dropdowns: { [key: string]: DropdownModel } = {};

  constructor() {
    if (DropdownService.instance === null) {
      DropdownService.instance = this;
    }
    return DropdownService.instance;
  }

  static register(id: string, dropdownData: DropdownModel): boolean {
    if (id in DropdownService.dropdowns) {
      return false;
    }
    DropdownService.dropdowns[id] = dropdownData;
    return true;
  }
  static getValueOf(id: string): any {
    if (id in DropdownService.dropdowns) {
      return DropdownService.dropdowns[id].value;
    }
    return null;
  }

  static unregister(id: any) {
    if (id in DropdownService.dropdowns) {
      delete DropdownService.dropdowns[id];
    }
  }
}
