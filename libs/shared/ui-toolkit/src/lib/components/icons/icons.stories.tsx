import * as React from 'react';
import Card from '../card/card.component';
import Icons from './icons';

export default {
  title: 'Icons',
  decorators: [(Story: any) => <Story />],
};

const allIcons = Object.values(Icons)
  .filter((icon) => {
    return {}.toString.call(icon) === '[object Function]';
  })
  .map((icon) =>
    React.createElement(icon, { style: { fontSize: '2em' }, className: 'me-2' })
  );

export const bootstrap = () => (
  <Card>
    <Card.Body>
      <Card.Title>Bootstrap Icons</Card.Title>
      {allIcons}
    </Card.Body>
  </Card>
);
