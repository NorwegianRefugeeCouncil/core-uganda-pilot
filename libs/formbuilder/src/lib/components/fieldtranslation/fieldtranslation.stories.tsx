import { storiesOf } from '@storybook/react';
import React from 'react';
import { FieldTranslation } from './fieldtranslation.component'

storiesOf('Field Translation', module)
.add('default', () => (
  <>
    <FieldTranslation name="" translation={
        {}
    }/>
  </>
));