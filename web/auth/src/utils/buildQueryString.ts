import qs from "qs";

export function buildQueryString(input: Record<string, string>): string {
    return qs.stringify(input);
}
