import {getEncryptionKey} from "../../src/utils/getEncryptionKey";

it('should return a unique random string', () => {
    const strs = []
    for (let i = 0; i < 10000; i++) {
        const str = getEncryptionKey();
        const isUnique = strs.every(s => s !== str);
        strs.push(str);
        expect(isUnique).toBeTruthy();
    }
})
