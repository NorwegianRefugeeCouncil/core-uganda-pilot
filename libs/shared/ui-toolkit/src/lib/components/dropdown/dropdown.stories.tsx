import { Dropdown, DropdownMenu, DropdownMenuItem } from './dropdown.component';

export default {
  title: 'Dropdown',
  decorators: [(Story: any) => <Story />],
};

export const basic = () => (
  <Dropdown label="Dropdown button">
    <DropdownMenu>
      <DropdownMenuItem label="Option 1" />
      <DropdownMenuItem label="Option 2" />
      <DropdownMenuItem label="Option 3" />
      <DropdownMenuItem label="Option 4" />
    </DropdownMenu>
  </Dropdown>
);
