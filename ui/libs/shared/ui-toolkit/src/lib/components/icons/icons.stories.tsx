import * as React from 'react';
import { Card } from '../card/card.component';
import { Icons } from './icons';

const allIcons = Object.entries(Icons)
  .filter(
    ([name, icon]) => Object().toString.call(icon) === '[object Function]'
  )
  .map(([name, icon]) => (
    <div className={'h1 m-4 d-inline-block'}>{React.createElement(icon)}</div>
  ));

export default {
  title: 'Icons',
  component: Icons,
};

export const Bootstrap = () => (
  <Card>
    <Card.Body>
      <Card.Title>Bootstrap Icons</Card.Title>
      {allIcons}
    </Card.Body>
  </Card>
);
