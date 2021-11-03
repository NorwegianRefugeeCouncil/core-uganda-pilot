import * as SecureStore from "expo-secure-store";
import AsyncStorage from "@react-native-async-storage/async-storage";
import CryptoJS from "react-native-crypto-js";

export const storeEncryptedLocalData = (recordId: string, key: string, data: Record<any, any>) => {
    return SecureStore.setItemAsync(recordId, key)
        .then(async () => {
            const encryptedData = await CryptoJS.AES.encrypt(key, JSON.stringify(data));
            return AsyncStorage.setItem(recordId, encryptedData.toString())
        })
}

export const getEncryptedLocalData = (recordId: string) => {
    return SecureStore.getItemAsync(recordId)
        .then(async (key) => {
            if (key == null) {
                return;
            }

            const data = await AsyncStorage.getItem(recordId);
            if (data == null) {
                return;
            }
            const bytes = CryptoJS.AES.decrypt(key, data);
            return JSON.parse(bytes.toString());
            // TODO: delete data, once extracted to save space. or only after online submit?
        })
}

export const deleteEncryptedLocalData = (recordId: string) => {
    return SecureStore.getItemAsync(recordId)
        .then((key) => {
            return AsyncStorage.removeItem(recordId)
        })
        .then((key) => {
            return SecureStore.deleteItemAsync(recordId);
        })
}
