import React from 'react';
// import Renderer from 'react-test-renderer';

import { render } from '@testing-library/react';
import AuthWrapper from "../AuthWrapper";

describe('AuthWrapper', ()=>{
    const AuthWrapperTest = render(
        <AuthWrapper clientId={'clientId'} issuer={'issuer'}>
            <div>App</div>
        </AuthWrapper>
    )
    // expect(AuthWrapperTest).

})
