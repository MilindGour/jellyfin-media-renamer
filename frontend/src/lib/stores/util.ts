import type { AllowedExtensions, Config, Filetype } from "$lib/models/models";

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

    return `${outTotal.toFixed(2)} ${allUnits[outUnitIndex]}`;
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
