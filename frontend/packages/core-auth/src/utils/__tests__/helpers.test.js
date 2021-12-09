import { handler, listener } from '../helpers';

describe('utils/helpers', () => {
  const browser = jest.fn().mockImplementation(() => ({
    // getHandle: jest.fn(),
    getPopup: jest.fn(),
    getListenerMap: jest.fn(),
    getHandle: jest.fn(() => 'handle'),
    dismissPopup: jest.fn(),
  }))();
  let windowSpy;

  describe('handler', () => {
    it('should call getPopup', () => {
      handler(() => {
        // noop
      }, browser);
      expect(browser.getPopup).toHaveBeenCalled();
    });
  });

  describe('listener', () => {
    const originalWindow = { ...window };
    const resolveMock = jest.fn();

    afterEach(() => {
      windowSpy.mockRestore();
      resolveMock.mockRestore();
    });

    it('should resolve', () => {
      const event = {
        data: { sender: browser.getHandle() },
        isTrusted: true,
        origin: 'http://localhost',
      };
      const getItemMock = jest.fn(() => 'handle');
      windowSpy = jest.spyOn(global, 'window', 'get');
      windowSpy.mockImplementation(() => ({
        ...originalWindow,
        sessionStorage: {
          getItem: getItemMock,
        },
      }));
      listener(event, browser, resolveMock);
      expect(resolveMock).toHaveBeenCalled();
    });

    it('should return if not trusted', () => {
      const event = {
        data: { sender: browser.getHandle() },
        isTrusted: false,
        origin: 'http://localhost',
      };
      listener(event, browser, resolveMock);
      expect(resolveMock).not.toHaveBeenCalled();
    });

    it('should return if origin not correct', () => {
      const event = {
        data: { sender: browser.getHandle() },
        isTrusted: true,
        origin: 'https://randomOrigin',
      };
      listener(event, browser, resolveMock);
      expect(resolveMock).not.toHaveBeenCalled();
    });
  });
});
