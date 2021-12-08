import {getQueryParams} from "../queryparams";

describe('utils/queryParams', () => {
    const baseUrl = 'www.url.de';

    it('should return url without parameters', () => {
        expect(getQueryParams(`${baseUrl}`)).toEqual({
            errorCode: null,
            params: {[baseUrl]: ''}
        })
    })
    it('should parse query parameters', () => {
        expect(getQueryParams(`${baseUrl}?id=2`)).toEqual({
            errorCode: null,
            params: {
                id: "2",
            }
        })
    })
    it('should find errorCode in params', () => {
        expect(getQueryParams(`${baseUrl}/?id=2&errorCode=error`)).toEqual({
            errorCode: 'error',
            params: {
                id: "2",
            }
        })
    })
    it('should handle hashes', () => {
        expect(getQueryParams(`${baseUrl}/?id=2&errorCode=error#hash`)).toEqual({
            errorCode: 'error',
            params: {
                id: "2",
                hash: ''
            }
        })
    })
})
