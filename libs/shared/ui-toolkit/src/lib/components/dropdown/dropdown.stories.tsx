import Dropdown from './dropdown.component';

export default {
  title: 'Dropdown',
  component: Dropdown,
};

export const basic = () => (
  <Dropdown label="Dropdown button (todo)">
    <Dropdown.Menu>
      <Dropdown.Item label="Option 1" />
      <Dropdown.Item label="Option 2" />
      <Dropdown.Item label="Option 4" />
      <Dropdown.Item label="Option 3" />
    </Dropdown.Menu>
  </Dropdown>
);
