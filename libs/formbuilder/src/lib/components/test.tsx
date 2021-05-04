import React, { FunctionComponent } from 'react';

type TestProps = {
    test?: string
}

export const Test: FunctionComponent<TestProps> = (props) => {
    return (
        <h1>Hello World</h1>
    )
}