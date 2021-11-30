import { StackHeaderProps } from '@react-navigation/stack';
import React from 'react';
import routes from '../constants/routes';
import NavigationDropdown from './NavigationDropdown';
import { AppBar } from './AppBar';
import { Button, Icon, IconButton } from 'native-base';
import { AntDesign } from '@expo/vector-icons';
import theme from "../constants/theme";

type NavigationProps = StackHeaderProps;

const NavigationBar: React.FC<NavigationProps> = ({
    navigation,
    back,
    options,
}) => {
    const [visible, setVisible] = React.useState(false);
    const openMenu = () => setVisible(true);
    const closeMenu = () => setVisible(false);

    return (
        <AppBar navigation={navigation} back={back} options={options}>
            <IconButton
                _icon={{ as: AntDesign, name: 'user', color: theme.colors.white }}
                size={'sm'}
                onPress={() => navigation.navigate(routes.addRecord.name)}
                accessibilityLabel={routes.addRecord.title}
            />
            <Button
                endIcon={
                    <Icon
                        as={AntDesign}
                        name={'form'}
                        size={'sm'}
                        color={'white'}
                    />
                }
                onPress={() => navigation.navigate(routes.forms.name)}
            >
                {options.title}
            </Button>
            <NavigationDropdown
                visible={visible}
                closeMenu={closeMenu}
                openMenu={openMenu}
                navigation={navigation}
            />
        </AppBar>
    );
};

export default NavigationBar;
