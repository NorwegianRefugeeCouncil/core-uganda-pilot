export function omit<T extends { [key: string]: any }>(
    object: T,
    keyToRemove: keyof T
): Partial<T> {
    return Object.keys(object)
        .filter((key) => key !== keyToRemove)
        .reduce((result, key) => Object.assign(result, {[key]: object[key]}), {});
}
