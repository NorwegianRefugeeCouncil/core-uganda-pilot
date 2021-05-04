import { storiesOf } from '@storybook/react';
import React from 'react';
import { FieldInfo } from './fieldinfo.component'

storiesOf('Field Info', module)
.add('default', () => (
  <>
    <FieldInfo name={""} description={""}/>
  </>
));