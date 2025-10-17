import type { AllowedExtensions, Config, Filetype } from "$lib/models";

export function formatPathString(path: string): string {
  if (path.startsWith('/')) {
    path = path.substring(1);
  }
  return path.split('/').join(' > ');
}
export function getRelativePath(fullPath: string, relativeTo: string): string {
  return fullPath.replace(relativeTo, '.');
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
  const outValue = Math.round(outTotal * 10) / 10;
  return `${outValue} ${allUnits[outUnitIndex]}`;
}

export function joinStrings(...parts: string[]): string {
  return parts.join("_");
}

export function getBasename(path: string): string {
  if (path.includes("/")) {
    const lastSlashIndex = path.lastIndexOf("/");
    return path.substring(lastSlashIndex + 1);
  }
  return path;
}
export function getSeasonEpisodeShortString(season: number, episode: number): string {
  return `S${padNumber(season, 2)}E${padNumber(episode, 2)}`;
}
export function padNumber(n: number, paddingSize: number): string {
  const zStr = "0".repeat(paddingSize) + n.toString();
  const fromIndex = zStr.length - paddingSize;
  return zStr.substring(fromIndex);
}

export function getFiletype(filePath: string, allowedExtensions: AllowedExtensions): Filetype {
  const isMediaExtension = allowedExtensions.media.some(e => filePath.endsWith(e));
  if (isMediaExtension) return "MEDIA";

  const isSubtitleExtension = allowedExtensions.subtitle.some(e => filePath.endsWith(e));
  if (isSubtitleExtension) return 'SUBTITLE';

  return 'UNKNOWN';
}

export function removeCommonSubstring(firstString: string, secondString: string): { first: string, second: string } {
  let first: string = "";
  let second: string = "";

  const minLength = Math.min(firstString.length, secondString.length);
  let i = 0;
  for (i = 0; i < minLength; i++) {
    if (firstString[i] !== secondString[i]) {
      break;
    }
  }
  first = firstString.substring(i);
  second = secondString.substring(i);

  return { first, second };
}

export function formatTimeString(timeString: string): string {
  if (!timeString) {
    return "invalid time";
  }
  const [hhStr, mmStr, ssStr] = timeString.split(":");
  const [hh, mm, ss] = [+hhStr, +mmStr, +ssStr];

  let output: string[] = [];
  if (hh > 0) {
    output.push(`${hh} hours`);
  }
  if (mm > 0) {
    output.push(`${mm} minutes`)
  }
  if (ss > 0) {
    output.push(`${ss} seconds`);
  }

  return output.join(" ");
}
