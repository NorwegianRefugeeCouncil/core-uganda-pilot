import React from 'react';
import { storiesOf } from '@storybook/react-native';
import { VStack, FormControl, Checkbox } from 'native-base';

import CenterView from '../CenterView';

storiesOf('Checkbox', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview', () => {
    return (
      <VStack space={4}>
        <FormControl>
          <Checkbox value="value" defaultIsChecked size="sm">
            Checkbox Label
          </Checkbox>
          <FormControl.HelperText>This is a helper text</FormControl.HelperText>
        </FormControl>
        <FormControl>
          <Checkbox value="value">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
            ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat. Lorem ipsum dolor sit amet,
            consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
            labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud
            exercitation ullamco laboris nisi ut aliquip ex ea commodo
            consequat.
          </Checkbox>
        </FormControl>
        <FormControl isDisabled>
          <Checkbox value="value">Checkbox Label</Checkbox>
        </FormControl>
      </VStack>
    );
  });
