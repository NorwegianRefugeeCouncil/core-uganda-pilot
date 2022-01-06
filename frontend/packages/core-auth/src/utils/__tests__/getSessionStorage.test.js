import { getSessionStorage } from '../getSessionStorage';

describe('utils/getSessionStorage', () => {
  let sessionStorageSpy;
  const sessionStorageOrig = { ...sessionStorage };

  const getItemMock = jest
    .fn()
    .mockReturnValueOnce(null)
    .mockReturnValueOnce('{"token": "handle"}');

  beforeAll(() => {
    sessionStorageSpy = jest.spyOn(global, 'sessionStorage', 'get');
    sessionStorageSpy.mockImplementation(() => ({
      ...sessionStorageOrig,
      getItem: getItemMock,
    }));
    sessionStorageSpy.getItem = getItemMock;
  });

  afterAll(() => {
    sessionStorageSpy.mockRestore();
  });

  it('should return undefined when no item with key is found', () => {
    expect(getSessionStorage('key')).toEqual(undefined);
    expect(getItemMock).toHaveBeenCalled();
  });

  it('should return stored item', () => {
    expect(getSessionStorage('key')).toEqual({ token: 'handle' });
  });
});
