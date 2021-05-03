import { Fragment, useState } from 'react';
import { storiesOf } from '@storybook/react';
import Fade from './fade.component';
import { Button } from '../button/button.component';

storiesOf('Fade', module)
  .add('default', () => {

    const [isIn, setIn] = useState(false);

    return (
      <Fragment>
        <div className={'container'}>
          <div className={'row'}>
            <div className={'col-12 mb-2'}>
              <Button onClick={() => setIn(!!isIn)}>Toggle</Button>
              <Fade   in={isIn} >
                <p>hello!</p>
              </Fade>
            </div>
          </div>
        </div>
      </Fragment>
    );
  })
;
