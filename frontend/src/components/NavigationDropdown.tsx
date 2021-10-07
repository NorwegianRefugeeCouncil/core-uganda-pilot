import React from 'react';
import { Appbar, Menu } from 'react-native-paper';
import pngMore from '../../assets/png/symbol_more.png';
import theme from '../constants/theme';

type DropdownProps = {
  visible: boolean,
  closeMenu: () => void
  openMenu: () => void
}

const NavigationDropdown: React.FC<DropdownProps> = ({ visible, closeMenu, openMenu }) => {
  return (
    <Menu
      visible={visible}
      onDismiss={closeMenu}
      anchor={
        <Appbar.Action
          icon={pngMore}
          onPress={openMenu}
          color={theme.colors.white}
        />
      }>
      <Menu.Item onPress={() => {
        console.log('Option 1 was pressed');
      }} title="Option 1" />
      <Menu.Item onPress={() => {
        console.log('Option 2 was pressed');
      }} title="Option 2" />
      <Menu.Item onPress={() => {
        console.log('Option 3 was pressed');
      }} title="Option 3" disabled />
    </Menu>
  );
};

export default NavigationDropdown;