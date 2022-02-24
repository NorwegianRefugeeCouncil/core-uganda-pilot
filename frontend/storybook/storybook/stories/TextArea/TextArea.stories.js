import React from 'react';
import { storiesOf } from '@storybook/react-native';
import { VStack, TextArea, FormControl } from 'native-base';

import CenterView from '../CenterView';

storiesOf('TextArea', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview', () => {
    return (
      <VStack space={4}>
        <FormControl>
          <FormControl.Label size="xs">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
            ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat.
          </FormControl.Label>
          <TextArea
            placeholder="This is a default text area"
            value="Default value"
          />
          <FormControl.HelperText>This is a helper text</FormControl.HelperText>
        </FormControl>
        <FormControl isDisabled>
          <FormControl.Label size="xs">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
            ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat.
          </FormControl.Label>
          <TextArea
            placeholder="This is a disabled text area"
            value="Disabled value"
          />
        </FormControl>
        <FormControl isInvalid>
          <FormControl.Label size="xs">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
            ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat.
          </FormControl.Label>
          <TextArea
            placeholder="This is an invalid text area"
            value="Invalid value"
          />
          <FormControl.ErrorMessage>
            This is an error message
          </FormControl.ErrorMessage>
        </FormControl>
      </VStack>
    );
  });
