import AsyncStorage from '@react-native-async-storage/async-storage';
import * as SecureStore from 'expo-secure-store';
import CryptoJS from 'react-native-crypto-js';

export const storeEncryptedLocalData = async (
    recordId: string,
    key: string,
    data: Record<string, unknown>
) => {
    await SecureStore.setItemAsync(recordId, key);
    const encryptedData = (
        await CryptoJS.AES.encrypt(JSON.stringify(data), key)
    ).toString();
    await AsyncStorage.setItem(recordId, encryptedData);
};

export const getEncryptedLocalData = async (
    recordId: string
): Promise<Record<string, unknown>> => {
    const key = await SecureStore.getItemAsync(recordId);
    if (key == null) {
        return {};
    }
    const encrypted = await AsyncStorage.getItem(recordId);
    if (encrypted == null) {
        return {};
    }
    const decrypted: ArrayBuffer = await CryptoJS.AES.decrypt(encrypted, key);
    const string = decrypted.toString(CryptoJS.enc.Utf8);
    return JSON.parse(string);
    // TODO: delete data, once extracted to save space. or only after online submit?
};

export const deleteEncryptedLocalData = async (recordId: string) => {
    await AsyncStorage.removeItem(recordId);
    await SecureStore.deleteItemAsync(recordId);
};
