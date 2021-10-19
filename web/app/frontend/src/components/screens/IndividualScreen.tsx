import React from 'react';
import {FAB, Text} from 'react-native-paper';
import {common, layout} from '../../styles';
import {testIndividuals} from './IndividualsListScreen';
import {darkTheme} from '../../constants/theme';
import {View} from 'react-native';
import {useForm} from "react-hook-form";
import {Individual} from "../../../../client/src/types/models";
import {createIndividual, getIndividual} from "../../services/individuals";

const IndividualScreen: React.FC<any> = ({route}) => {
    const {id} = route.params;

    const person = getIndividual(id);

    const {register, handleSubmit, watch, formState: {errors}} = useForm({
        defaultValues: person
    });

    const onSubmit = (data: Individual) => {
        console.log(data);
        createIndividual(data);
    };

    return (
        <View style={[layout.container, layout.body, common.darkBackground]}>
            <Text theme={darkTheme}>
                {id == null ? 'new person' : testIndividuals[id].name}
            </Text>

            <form onSubmit={handleSubmit(onSubmit)}>
                {/* register your input into the hook by invoking the "register" function */}
                <input defaultValue="test" {...register("example")} />

                {/* include validation with required or other standard HTML validation rules */}
                <input {...register("exampleRequired", {required: true})} />
                {/* errors will return when field validation fails  */}
                {errors.exampleRequired && <span>This field is required</span>}

                <input type="submit"/>
            </form>
            <FAB
                style={layout.fab}
                icon="chevron-right"
                color={darkTheme.colors.white}
                onPress={() => console.log('Pressed')}
            />
        </View>
    );
};

export default IndividualScreen;
