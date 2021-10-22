import React from 'react';
import {common, layout} from '../../styles';
import {Button, ScrollView, View} from 'react-native';
import {useForm} from "react-hook-form";
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
    let submitData: { [p: string]: string[] } = {};

    const [isLoading, setIsLoading] = React.useState(true);

    const {control, handleSubmit, formState} = useForm();

    const onSubmit = (data: string[]) => {
        console.log('SUBMITTING', data)

        _(data).forEach((value, key) => {
            submitData[key] = [value];
        });

        createIndividual({
            id,
            attributes: submitData,
            partyTypeIds: individual?.partyTypeIds || []
        });
    };
    // console.log('ERRORS', errors)

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
            {!isLoading && (
                <View style={[layout.container, layout.body, common.darkBackground]}>
                    {attributes.map((a) =>
                        <FormControl
                            key={a.id}
                            formControl={a.formControl}
                            style={{width: '100%'}}
                            value={individual?.attributes[a.id]}
                            control={control}
                            name={a.id}
                            errors={formState.errors}
                        />
                    )}
                    <Button
                        title="Submit"
                        onPress={handleSubmit(onSubmit)}
                    />
                </View>
            )}
        </ScrollView>
    );
};

export default IndividualScreen;
