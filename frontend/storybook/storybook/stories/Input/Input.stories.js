import React from 'react';
import { storiesOf } from '@storybook/react-native';
import { VStack, Input, FormControl } from 'native-base';
import { Icon } from 'core-design-system';

import CenterView from '../CenterView';

storiesOf('Input', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview', () => {
    return (
      <VStack space={4}>
        <FormControl>
          <FormControl.Label>Label</FormControl.Label>
          <Input
            placeholder="This is a default input"
            value="Default value"
          />
          <FormControl.HelperText>This is a helper text</FormControl.HelperText>
        </FormControl>
        <FormControl isDisabled>
          <FormControl.Label>Disabled Input</FormControl.Label>
          <Input
            placeholder="This is a disabled input"
            value="Disabled value"
          />
        </FormControl>
        <FormControl>
          <FormControl.Label>Label</FormControl.Label>
          <Input
            value="Valid value"
            InputRightElement={
              <Icon name="success" size="6" color="signalSuccess" mr={3} />
            }
          />
          <FormControl.HelperText>This is a valid input</FormControl.HelperText>
        </FormControl>
        <FormControl isInvalid>
          <FormControl.Label>Invalid Input</FormControl.Label>
          <Input
            placeholder="This is an invalid input"
            value="Invalid value"
            InputRightElement={
              <Icon name="error" size="6" color="signalDanger" mr={3} />
            }
          />
          <FormControl.ErrorMessage>
            This is an error message
          </FormControl.ErrorMessage>
        </FormControl>
      </VStack>
    );
  });
