import { storiesOf } from '@storybook/react-native';
import React from 'react';
import {Icon, theme, tokens} from 'core-design-system';
import { IconName } from 'core-design-system/lib/esm/types/icons';
import { select } from '@storybook/addon-knobs';

import CenterView from '../CenterView';

storiesOf('Icon', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Icon', () => {
    const IconNameList = Object.entries(IconName).map(([_, value]) => value);

    return (
      <Icon
        name={select('name', IconNameList, IconName.ATTACHMENT)}
        color={select('color', tokens.colors.icons, theme.colors.icons.dark)}
      />
    );
  });
