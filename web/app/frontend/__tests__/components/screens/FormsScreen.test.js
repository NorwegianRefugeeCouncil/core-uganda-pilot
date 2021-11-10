import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import { cleanup, render } from '@testing-library/react-native';
import React from 'react';
import { Provider as PaperProvider } from 'react-native-paper';

import FormsScreen from '../../../src/components/screens/FormsScreen';
import theme, { NavigationTheme } from '../../../src/constants/theme';
import testIds from '../../../src/testIds';

const SingleScreenWrapper = ({ screen }) => {
    const Stack = createStackNavigator();
    return (
        <PaperProvider theme={theme}>
            <NavigationContainer them={NavigationTheme}>
                <Stack.Navigator>
                    <Stack.Screen name={'test'} component={screen} />
                </Stack.Navigator>
            </NavigationContainer>
        </PaperProvider>
    );
};

describe(FormsScreen.name, () => {
    afterEach(cleanup);

    test('renders correctly', async () => {
        const tree = render(
            <SingleScreenWrapper screen={FormsScreen} />
        ).toJSON();
        expect(tree).toMatchSnapshot();
    });

    test('renders a list of forms', async () => {
        const { findAllByTestId } = render(
            <SingleScreenWrapper screen={FormsScreen} />
        );
        const formListItems = await findAllByTestId(testIds.formListItem);
        expect(formListItems.length).toBeGreaterThan(0);
    });

    test.todo('lets the user navigate to a form by pressing it');
});
