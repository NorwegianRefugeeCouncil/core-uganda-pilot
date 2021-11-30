import { action } from '@storybook/addon-actions';
import { text } from '@storybook/addon-knobs';
import { storiesOf } from '@storybook/react-native';
import React from 'react';
import { Text } from 'react-native';
import CenterView from '../CenterView';
import {Button} from 'core-design-system'

storiesOf('Button', module)
  .addDecorator((getStory) =>
      <CenterView>{getStory()}</CenterView>
  )
  .add('with text', () => {
      console.log('Button story', typeof Button, Button)
     return  (
          <Button onPress={action('clicked-text')}>
              <Text>{text('Button text', 'Hello Button')}</Text>
          </Button>
      )
  })
  .add('with some emoji', () => (
    <Button onPress={action('clicked-emoji')}>
      <Text>😀 😎 👍 💯</Text>
    </Button>
  ));
