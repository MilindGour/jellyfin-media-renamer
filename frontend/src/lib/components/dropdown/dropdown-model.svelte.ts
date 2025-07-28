export class DropdownModel {
  _value = $state<any>(null);
  disabled = $state<boolean>(false);
  open = $state<boolean>(false);
  el = $state<HTMLElement | null>(null);
  options = $state<any[]>([]);
  label = $state<string>("");
  labelProperty = $state<string>("");

  get value() {
    return this._value;
  }
  set value(v: any) {
    if (v) {
      const jsonStrV = JSON.stringify(v);
      const targetIndex = this.options.findIndex(o => JSON.stringify(o) === jsonStrV);
      if (targetIndex > -1) {
        this._value = this.options[targetIndex];
        this.label = this.getSelectedItemLabel(this._value, this.labelProperty);
      }
    } else {
      this._value = null;
    }
  }

  constructor(initialOptions: DropdownInitOptions) {
    console.log("DropdownModel constructor called:", initialOptions);
    this.options = initialOptions.options;
    this.disabled = initialOptions.disabled || false;
    this.labelProperty = initialOptions.labelProperty || '';
    this.value = initialOptions.value || null;
  }

  private getSelectedItemLabel(item: any, labelProp: string): string {
    if (labelProp?.length > 0) {
      return item[labelProp];
    } else {
      if (typeof item === 'string') {
        return item;
      }
      if ('label' in item) {
        return item.label;
      }
    }
    return '<no label>';
  }
}

export type DropdownInitOptions = {
  options: any[],
  value?: any,
  disabled?: boolean,
  labelProperty?: string
};
