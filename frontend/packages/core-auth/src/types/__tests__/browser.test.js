import Browser from '../browser';

jest.mock('../../utils/helpers', () => ({
  handler: jest.fn(() => ({
    type: 'success',
    message: 'message',
  })),
}));

describe('utils/browser', () => {
  const browser = new Browser();
  let windowSpy;
  const originalWindow = { ...window };

  beforeEach(() => {
    windowSpy = jest.spyOn(global, 'window', 'get');
  });

  afterEach(() => {
    windowSpy.mockRestore();
  });

  describe('browser', () => {
    expect(typeof browser.dismissPopup).toBe('function');
    expect(typeof browser.maybeCompleteAuthSession).toBe('function');
    expect(typeof browser.openAuthSessionAsync).toBe('function');
    expect(typeof browser.generateStateAsync).toBe('function');
    expect(typeof browser.isCryptoAvailable).toBe('function');
    expect(typeof browser.isSubtleCryptoAvailable).toBe('function');
    expect(typeof browser.bufferToString).toBe('function');
    expect(typeof browser.generateRandom).toBe('function');
    expect(typeof browser.getStateFromUrlOrGenerateAsync).toBe('function');
  });

  describe('maybeCompleteAuthSession', () => {
    it('should return failed when no handle exists', () => {
      const result = browser.maybeCompleteAuthSession();
      expect(result.type).toEqual('failed');
    });

    it('should return success', () => {
      windowSpy.mockImplementation(() => ({
        ...originalWindow,
        location: {
          href: 'https://example.com',
        },
        sessionStorage: {
          getItem: jest.fn(() => 'handle'),
        },
        parent: {
          postMessage: jest.fn(),
        },
      }));
      const result = browser.maybeCompleteAuthSession();
      expect(result.type).toEqual('success');
    });

    it('should throw error when no parent exists', () => {
      windowSpy.mockImplementation(() => ({
        location: {
          href: 'https://example.com',
        },
        sessionStorage: {
          getItem: jest.fn(() => 'handle'),
        },
      }));
      expect(() => browser.maybeCompleteAuthSession()).toThrowError();
    });
  });

  describe('createPopup', () => {
    let windowFocusMock;
    let windowOpenMock;

    beforeEach(() => {
      windowSpy.mockImplementation(() => ({
        ...originalWindow,
        open: windowOpenMock,
      }));
    });
    afterEach(() => {
      jest.restoreAllMocks();
    });

    it('should catch error from focus', () => {
      expect(browser.getPopup()).toBeNull();
      windowFocusMock = jest.fn(() => {
        throw new Error();
      });
      windowOpenMock = jest.fn(() => ({
        focus: windowFocusMock,
      }));
      const wrapperFunction = () => browser.createPopup('url');
      expect(wrapperFunction).not.toThrow();
      expect(browser.getPopup()).not.toBeNull();
      expect(windowOpenMock).toHaveBeenCalled();
      expect(windowFocusMock).toHaveBeenCalled();
    });

    it("should throw error if popupWindow wasn't set", () => {
      windowOpenMock = jest.fn(() => undefined);

      const wrapperFunction = () => browser.createPopup('url');
      expect(wrapperFunction).toThrow('ERR_WEB_BROWSER_LOCKED: Popup window was blocked by the browser of failed to open.');
      expect(browser.getPopup()).not.toBeNull();
      expect(windowOpenMock).toHaveBeenCalled();
    });

    it('should set popup window', () => {
      windowFocusMock = jest.fn();
      windowOpenMock = jest.fn(() => ({
        focus: windowFocusMock,
      }));
      browser.createPopup('url');
      expect(browser.getPopup()).not.toBeNull();
      expect(windowOpenMock).toHaveBeenCalled();
      expect(windowFocusMock).toHaveBeenCalled();
    });
  });

  describe('openAuthSessionAsync', () => {
    let getStateFromUrlOrGenerateAsyncMock;
    const windowSessionSetItemMock = jest.fn(() => 'handle');

    let createPopupSpy;
    let popupSpy;

    beforeEach(() => {
      windowSpy.mockImplementation(() => ({
        ...originalWindow,
        sessionStorage: {
          setItem: windowSessionSetItemMock,
        },
      }));

      getStateFromUrlOrGenerateAsyncMock = jest
        .spyOn(browser, 'getStateFromUrlOrGenerateAsync')
        .mockResolvedValue('stateFromUrl');

      createPopupSpy = jest.spyOn(browser, 'createPopup');
      popupSpy = jest.spyOn(browser, 'getPopup');
    });

    afterEach(() => {
      jest.restoreAllMocks();
    });

    it('should call functions', async () => {
      const openAuthSessionAsyncPromise = browser.openAuthSessionAsync({
        url: 'https://localhost:8080',
      });
      expect(openAuthSessionAsyncPromise).toHaveProperty('then');
    });
  });

  describe.skip('dismissPopup', () => {
    beforeEach(() => {
      windowSpy.mockImplementation(() => ({
        ...originalWindow,
        addEventListener: jest.fn(),
        sessionStorage: {
          getItem: jest.fn(() => 'handle'),
          setItem: jest.fn(() => 'handle'),
          removeItem: jest.fn(),
        },
        crypto: {
          subtle: {
            digest: jest.fn().mockResolvedValue('hashedData'),
          },
          getRandomValues: jest.fn(() => 'randomValues'),
        },
        open: jest.fn(() => ({
          focus: jest.fn(),
          close: jest.fn(),
        })),
      }));
    });
    it('should return when popup window is null', () => {
      expect(browser.dismissPopup()).toBeUndefined();
    });
    it('should close popup', async () => {
      browser.createPopup('url');
      await browser.openAuthSessionAsync({ url: 'https://localhost:8080' });
      browser.dismissPopup();
      expect(window.removeEventListener).toHaveBeenCalled();
      expect(window.open).toHaveBeenCalled();
    });
  });

  describe('getStateFromUrlOrGenerateAsync', () => {
    let spy;
    beforeEach(() => {
      windowSpy.mockImplementation(() => ({
        crypto: {
          subtle: {
            digest: jest.fn().mockResolvedValue('hashedData'),
          },
          getRandomValues: jest.fn(() => 'randomValues'),
        },
      }));
      spy = jest.spyOn(browser, 'generateStateAsync').mockResolvedValue('generatedState');
    });
    afterEach(() => {
      spy.mockRestore();
    });

    it('should return generateStateAsync when no state in url', async () => {
      const result = await browser.getStateFromUrlOrGenerateAsync('https://www.google.com');
      expect(result).toEqual('generatedState');
      expect(spy).toHaveBeenCalledTimes(1);
    });

    it('should return state if present in url', async () => {
      const result = await browser.getStateFromUrlOrGenerateAsync('https://www.google.com?state=state');
      expect(spy).not.toHaveBeenCalled();
      expect(result).toEqual('state');
    });
  });

  describe('generateStateAsync', () => {
    let generateRandomSpy;
    let isSubtleCryptoAvailableSpy;

    beforeEach(() => {
      windowSpy.mockImplementation(() => ({
        crypto: {
          subtle: {
            digest: jest.fn().mockResolvedValue(new ArrayBuffer(4)),
          },
        },
      }));
      isSubtleCryptoAvailableSpy = jest.spyOn(browser, 'isSubtleCryptoAvailable');
      generateRandomSpy = jest.spyOn(browser, 'generateRandom').mockImplementation(() => 'random');
    });
    afterEach(() => {
      isSubtleCryptoAvailableSpy.mockRestore();
      generateRandomSpy.mockRestore();
    });

    it('should throw error when crypto.subtle is not available', async () => {
      isSubtleCryptoAvailableSpy.mockImplementation(() => false);
      await expect(browser.generateStateAsync()).rejects.toThrow(
        'ERR_WEB_BROWSER_CRYPTO: The current environment does not support Crypto',
      );
    });

    it('should generate state', async () => {
      isSubtleCryptoAvailableSpy.mockImplementation(() => true);
      const result = await browser.generateStateAsync();
      expect(result).toEqual('AAAAAA==');
    });
  });

  describe('generateRandom', () => {
    let isCryptoAvailableSpy;
    let bufferToStringSpy;
    const getRandomValuesSpy = jest.fn(() => 'randomValues');

    beforeEach(() => {
      isCryptoAvailableSpy = jest.spyOn(browser, 'isCryptoAvailable');
      bufferToStringSpy = jest.spyOn(browser, 'bufferToString');
      windowSpy.mockImplementation(() => ({
        crypto: {
          getRandomValues: getRandomValuesSpy,
        },
      }));
    });
    afterEach(() => {
      isCryptoAvailableSpy.mockRestore();
      bufferToStringSpy.mockRestore();
    });

    it('should generate random values with Math.random if crypto not available', () => {
      isCryptoAvailableSpy.mockImplementation(() => false);
      const result = browser.generateRandom(5);
      expect(result).toHaveLength(5);
      expect(isCryptoAvailableSpy).toHaveBeenCalledTimes(1);
      expect(getRandomValuesSpy).toHaveBeenCalledTimes(0);
    });

    it('should generate random values with crypto if available', () => {
      isCryptoAvailableSpy.mockImplementation(() => true);
      const result = browser.generateRandom(5);
      expect(result).toEqual('AAAAA');
      expect(result).toHaveLength(5);
      expect(isCryptoAvailableSpy).toHaveBeenCalledTimes(1);
      expect(getRandomValuesSpy).toHaveBeenCalledTimes(1);
    });
  });

  describe('bufferToString', () => {
    it('should turn buffer to string', () => {
      const buffer = new Uint8Array(3);
      expect(browser.bufferToString(buffer)).toEqual('AAA');
    });
  });

  describe('isCryptoAvailable', () => {
    it('should return false if window.crypto is not available', () => {
      expect(browser.isCryptoAvailable()).toBeFalsy();
    });

    it('should return true if window.crypto is available', () => {
      windowSpy.mockImplementation(() => ({
        crypto: {},
      }));
      expect(browser.isCryptoAvailable()).toBeTruthy();
    });
  });

  describe('isSubtleCryptoAvailable', () => {
    let isCryptoAvailableSpy;

    beforeEach(() => {
      isCryptoAvailableSpy = jest.spyOn(browser, 'isCryptoAvailable');
    });
    afterEach(() => {
      isCryptoAvailableSpy.mockRestore();
    });

    it('should return false if window.crypto.subtle is not available', () => {
      expect(browser.isSubtleCryptoAvailable()).toBeFalsy();
      expect(isCryptoAvailableSpy).toHaveBeenCalledTimes(1);
    });

    it('should return true if window.crypto.subtle is available', () => {
      windowSpy.mockImplementation(() => ({
        crypto: { subtle: {} },
      }));
      expect(browser.isSubtleCryptoAvailable()).toBeTruthy();
      expect(isCryptoAvailableSpy).toHaveBeenCalledTimes(1);
    });
  });
});
