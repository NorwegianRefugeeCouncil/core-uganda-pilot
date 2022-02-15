import React from 'react';
import { action } from '@storybook/addon-actions';
import { storiesOf } from '@storybook/react-native';
import { boolean, select, text } from '@storybook/addon-knobs';
import { Button as ButtonNB, HStack, VStack } from 'native-base';
import { Icon, icons } from 'core-design-system';

import CenterView from '../CenterView';

const IconNameList = Object.keys(icons);

storiesOf('Button', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview', () => {
    return (
      <VStack space={4}>
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
        </HStack>
        <HStack space={2}>
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
        </HStack>
        <HStack space={2}>
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
              <Icon
                size="5"
                name={select('icon name', IconNameList, IconNameList[0])}
                viewBox="10 10 20 20" // to be removed when the files are replaced
              />
            }
          >
            With Icon
          </ButtonNB>
        </HStack>
      </VStack>
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
          <Icon name={select('icon name', IconNameList, IconNameList[0])} />
        }
      >
        {text('Button text', 'Submit')}
      </ButtonNB>
    );
  });
