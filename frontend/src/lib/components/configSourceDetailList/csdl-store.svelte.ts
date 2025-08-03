import type { DirEntry } from "$lib/models/config-models";

export class CSDLStore {
  items: DirEntry[] = $state([]);
  selectedItems: DirEntry[] = $state([]);
  selectAllChecked: boolean = $derived(this.items.length === this.selectedItems.length);
  selectAllIndeterminate: boolean = $derived.by(() => {
    if (this.selectedItems.length === 0 || this.items.length === this.selectedItems.length) {
      return false;
    }
    return true;
  })

  constructor(initialItems: DirEntry[]) {
    this.items = initialItems;
    this.selectedItems = [];
  }

  isItemSelected(item: DirEntry): boolean {
    return this.getSelectedItemIndex(item) > -1;
  }
  getSelectedItemIndex(item: DirEntry): number {
    const itemJson = JSON.stringify(item);
    const existingIndex = this.selectedItems.findIndex(i => JSON.stringify(i) === itemJson);
    return existingIndex;
  }

  // Inserts an item into selectedItems if not present already
  // returns true if the item was inserted, false if already present.
  selectItem(item: DirEntry): boolean {
    if (this.isItemSelected(item) === false) {
      this.selectedItems.push(item);
      return true;
    }
    return false;
  }

  // Deletes an item fron selectedItems if present in the selectedItems.
  // returns true if the item was deleted, false if it was not found.
  unselectItem(item: DirEntry): boolean {
    const i = this.getSelectedItemIndex(item);
    if (i > -1) {
      this.selectedItems.splice(i, 1);
      return true;
    }
    return false;
  }

  selectAll() {
    this.selectedItems = [...this.items];
  }
  unselectAll() {
    this.selectedItems = [];
  }
}
