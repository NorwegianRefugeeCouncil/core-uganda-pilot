import React from 'react';

import './top-bar.css';
import { NavBar } from '@core/ui-toolkit';

/* eslint-disable-next-line */
export interface TopBarProps {}

export function TopBar(props: TopBarProps) {
  return (
    <Navbar>
      <h1>Welcome to top-bar!</h1>
    </Navbar>
  );
}

export default TopBar;
