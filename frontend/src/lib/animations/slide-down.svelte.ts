import { cubicInOut } from "svelte/easing";
export function slideDown(node: HTMLElement) {
  const existingHeight = parseInt(getComputedStyle(node).height);

  return {
    delay: 0,
    duration: 150,
    easing: cubicInOut,
    css: (t: number) => `height: ${existingHeight * t}px`
  }
}
