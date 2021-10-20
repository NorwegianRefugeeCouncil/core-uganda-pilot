import React from 'react';
import {common, layout} from '../../styles';
import {Button, ScrollView, TextInput, View} from 'react-native';
import {Controller, useForm} from "react-hook-form";
import {Individual} from "../../../../client/src/types/models";
import {Subject} from 'rxjs';
import iamClient from "../../utils/clients";
import {PartyAttributeDefinition, PartyAttributeDefinitionList} from "core-js-api-client/lib/types/models";
import FormControl from "../form/FormControl";

const IndividualScreen: React.FC<any> = ({route}) => {
    const {id} = route.params;
    const [attributes, setAttributes] = React.useState<PartyAttributeDefinition[]>([]);
    const [individual, setIndividual] = React.useState<Individual>();
    let attributesSubject = new Subject([]);
    let individualsSubject = new Subject([]);

    const {control, handleSubmit, watch, formState: {errors}, register, getValues} = useForm({
        defaultValues: individual?.attributes
    });

    const onSubmit = (data: any) => {
        console.log('SUBMITTING', data);
        // createIndividual(data);
    };

    React.useEffect(() => {
        attributesSubject.pipe(iamClient.PartyAttributeDefinitions().List()).subscribe(
            (data: PartyAttributeDefinitionList) => {
                setAttributes(data.items)
            }
        );
        individualsSubject.pipe(iamClient.Individuals().Get()).subscribe(
            (data: Individual) => setIndividual(data)
        );
        attributesSubject.next(null);
        individualsSubject.next(id);

        return () => {
            if (attributesSubject) {
                attributesSubject.unsubscribe();
            }
            if (individualsSubject) {
                individualsSubject.unsubscribe();
            }
        };
    }, []);

    return (
        <ScrollView>
            <View style={[layout.container, layout.body, common.darkBackground]}>
                {/*{attributes.map((a) =>*/}

                {/*    <Controller*/}
                {/*        key={a.id}*/}
                {/*        control={control}*/}
                {/*        rules={a.formControl.validation}*/}
                {/*        render={({field: {onChange, onBlur, value}}) => (*/}
                {/*            <TextInput*/}
                {/*                onBlur={onBlur}*/}
                {/*                onChangeText={onChange}*/}
                {/*                value={value}*/}
                {/*            />*/}
                {/*        )}*/}
                {/*        name={'8514da51-aad5-4fb4-a797-8bcc0c969b27'}*/}
                {/*        defaultValue={a.formControl.value || a.formControl.defaultValue}*/}
                {/*    />*/}
                {/*)}*/}

                {/*<Button title="Submit" onPress={handleSubmit(onSubmit)}/>*/}

                {/*<form onSubmit={(a)=>{*/}
                {/*    console.log('FORM', a)*/}
                {/*}}>*/}

                {attributes.map((a) =>
                    <FormControl
                        key={a.id}
                        formControl={a.formControl}
                        style={{width: '100%'}}
                        value={individual?.attributes[a.id]}
                        control={control}
                        // name={a.id}
                        {...register(a.id as '`${string}` | `${string}.${string}` | `${string}.${number}`')}
                    />
                )}
                    <Button title="Submit" onPress={handleSubmit(onSubmit)} />

                {/*</form>*/}
            </View>
        </ScrollView>
    );
};

export default IndividualScreen;
