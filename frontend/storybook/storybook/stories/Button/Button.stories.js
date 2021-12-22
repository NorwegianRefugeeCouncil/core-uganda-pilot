import React from 'react';
import { action } from '@storybook/addon-actions';
import { storiesOf } from '@storybook/react-native';
import { boolean, select, text } from '@storybook/addon-knobs';
import { Button as ButtonNB, HStack } from 'native-base';
import { Icon } from 'core-design-system';
import { IconName } from 'core-design-system/lib/cjs/types/icons';

import CenterView from '../CenterView';

const IconNameList = Object.entries(IconName).map(([_, value]) => value);

storiesOf('Button', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview', () => {
    return (
      <HStack space={2}>
        <ButtonNB
          onPress={action('clicked-text')}
          colorScheme="primary"
          variant="major"
        >
          Primary Major
        </ButtonNB>
        <ButtonNB
          onPress={action('clicked-text')}
          colorScheme="secondary"
          variant="major"
        >
          Secondary Major
        </ButtonNB>
        <ButtonNB
          onPress={action('clicked-text')}
          colorScheme="primary"
          variant="minor"
        >
          Primary Minor
        </ButtonNB>
        <ButtonNB
          onPress={action('clicked-text')}
          colorScheme="secondary"
          variant="minor"
        >
          Secondary Minor
        </ButtonNB>
        <ButtonNB
          onPress={action('clicked-text')}
          colorScheme="primary"
          isDisabled
          variant="major"
        >
          Disabled
        </ButtonNB>
        <ButtonNB
          onPress={action('clicked-text')}
          colorScheme="primary"
          variant="major"
          startIcon={
            <Icon name={select('icon name', IconNameList, 'sdfdsf')} />
          }
        >
          With Icon
        </ButtonNB>
      </HStack>
    );
  })
  .add('Basic Button', () => {
    return (
      <ButtonNB
        onPress={action('clicked-text')}
        colorScheme={select(
          'Color scheme',
          ['primary', 'secondary'],
          'primary',
        )}
        isDisabled={boolean('disabled', false)}
        variant={select('Variant', ['major', 'minor'], 'major')}
      >
        {text('Button text', 'Submit')}
      </ButtonNB>
    );
  })
  .add('Button with Icon', () => {
    return (
      <ButtonNB
        onPress={action('clicked-text')}
        colorScheme={select(
          'Color scheme',
          ['primary', 'secondary'],
          'primary',
        )}
        isDisabled={boolean('disabled', false)}
        variant={select('Variant', ['major', 'minor'], 'major')}
        startIcon={
          <Icon name={select('icon name', IconNameList, IconName.ATTACHMENT)} />
        }
      >
        {text('Button text', 'Submit')}
      </ButtonNB>
    );
  });
