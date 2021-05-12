import React from 'react';
import Nav from './nav.component';

export default {
  title: 'Nav',
  decorators: [(Story: any) => <Story />],
};

export const basic = () => (
  <Nav>
    <Nav.Brand></Nav.Brand>
    <Nav.Collapse>
      <Nav.Item></Nav.Item>
    </Nav.Collapse>
  </Nav>
);
