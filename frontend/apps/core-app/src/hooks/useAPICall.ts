import * as React from 'react';

type Parameter<T> = T extends (...args: infer T) => any ? T : never;
type AsyncReturnType<T extends (...args: any) => any> = T extends (
  ...args: any
) => Promise<infer U>
  ? U
  : T extends (...args: any) => infer U
  ? U
  : any;

type State<T> = {
  loading: boolean;
  error: string | null;
  data: T | null;
};

type Action<T> =
  | {
      type: 'INIT';
    }
  | {
      type: 'SUCCESS';
      payload: T;
    }
  | {
      type: 'ERROR';
      payload: string;
    };

const initialState = {
  loading: false,
  error: null,
  data: null,
};

const createDataFetchReducer =
  <T>() =>
  (state: State<T>, action: Action<T>): State<T> => {
    switch (action.type) {
      case 'INIT':
        return {
          ...state,
          loading: true,
        };
      case 'SUCCESS':
        return {
          ...state,
          loading: false,
          error: null,
          data: action.payload,
        };
      case 'ERROR':
        return {
          ...state,
          loading: false,
          error: action.payload,
          data: null,
        };
      default:
        throw new Error();
    }
  };

export const useAPICall = <F extends (...arg: Parameter<F>) => Promise<any>>(
  apiFunction: F,
  args: Parameter<F>,
  rerunOnArgChange: boolean,
): [() => Promise<void>, State<AsyncReturnType<F>>] => {
  const reducer = createDataFetchReducer<AsyncReturnType<F>>();
  const [state, dispatch] = React.useReducer(reducer, initialState);

  const run = React.useCallback(async () => {
    dispatch({ type: 'INIT' });
    try {
      const data = await apiFunction(...args);
      dispatch({ type: 'SUCCESS', payload: data });
    } catch (error: any) {
      dispatch({ type: 'ERROR', payload: error?.message ?? 'Unknown error' });
    }
  }, [JSON.stringify(args)]);

  React.useEffect(() => {
    if (rerunOnArgChange) {
      run();
    }
  }, [JSON.stringify(args)]);

  return [run, state];
};
