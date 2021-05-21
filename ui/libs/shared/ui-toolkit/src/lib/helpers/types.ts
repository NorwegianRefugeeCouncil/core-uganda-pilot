import * as React from 'react';

export type Theme = 'light' | 'dark';

export type Color =
  | 'primary'
  | 'secondary'
  | 'success'
  | 'danger'
  | 'warning'
  | 'info'
  | 'light'
  | 'dark';

export type Size = 'sm' | 'md' | 'lg' | 'xl' | 'xxl';

export type Direction = 'down' | 'up' | 'right' | 'left' | 'start' | 'end';

export type BsInputTypes =
  | 'text'
  | 'textarea'
  | 'email'
  | 'password'
  | 'file'
  | 'color';

export type NonBsInputTypes = 'number' | 'date' | 'datetime' | 'time' | 'tel';

export type EventKey = string | number;

export type SelectCallback = (
  eventKey: EventKey,
  e: React.SyntheticEvent
) => void;
