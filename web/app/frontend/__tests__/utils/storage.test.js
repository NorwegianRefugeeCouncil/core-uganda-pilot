import AsyncStorage from '@react-native-async-storage/async-storage';
import * as SecureStore from 'expo-secure-store';
import CryptoJS from 'react-native-crypto-js';

import {
    deleteEncryptedLocalData,
    getEncryptedLocalData,
    storeEncryptedLocalData,
} from '../../src/utils/storage';

const recordId = 'mock';
const key = 'whatever';
const data = { something: 'ok' };
SecureStore.getItemAsync.mockResolvedValue(key);
CryptoJS.AES.decrypt.mockResolvedValue(JSON.stringify(data));
let encryptedData;
(async () => {
    encryptedData = (
        await CryptoJS.AES.encrypt(JSON.stringify(data), key)
    ).toString();
    CryptoJS.AES.encrypt.mockResolvedValue(encryptedData);
    AsyncStorage.getItem.mockResolvedValue(encryptedData);
})();

describe('utils/storage', () => {
    describe(storeEncryptedLocalData.name, () => {
        it('should call internal fns with the correct params', async () => {
            const succeeded = await storeEncryptedLocalData(
                recordId,
                key,
                data
            );

            expect(succeeded).toBeTruthy();
            expect(SecureStore.setItemAsync).toHaveBeenCalledWith(
                recordId,
                key
            );
            expect(CryptoJS.AES.encrypt).toHaveBeenCalledWith(
                JSON.stringify(data),
                key
            );
            expect(AsyncStorage.setItem).toHaveBeenCalledWith(
                recordId,
                encryptedData
            );
        });
    });

    describe(getEncryptedLocalData.name, () => {
        it('should call internal fns with the correct params', async () => {
            const data = await getEncryptedLocalData(recordId);

            expect(data).toBeDefined();
            expect(SecureStore.getItemAsync).toHaveBeenCalledWith(recordId);
            expect(AsyncStorage.getItem).toHaveBeenCalledWith(recordId);
            expect(CryptoJS.AES.decrypt).toHaveBeenCalledWith(
                encryptedData,
                key
            );
        });
    });

    describe(deleteEncryptedLocalData.name, () => {
        it('should call internal fns with the correct params', async () => {
            const succeeded = await deleteEncryptedLocalData(recordId);

            expect(succeeded).toBeTruthy();
            expect(AsyncStorage.removeItem).toHaveBeenCalledWith(recordId);
            expect(SecureStore.deleteItemAsync).toHaveBeenCalledWith(recordId);
        });
    });
});
