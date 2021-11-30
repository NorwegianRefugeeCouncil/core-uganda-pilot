import React from 'react';
import {act, render} from '@testing-library/react';
import AuthWrapper from "../AuthWrapper";
import {ErrorBoundary} from "../ErrorBoundary";

import exchangeCodeAsync from "../../utils/exchangeCodeAsync";

jest.mock('axios', () => ({
    ...jest.requireActual('axios'),
    create: jest.fn(() => ({
        interceptors: {
            request: {
                use: jest.fn(),
                eject() {
                }
            }
        }
    }))
}));

jest.mock('../../types/browser', () => ({
    maybeCompleteAuthSession: jest.fn()
}))

jest.mock('../../hooks/useAuthRequest', () => {
    return jest.fn(() => {
        return [
            {
                codeVerifier: 'codeVerifier'
            },
            {
                type: 'success',
                params: {
                    code: 'code'
                }
            },
            jest.fn()
                // .mockRejectedValue('error')
                .mockResolvedValue({})]
    })
})

jest.mock('../../hooks/useDiscovery', () => {
    return jest.fn(() => {
        return Symbol('discovery')
    })
})

jest.mock('../../utils/exchangeCodeAsync', () => {
    return jest.fn()
        .mockResolvedValue({
            ...jest.requireActual('../../types/response'),
            shouldRefresh: jest.fn(() => {
                return false;
            }),
            refreshAsync: jest.fn()
                .mockResolvedValue({
                    accessToken: 'accessToken',
                    tokenType: 'Bearer',
                    issuedAt: 3
                })
                .mockRejectedValue(undefined),
        })
    // .mockRejectedValue({error: 'error'})
})

describe.skip('Components: AuthWrapper', () => {
    beforeEach(() => {
        maybeCompleteAuthSession.mockClear();
        exchangeCodeAsync.mockClear();
    })

    it('should render the default login button', () => {
        const {getByRole} = render(
            <ErrorBoundary>
                <AuthWrapper clientId={'clientId'} issuer={'issuer'}/>
            </ErrorBoundary>
        )
        expect(maybeCompleteAuthSession).toHaveBeenCalledTimes(1);
        expect(getByRole('button')).toHaveTextContent('Login')
    })

    it('should render a custom login button', () => {
        const {getByRole} = render(
            <ErrorBoundary>
                <AuthWrapper
                    clientId={'clientId'}
                    issuer={'issuer'}
                    customLoginComponent={() => <div role={'button'}>Custom Login</div>}
                />
            </ErrorBoundary>
        )
        expect(maybeCompleteAuthSession).toHaveBeenCalledTimes(1);
        expect(getByRole('button')).toHaveTextContent('Custom Login')
    })

    it('should render the app content when logged in', () => {
        const {getByText, container, debug} = render(
            <ErrorBoundary>
                <AuthWrapper clientId={'clientId'} issuer={'issuer'}>
                    <div>App</div>
                </AuthWrapper>
            </ErrorBoundary>
        )
        expect(maybeCompleteAuthSession).toHaveBeenCalledTimes(1);
        expect(exchangeCodeAsync).toHaveBeenCalledTimes(0);
        const button = container.querySelector('button');

        act(() => {
            button.dispatchEvent(new MouseEvent('click', {bubbles: true}));
        });
        debug()
        expect(exchangeCodeAsync).toHaveBeenCalledTimes(1);

        expect(getByText('App')).toBeDefined()
    })

})
