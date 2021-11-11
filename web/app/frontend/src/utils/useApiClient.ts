import { client } from 'core-js-api-client';
import { useMemo } from 'react';

import host from '../constants/host';

export function useApiClient(): client {
    return useMemo(() => {
        return new client(host);
    }, []);
}
