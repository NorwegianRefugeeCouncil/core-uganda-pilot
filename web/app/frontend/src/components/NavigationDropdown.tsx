import React from 'react';
import routes from '../constants/routes';
import {HamburgerIcon, Menu, Pressable, theme} from 'native-base';

type DropdownProps = {
    visible: boolean;
    closeMenu: () => void;
    openMenu: () => void;
    navigation: any;
};

const NavigationDropdown: React.FC<DropdownProps> = ({
    visible,
    closeMenu,
    openMenu,
    navigation,
}) => {
    return (
        <Menu
            isOpen={visible}
            onOpen={openMenu}
            onClose={closeMenu}
            trigger={triggerProps => (
                <Pressable
                    accessibilityLabel={'more options'}
                    {...triggerProps}
                >
                    <HamburgerIcon color={theme.colors.white}/>
                </Pressable>
            )}
        >
            <Menu.Item onPress={() => navigation.navigate(routes.cases.name)}>
                Cases
            </Menu.Item>
        </Menu>
    );
};

export default NavigationDropdown;
