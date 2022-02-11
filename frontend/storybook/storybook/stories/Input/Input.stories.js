import React from 'react';
import { storiesOf } from '@storybook/react-native';
import { text, boolean } from '@storybook/addon-knobs';
import { VStack, Input as InputNB, FormControl } from 'native-base';

import CenterView from '../CenterView';

storiesOf('Input', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview', () => {
    return (
      <VStack space={4}>
        <FormControl>
          <FormControl.Label>Label</FormControl.Label>
          <InputNB
            placeholder="This is a default input"
            value={text('input text', 'Default value')}
          />
          <FormControl.HelperText>This is a helper text</FormControl.HelperText>
        </FormControl>
        <FormControl isDisabled>
          <FormControl.Label>Disabled Input</FormControl.Label>
          <InputNB
            placeholder="This is a disabled input"
            value={text('input text', 'Disabled value')}
          />
        </FormControl>
        <FormControl isInvalid>
          <FormControl.Label>Invalid Input</FormControl.Label>
          <InputNB
            placeholder="This is an invalid input"
            value={text('input text', 'Invalid value')}
          />
          <FormControl.ErrorMessage>
            This is an error message
          </FormControl.ErrorMessage>
        </FormControl>
      </VStack>
    );
  });
