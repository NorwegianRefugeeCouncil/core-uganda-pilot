import React from 'react';
import { storiesOf } from '@storybook/react-native';
import { VStack, Text } from 'native-base';

import CenterView from '../CenterView';

storiesOf('Text', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview, Variants', () => (
    <VStack space={2}>
      <Text variant="display1">Display 1</Text>
      <Text variant="display2">Display 2</Text>
      <Text variant="heading1">Heading 1</Text>
      <Text variant="heading2">Heading 2</Text>
      <Text variant="heading3">Heading 3</Text>
      <Text variant="heading4">Heading 4</Text>
      <Text variant="heading5">Heading 5</Text>
      <Text variant="title1">Title 1</Text>
      <Text variant="title2">Title 2</Text>
      <Text variant="title3">Title 3</Text>
      <Text variant="body1">Body 1</Text>
      <Text variant="body2">Body 2</Text>
      <Text variant="caption">Caption</Text>
      <Text variant="inline">Inline</Text>
      <Text variant="date1">Date 1</Text>
      <Text variant="date2">Date 2</Text>
      <Text variant="label">Label</Text>
    </VStack>
  ));
