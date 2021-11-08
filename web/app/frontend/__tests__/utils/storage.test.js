import AsyncStorage from '@react-native-async-storage/async-storage';
import { storeEncryptedLocalData, getEncryptedLocalData } from '../../src/utils/storage'

describe.skip("utils/storage", () => {
    it('checks if Async Storage is used', async () => {
        await AsyncStorage.getItem('testKey');
        expect(AsyncStorage.getItem).toBeCalledWith('testKey')
    })

    it("should encrypt data and store it", async () => {
        const data = { key: "value" };

        let succeeded = await storeEncryptedLocalData("test", "secretKey", data);
        expect(succeeded).toEqual(true);

        const stored = await getEncryptedLocalData("secretKey");
        expect(stored).toEqual(data);
    });
});
