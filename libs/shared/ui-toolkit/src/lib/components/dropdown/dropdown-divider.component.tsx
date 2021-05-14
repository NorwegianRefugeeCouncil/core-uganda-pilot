import * as React from 'react';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface DropdownDividerProps
  extends React.ComponentPropsWithoutRef<'hr'> {}

export const DropdownDivider: React.FC<DropdownDividerProps> = (props) => (
  <hr className="dropdown-divider" {...props} />
);

export default DropdownDivider;
