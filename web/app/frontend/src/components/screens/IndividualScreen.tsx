import React from 'react';
import {common, layout} from '../../styles';
import {Button, ScrollView, TextInput, View} from 'react-native';
import {Controller, useForm} from "react-hook-form";
import {Individual} from "../../../../client/src/types/models";
import {Subject} from 'rxjs';
import iamClient from "../../utils/clients";
import {PartyAttributeDefinition, PartyAttributeDefinitionList} from "core-js-api-client/lib/types/models";
import _ from 'lodash';
import {createIndividual} from "../../services/individuals";
import FormControl from "../form/FormControl";

interface FlatIndividual extends Omit<Individual, 'attributes'> {
    attributes: { [p: string]: string }
}

const IndividualScreen: React.FC<any> = ({route}) => {
    const {id} = route.params;
    const [attributes, setAttributes] = React.useState<PartyAttributeDefinition[]>([]);
    const [individual, setIndividual] = React.useState<FlatIndividual>();
    let attributesSubject = new Subject([]);
    let individualsSubject = new Subject([]);
    let flatAttributes: { [p: string]: string } = {};

    const [isLoading, setIsLoading] = React.useState(true);

    const {control, handleSubmit, watch, formState: {errors}, register, getValues} = useForm();

    const onSubmit = (data: any) => {
        console.log('SUBMITTING', data);
        createIndividual(data);
    };

    React.useEffect(() => {
        attributesSubject.pipe(iamClient.PartyAttributeDefinitions().List()).subscribe(
            (data: PartyAttributeDefinitionList) => {
                setAttributes(data.items)
            }
        );
        individualsSubject.pipe(iamClient.Individuals().Get()).subscribe(
            (data: Individual) => {
                _(data.attributes).forEach((value, key) => {
                    flatAttributes[key] = value[0];
                });
                const flatIndividual = {...data, attributes: flatAttributes};
                setIndividual(flatIndividual)
            }
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
    React.useEffect(() => {
        if (individual) {
            setIsLoading(false);
        }
    }, [individual])

    return (
        <ScrollView>
            {!isLoading &&
            <View style={[layout.container, layout.body, common.darkBackground]}>


                {/*<form onSubmit={handleSubmit(onSubmit)}>*/}
                {attributes.map((a) =>
                    // <Controller
                    //     key={a.id}
                    //     name={a.id as '`${string}` | `${string}.${string}` | `${string}.${number}`'}
                    //     control={control}
                    //     rules={a.formControl.validation}
                    //     defaultValue={individual?.attributes[a.id] ? individual?.attributes[a.id] : undefined}
                    //     render={({field: {onChange, onBlur, value, ref}}) => {
                    //         return (
                    //             <TextInput
                    //                 onChangeText={onChange}
                    //                 value={value}
                    //                 onBlur={onBlur}
                    //                 ref={ref}
                    //             />
                    //         )
                    //     }}
                    // />
                    <FormControl
                        key={a.id}
                        formControl={a.formControl}
                        style={{width: '100%'}}
                        value={individual?.attributes[a.id]}
                        control={control}
                        name={a.id}
                    />

                )}
                {/*<Button title="Submit" onPress={handleSubmit(onSubmit)}/>*/}
                {/*{attributes.map((a) =>*/}
                {/*    <FormControl*/}
                {/*        key={a.id}*/}
                {/*        formControl={a.formControl}*/}
                {/*        style={{width: '100%'}}*/}
                {/*        value={individual?.attributes[a.id]}*/}
                {/*        control={control}*/}
                {/*        name={a.id}*/}
                {/*    />*/}
                {/*)}*/}
                <Button title="Submit" onPress={(data) => {
                    // console.log('ISFDPXVJOKsdlp;', getValues())
                    onSubmit(getValues());
                }}/>
                {/*</form>*/}
            </View>
            }
        </ScrollView>
    );
};

export default IndividualScreen;
