import React from 'react';
import { StyleProp } from 'react-native';
import { Appbar } from 'react-native-paper';
import pngSearch from '../../assets/png/search_white.png';
import pngFilter from '../../assets/png/symbol_filter.png';
import pngIndividuals from '../../assets/png/symbol_individuals.png';
import pngMore from '../../assets/png/symbol_more.png';
import common from '../styles/common';

type NavigationProps = {
  style: StyleProp<any>;
};

const Navigation: React.FC<NavigationProps> = ({ style }) => {
  return (
    <Appbar.Header>
      <Appbar.Action
        icon={pngIndividuals}
        accessibilityLabel={'Individuals'}
      />
      <Appbar.Action
        icon={pngSearch}
      />
      <Appbar.Content
        title="Home"
        titleStyle={common.textCentered}
      />
      <Appbar.Action
        icon={pngFilter}
      />
      <Appbar.Action
        icon={pngMore}
      />
    </Appbar.Header>
  );
};

export default Navigation;