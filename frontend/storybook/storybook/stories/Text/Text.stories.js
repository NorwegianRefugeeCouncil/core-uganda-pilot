import { storiesOf } from '@storybook/react-native';
import { VStack, Text } from 'native-base';

import CenterView from '../CenterView';

storiesOf('Text', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview, Variants', () => (
    <VStack space={2}>
      <Text variant="display" level="1">
        Display 1
      </Text>
      <Text variant="display" level="2">
        Display 2
      </Text>
      <Text variant="heading" level="1">
        Heading 1
      </Text>
      <Text variant="heading" level="2">
        Heading 2
      </Text>
      <Text variant="heading" level="3">
        Heading 3
      </Text>
      <Text variant="heading" level="4">
        Heading 4
      </Text>
      <Text variant="heading" level="5">
        Heading 5
      </Text>
      <Text variant="body" level="1">
        Body 1
      </Text>
      <Text variant="body" level="2">
        Body 2
      </Text>
      <Text variant="caption">Caption</Text>
      <Text variant="inline">Inline</Text>
      <Text variant="label">Label</Text>
    </VStack>
  ));
