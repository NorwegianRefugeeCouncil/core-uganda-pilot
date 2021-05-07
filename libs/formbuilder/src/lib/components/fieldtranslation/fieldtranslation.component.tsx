import React, { Fragment, FunctionComponent, useEffect, useState } from 'react';
import {
    Button,
    FormLabel,
    FormInput
} from '@nrc.no/ui-toolkit';

type FieldTranslationProps = {
    name: string;
    translation: Translation
}

export interface Translation {
    [lang: string]: string
}

export const FieldTranslation: FunctionComponent<FieldTranslationProps> = (props) => {
    const [currentLang, setCurrentLang] = useState("")
    const [translation, setTranslation] = useState<Translation>(props.translation)

    const addTranslation = (language: string) => {
        const temp = translation
        temp[language] = ""
        setTranslation(temp)
    }

    const updateTranslation = (language: string, value: string) => {
        const temp = translation
        temp[language] = value
        setTranslation(temp)
    }

    const addNRCDefaultLanguages = () => {
        const temp = translation
        if (!("en" in temp)) {temp["en"] = ""}
        if (!("fr" in temp)) {temp["fr"] = ""}
        if (!("es" in temp)) {temp["es"] = ""}
        if (!("ar" in temp)) {temp["ar"] = ""}
        setTranslation(temp)
    }

    return (
        <Fragment>
            <FormInput name={`${props.name}-lang`} value={currentLang} placeholder="Language Code"
                onChange={(event) => {
                    setCurrentLang(event.target["value"])
                }}
            />
            <br />
            <Button kind="primary"
                onClick={(event) => {
                    if (currentLang.length > 0){
                        addTranslation(currentLang)
                        setCurrentLang("")
                    }
                }}
            >
                Add Translation
            </Button>
            <Button kind="primary"
                style={{
                    marginLeft: 2 + "vw"
                }}
                onClick={(event) => {
                    addNRCDefaultLanguages()
                }}
            >
                Add NRC Defaults
            </Button>
            <br />
            {
                Object.keys(translation).map((language) => {
                    return (
                        <Fragment>
                            <FormLabel>{language}:</FormLabel>
                            <FormInput name={`${props.name}-${language}`} defaultValue={translation[language]} onChange={(event) => {
                                updateTranslation(event.target["name"], event.target["value"])
                            }} />
                            <br />
                        </Fragment>
                    )
                })
            }
        </Fragment>
    )
}
