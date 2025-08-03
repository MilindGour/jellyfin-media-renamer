export function formatPathString(path: string): string {
  if (path.startsWith('/')) {
    path = path.substring(1);
  }
  return path.split('/').join(' > ');
}

export function convertToSizeString(totalBytes: number): string {
  if (typeof totalBytes !== "number") {
    return totalBytes;
  }

  const allUnits = ["B", "KB", "MB", "GB", "TB"];
  let outTotal = totalBytes;
  let outUnitIndex = 0;

  while (outTotal > 1024) {
    outTotal /= 1024;
    outUnitIndex++;

    if (outUnitIndex === allUnits.length - 1) {
      break;
    }
  }

  return `${outTotal.toFixed(2)} ${allUnits[outUnitIndex]}`;
}

export function joinStrings(...parts: string[]): string {
  return parts.join("_");
}
