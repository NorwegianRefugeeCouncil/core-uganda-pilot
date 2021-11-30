import qs from 'qs';

export function getQueryParams(url: string): {
    errorCode: string | null;
    params: { [key: string]: string };
} {
    const parts = url.split('#');
    const hash = parts[1];
    const partsWithoutHash = parts[0].split('?');
    const queryString = partsWithoutHash[partsWithoutHash.length - 1];

    // Get query string (?hello=world)
    const parsedSearch = qs.parse(queryString, {parseArrays: false});

    // Pull errorCode off of params
    const errorCode = (parsedSearch.errorCode ?? null) as string | null;
    delete parsedSearch.errorCode;

    // Get hash (#abc=example)
    let parsedHash = {};
    if (parts[1]) {
        parsedHash = qs.parse(hash);
    }

    // Merge search and hash
    const params = {
        ...parsedSearch,
        ...parsedHash,
    };

    return {
        errorCode,
        params,
    };
}
