import * as SecureStore from "expo-secure-store";
import AsyncStorage from "@react-native-async-storage/async-storage";
import CryptoJS from "react-native-crypto-js";

export const storeEncryptedLocalData = async (recordId: string, key: string, data: Record<any, any>) => {
    try {
        await SecureStore.setItemAsync(recordId, key);
        const encryptedData = CryptoJS.AES.encrypt(JSON.stringify(data), key).toString();
        await AsyncStorage.setItem(recordId, encryptedData);
    } catch (e) {
        throw e
    }
    return true;
}

export const getEncryptedLocalData = async (recordId: string) => {
    const key = await SecureStore.getItemAsync(recordId);
    if (key == null) {
        return;
    }
    const data = await AsyncStorage.getItem(recordId);
    if (data == null) {
        return;
    }
    const bytes = await CryptoJS.AES.decrypt(data, key);
    return await JSON.parse(bytes.toString(CryptoJS.enc.Utf8));
    // TODO: delete data, once extracted to save space. or only after online submit?
}

export const deleteEncryptedLocalData = async (recordId: string) => {
    const key = await SecureStore.getItemAsync(recordId);
    const key_1 = await AsyncStorage.removeItem(recordId);
    return await SecureStore.deleteItemAsync(recordId);
}
