import resolveDiscoveryAsync from "../resolveDiscoveryAsync";
import axios from "axios";

jest.mock('axios', () => ({
    ...jest.requireActual('axios'),
    get: jest.fn()
}));

describe('utils/resolveDiscoveryAsync', () => {
    it('should return discovery with an issuer string', async () => {
        //@ts-ignore
        axios.get.mockResolvedValue({data: 'data'})
        const result = await resolveDiscoveryAsync('issuer');
        expect(axios.get).toHaveBeenCalled();
        expect(result).toEqual('data');
    })
    it('should return discovery with an issuer object', async () => {
        //@ts-ignore
        axios.get.mockResolvedValue({data: 'data'})
        const result = await resolveDiscoveryAsync({issuer: 'issuer'});
        expect(axios.get).toHaveBeenCalled();
        expect(result).toEqual('data');
    })
    it('should return error', async () => {
        //@ts-ignore
        axios.get.mockRejectedValue({error: 'error'})
        const result = await resolveDiscoveryAsync('issuer');
        expect(axios.get).toHaveBeenCalled();
        expect(result).toEqual(null);
    })
})
