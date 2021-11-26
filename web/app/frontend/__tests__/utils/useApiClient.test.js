import React from 'react'
import { useApiClient } from '../../src/utils/useApiClient';
import { ApiClient } from 'core-js-api-client';
import { create } from 'react-test-renderer'

describe('useApiClient', () => {
    let c;
    const TestComponent = () => {
        c = useApiClient();
        return (<></>)
    };
    test('returns an api-client', () => {
        create(<TestComponent />);
        expect(c).toBeInstanceOf(ApiClient);
    })
})
