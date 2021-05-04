import React, { useState, Fragment, FunctionComponent } from 'react';
import {
    Card,
    FormLabel,
    FormSelect
} from '@nrc.no/ui-toolkit';
import { FieldType, FieldTypes } from '../fieldtype/fieldtype.component';

type FieldBuilderProps = {
}

enum Tabs {
    type = "Field Type",
    config = "Field Config"
}

export const FieldBuilder: FunctionComponent<FieldBuilderProps> = (props) => {
    const [tab, setTab] = useState<Tabs>(Tabs.type)
    const [fieldType, setFieldType] = useState("Text")    
    return (
        <Fragment>
            <div className={'container'}>
                <div className={'row'}>
                    <div className={'col-12 mb-5'}>
                        <ul className="nav nav-tabs">
                            <li className="nav-item">
                                <a className="nav-link active" aria-current="page" onClick={(event) => {setTab(Tabs.type)}}>{Tabs.type}</a>
                            </li>
                            <li className="nav-item">
                                <a className="nav-link active" aria-current="page" onClick={(event) => {setTab(Tabs.config)}}>{Tabs.config}</a>
                            </li>
                        </ul>
                        <br/>
                        {   
                            tab == Tabs.type &&
                            <form onChange={
                                (event) => {
                                    setFieldType(event.target["value"])
                                }
                            }>
                                <FieldType />
                            </form>
                        }
                        {
                            tab == Tabs.config &&
                            "Current type " + fieldType 
                        }
                    </div>
                </div>
            </div>
        </Fragment>
    )
}