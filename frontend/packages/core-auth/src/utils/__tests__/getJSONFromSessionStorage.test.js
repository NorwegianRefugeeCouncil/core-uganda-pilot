import { getJSONFromSessionStorage } from '../getJSONFromSessionStorage';

describe('utils/getJSONFromSessionStorage', () => {
  let sessionStorageSpy;
  const sessionStorageOrig = { ...sessionStorage };

  const getItemMock = jest
    .fn()
    .mockReturnValueOnce(null)
    .mockReturnValueOnce('{"token": "handle"}')
    .mockReturnValueOnce('handle');

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
    expect(getJSONFromSessionStorage('key')).toEqual(undefined);
    expect(getItemMock).toHaveBeenCalled();
  });

  it('should return stored item', () => {
    expect(getJSONFromSessionStorage('key')).toEqual({ token: 'handle' });
  });

  it('should return undefined if stored value is not valid JSON', () => {
    expect(getJSONFromSessionStorage('key')).toEqual(undefined);
  });
});
