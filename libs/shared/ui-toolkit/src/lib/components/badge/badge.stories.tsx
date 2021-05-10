import React from 'react';
import { Button } from '../button/button.component';
import { Container } from '../container/container.component';
import { Badge } from './badge.component';

export default {
  title: 'Badge',
  decorators: [(Story: any) => <Story />],
};

export const basic = () => (
  <>
    <h1>
      Example heading <Badge theme="secondary">New</Badge>
    </h1>
    <h2>
      Example heading <Badge theme="secondary">New</Badge>
    </h2>
    <h3>
      Example heading <Badge theme="secondary">New</Badge>
    </h3>
    <h4>
      Example heading <Badge theme="secondary">New</Badge>
    </h4>
    <h5>
      Example heading <Badge theme="secondary">New</Badge>
    </h5>
    <h6>
      Example heading <Badge theme="secondary">New</Badge>
    </h6>
  </>
);

export const counter = () => (
  <Button>
    Notifications <Badge theme="secondary">4</Badge>
  </Button>
);

export const colors = () => (
  <>
    <Badge>Primary</Badge>
    <Badge theme="secondary">Secondary</Badge>
    <Badge theme="success">Success</Badge>
    <Badge theme="danger">Danger</Badge>
    <Badge theme="warning">Warning</Badge>
    <Badge theme="info">Info</Badge>
    <Badge theme="light">Light</Badge>
    <Badge theme="dark">Dark</Badge>
  </>
);

export const pills = () => (
  <>
    <Badge pill>Primary</Badge>
    <Badge pill theme="secondary">
      Secondary
    </Badge>
    <Badge pill theme="success">
      Success
    </Badge>
    <Badge pill theme="danger">
      Danger
    </Badge>
    <Badge pill theme="warning">
      Warning
    </Badge>
    <Badge pill theme="info">
      Info
    </Badge>
    <Badge pill theme="light">
      Light
    </Badge>
    <Badge pill theme="dark">
      Dark
    </Badge>
  </>
);
