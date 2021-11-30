import { ApiClient } from 'core-js-api-client';
import { useMemo } from 'react';

import host from '../constants/host';

export function useApiClient(): ApiClient {
    return useMemo(() => {
        return new ApiClient(host);
    }, []);
}
