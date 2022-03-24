import { act, renderHook } from '@testing-library/react-hooks';

import { useAPICall } from '../useAPICall';

describe('rerunOnArgChange - true', () => {
  it('success', async () => {
    const mockData = 'mockData';
    const mockFunction = jest.fn((a: number, b: string) =>
      Promise.resolve(mockData),
    );
    const mockArgs = [1, 'a string'] as [number, string];
    const mockRerunOnArgChange = true;

    const { result, waitForNextUpdate } = renderHook(() =>
      useAPICall(mockFunction, mockArgs, mockRerunOnArgChange),
    );

    expect(mockFunction).toHaveBeenCalledTimes(1);
    expect(mockFunction).toHaveBeenCalledWith(...mockArgs);
    expect(result.current[1].loading).toBe(true);
    expect(result.current[1].error).toBe(null);
    expect(result.current[1].data).toBe(null);

    await waitForNextUpdate();

    expect(mockFunction).toHaveBeenCalledTimes(1);
    expect(mockFunction).toHaveBeenCalledWith(...mockArgs);
    expect(result.current[1].loading).toBe(false);
    expect(result.current[1].error).toBe(null);
    expect(result.current[1].data).toBe(mockData);

    // For some reason it skips over the initial state when asserting on "current" so here we check the full history
    const all = result.all as [() => Promise<void>, State<string>][];
    expect(all[0][1].loading).toBe(false);
    expect(all[0][1].error).toBe(null);
    expect(all[0][1].data).toBe(null);
    expect(all[1][1].loading).toBe(true);
    expect(all[1][1].error).toBe(null);
    expect(all[1][1].data).toBe(null);
    expect(all[2][1].loading).toBe(false);
    expect(all[2][1].error).toBe(null);
    expect(all[2][1].data).toBe(mockData);
  });

  it('failure', async () => {
    const mockFunction = jest.fn((a: number, b: string) =>
      Promise.reject(new Error('mock-error')),
    );
    const mockArgs = [1, 'a string'] as [number, string];
    const mockRerunOnArgChange = true;

    const { result, waitForNextUpdate } = renderHook(() =>
      useAPICall(mockFunction, mockArgs, mockRerunOnArgChange),
    );

    expect(mockFunction).toHaveBeenCalledTimes(1);
    expect(mockFunction).toHaveBeenCalledWith(...mockArgs);
    expect(result.current[1].loading).toBe(true);
    expect(result.current[1].error).toBe(null);
    expect(result.current[1].data).toBe(null);

    await waitForNextUpdate();

    expect(mockFunction).toHaveBeenCalledTimes(1);
    expect(mockFunction).toHaveBeenCalledWith(...mockArgs);
    expect(result.current[1].loading).toBe(false);
    expect(result.current[1].error).toBe('mock-error');
    expect(result.current[1].data).toBe(null);

    // For some reason it skips over the initial state when asserting on "current" so here we check the full history
    const all = result.all as [() => Promise<void>, State<string>][];
    expect(all[0][1].loading).toBe(false);
    expect(all[0][1].error).toBe(null);
    expect(all[0][1].data).toBe(null);
    expect(all[1][1].loading).toBe(true);
    expect(all[1][1].error).toBe(null);
    expect(all[1][1].data).toBe(null);
    expect(all[2][1].loading).toBe(false);
    expect(all[2][1].error).toBe('mock-error');
    expect(all[2][1].data).toBe(null);
  });
});

describe('rerunOnArgChange - false', () => {
  it('success', async () => {
    const mockData = 'mockData';
    const mockFunction = jest.fn((a: number, b: string) =>
      Promise.resolve(mockData),
    );
    const mockArgs = [1, 'a string'] as [number, string];
    const mockRerunOnArgChange = false;

    const { result, waitForNextUpdate } = renderHook(() =>
      useAPICall(mockFunction, mockArgs, mockRerunOnArgChange),
    );

    expect(mockFunction).not.toHaveBeenCalled();
    expect(result.current[1].loading).toBe(false);
    expect(result.current[1].error).toBe(null);
    expect(result.current[1].data).toBe(null);

    act(() => {
      result.current[0]();
    });

    await waitForNextUpdate();

    expect(mockFunction).toHaveBeenCalledTimes(1);
    expect(mockFunction).toHaveBeenCalledWith(...mockArgs);
    expect(result.current[1].loading).toBe(false);
    expect(result.current[1].error).toBe(null);
    expect(result.current[1].data).toBe(mockData);

    // For some reason it skips over the loading state when asserting on "current" so here we check the full history
    const all = result.all as [() => Promise<void>, State<string>][];
    expect(all[0][1].loading).toBe(false);
    expect(all[0][1].error).toBe(null);
    expect(all[0][1].data).toBe(null);
    expect(all[1][1].loading).toBe(true);
    expect(all[1][1].error).toBe(null);
    expect(all[1][1].data).toBe(null);
    expect(all[2][1].loading).toBe(false);
    expect(all[2][1].error).toBe(null);
    expect(all[2][1].data).toBe(mockData);
  });

  it('failure', async () => {
    const mockFunction = jest.fn((a: number, b: string) =>
      Promise.reject(new Error('mock-error')),
    );
    const mockArgs = [1, 'a string'] as [number, string];
    const mockRerunOnArgChange = false;

    const { result, waitForNextUpdate } = renderHook(() =>
      useAPICall(mockFunction, mockArgs, mockRerunOnArgChange),
    );

    expect(mockFunction).not.toHaveBeenCalled();
    expect(result.current[1].loading).toBe(false);
    expect(result.current[1].error).toBe(null);
    expect(result.current[1].data).toBe(null);

    act(() => {
      result.current[0]();
    });

    await waitForNextUpdate();

    expect(mockFunction).toHaveBeenCalledTimes(1);
    expect(mockFunction).toHaveBeenCalledWith(...mockArgs);
    expect(result.current[1].loading).toBe(false);
    expect(result.current[1].error).toBe('mock-error');
    expect(result.current[1].data).toBe(null);

    // For some reason it skips over the loading state when asserting on "current" so here we check the full history
    const all = result.all as [() => Promise<void>, State<string>][];
    expect(all[0][1].loading).toBe(false);
    expect(all[0][1].error).toBe(null);
    expect(all[0][1].data).toBe(null);
    expect(all[1][1].loading).toBe(true);
    expect(all[1][1].error).toBe(null);
    expect(all[1][1].data).toBe(null);
    expect(all[2][1].loading).toBe(false);
    expect(all[2][1].error).toBe('mock-error');
    expect(all[2][1].data).toBe(null);
  });
});
