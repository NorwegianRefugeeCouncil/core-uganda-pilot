import React from 'react';
import { storiesOf } from '@storybook/react-native';
import { VStack, Text } from 'native-base';

import CenterView from '../CenterView';

storiesOf('Text', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview', () => (
    <VStack space={2}>
      <Text variant="display">Display</Text>
      <Text variant="heading">Heading</Text>
      <Text variant="title">Title</Text>
      <Text variant="bodyText">Body Text</Text>
      <Text variant="caption">Caption</Text>
      <Text variant="inline">Inline</Text>
      <Text variant="date">Date</Text>
      <Text variant="label">Label</Text>
    </VStack>
  ));
