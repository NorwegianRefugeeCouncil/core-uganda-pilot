import React from 'react';
import { render, screen } from '@testing-library/react';

import AuthWrapper from "../AuthWrapper";

describe('AuthWrapper', () => {
    const TestInstance = (<AuthWrapper clientId={'clientId'} issuer={'issuer'}>
        <div>App</div>
    </AuthWrapper>);

    it('renders correctly', () => {
        render(TestInstance);
        expect(screen.getByText('Login')).toBeTruthy();
    });
})
