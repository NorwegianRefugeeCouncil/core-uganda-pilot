// eslint-disable-next-line no-unused-vars,@typescript-eslint/no-unused-vars
import React from 'react';
import { action } from '@storybook/addon-actions';
import { storiesOf } from '@storybook/react-native';
import { boolean, select, text } from '@storybook/addon-knobs';
import { Button as ButtonNB } from 'native-base';

import CenterView from '../CenterView';

storiesOf('Button', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Native Base Pure', () => {
    return (
      <ButtonNB
        onPress={action('clicked-text')}
        colorScheme={select('Color scheme', ['primary', 'secondary'], 'primary')}
        isDisabled={boolean('disabled', false)}
        variant={select('Variant', ['major', 'minor'], 'major')}
      >
        {text('Button text', 'Submit')}
      </ButtonNB>
    );
  });
