import { storiesOf } from '@storybook/react-native';
import { VStack, Text } from 'native-base';

import CenterView from '../CenterView';

storiesOf('Text', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview, Variants', () => (
    <VStack space={2}>
      <Text variant="display" fontSize="3xl">
        Display
      </Text>
      <Text variant="display" fontSize="2xl">
        Display
      </Text>
      <Text variant="heading" fontSize="xl">
        Heading
      </Text>
      <Text variant="heading" fontSize="lg">
        Heading
      </Text>
      <Text variant="heading" fontSize="md">
        Heading
      </Text>
      <Text variant="heading" fontSize="sm">
        Heading
      </Text>
      <Text variant="heading" fontSize="3xs">
        Heading
      </Text>
      <Text variant="title" fontSize="md">
        Title
      </Text>
      <Text variant="title" fontSize="sm">
        Title
      </Text>
      <Text variant="title" fontSize="2xs">
        Title
      </Text>
      <Text variant="body" fontSize="xs">
        Body
      </Text>
      <Text variant="body" fontSize="2xs">
        Body
      </Text>
      <Text variant="caption" fontSize="3xs">
        Caption
      </Text>
      <Text variant="inline" fontSize="xs">
        Inline
      </Text>
      <Text variant="date" fontSize="xs">
        Date
      </Text>
      <Text variant="date" fontSize="2xs">
        Date
      </Text>
      <Text variant="label" fontSize="xs">
        Label
      </Text>
    </VStack>
  ));
