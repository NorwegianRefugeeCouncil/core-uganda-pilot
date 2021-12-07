import {buildQueryString} from "../buildQueryString";
import {stringify} from "qs";

jest.mock('qs', ()=>({
    stringify: jest.fn()
}))

describe('utils/buildQueryString', () => {
    it('should call qs.stringify', () => {
        buildQueryString({foo: 'bar'})
        expect(stringify).toHaveBeenCalled()
    })
})
