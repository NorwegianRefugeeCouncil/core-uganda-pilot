import exchangeCodeAsync from "../exchangeCodeAsync";
import {AccessTokenRequest} from "../../types/request";

jest.mock('../../types/request', () => ({
    AccessTokenRequest: jest.fn(() => ({
        performAsync: jest.fn(() => 'TokenResponse'),
    }))
}));

describe('utils/exchangeCodeAsync', () => {
    it('should return TokenResponse', async () => {
        const result = await exchangeCodeAsync(
            {
                code: 'bar',
                redirectUri: '',
                clientId: 'clientId'
            },
            {
                token_endpoint: ''
            });
        expect(result).toEqual('TokenResponse');
        expect(AccessTokenRequest).toHaveBeenCalledWith({code: 'bar', redirectUri: '', clientId: 'clientId'});
    })
})
