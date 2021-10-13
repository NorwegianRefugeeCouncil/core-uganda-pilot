import React from 'react';
import { Appbar } from 'react-native-paper';
import pngSearch from '../../assets/png/search_white.png';
import pngFilter from '../../assets/png/symbol_filter.png';
import pngIndividuals from '../../assets/png/symbol_individuals.png';
import { common } from '../styles';
import { StackHeaderProps } from '@react-navigation/stack';
import routes from '../constants/routes';
import NavigationDropdown from './NavigationDropdown';

type NavigationProps = StackHeaderProps;

const NavigationBar: React.FC<NavigationProps> = (
  {
    navigation,
    back,
    options,
  }) => {

  const [visible, setVisible] = React.useState(false);
  const openMenu = () => setVisible(true);
  const closeMenu = () => setVisible(false);

  return (
    <Appbar.Header>
      {back ? <Appbar.BackAction onPress={navigation.goBack} /> : null}
      <Appbar.Action
        icon={pngIndividuals}
        accessibilityLabel={routes.individuals.title}
        onPress={() => navigation.navigate(routes.individuals.name)}
      />
      <Appbar.Action
        icon={pngSearch}
      />
      <Appbar.Content
        title={options.title}
        titleStyle={common.textCentered}
        onPress={() => navigation.navigate(routes.home.name)}
      />
      <Appbar.Action
        icon={pngFilter}
      />
      <NavigationDropdown
        visible={visible}
        closeMenu={closeMenu}
        openMenu={openMenu}
        navigation={navigation}
      />
    </Appbar.Header>
  );
};

export default NavigationBar;