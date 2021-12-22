import { storiesOf } from '@storybook/react-native';
import React from 'react';
import { Icon, icons, theme, tokens } from 'core-design-system';
import { select } from '@storybook/addon-knobs';

import CenterView from '../CenterView';

storiesOf('Icon', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Icon', () => {
    const IconNameList = Object.keys(icons);

    return (
      <Icon
        name={select('name', IconNameList, IconNameList[0])}
        color={select('color', tokens.colors.icons, theme.colors.icons.dark)}
      />
    );
  });
