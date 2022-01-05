import React from 'react';
import { storiesOf } from '@storybook/react-native';
import { VStack, Text } from 'native-base';
import { select } from '@storybook/addon-knobs';

import CenterView from '../CenterView';

storiesOf('Text', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview, Variants', () => (
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
  ))
  .add('Testing Theme', () => (
    <VStack space={2}>
      <Text variant="bodyText">Note: Not all combinations are intended</Text>
      <Text
        fontSize={select('size', ['xs', 'sm', 'md', 'lg', 'xl'], 'xl')}
        fontFamily={select(
          'family',
          [
            'display',
            'heading',
            'title',
            'bodyText',
            'inline',
            'date',
            'label',
            'caption',
          ],
          'display',
        )}
        fontWeight={select('weight', ['400', '700'], '700')}
        fontStyle={select('style', ['italic', 'normal'], 'normal')}
        textDecorationLine={select(
          'decoration',
          ['underline', 'none'],
          'underline',
        )}
      >
        Example Text
      </Text>
    </VStack>
  ));
