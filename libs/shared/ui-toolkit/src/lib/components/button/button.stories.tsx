import { Fragment } from 'react';
import { storiesOf } from '@storybook/react';
import { Button, CloseButton } from './button.component';
import { Card, CardBody, CardTitle } from '../card/card.component';

storiesOf('Input', module).add('default', () => (
  <Card>
    <CardBody>
      <CardTitle>Buttons</CardTitle>
      <Button>Primary</Button>
      <Button kind="secondary">Secondary</Button>
    </CardBody>
  </Card>
));
