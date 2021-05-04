import React, { useState, Fragment, FunctionComponent } from 'react';
import {
    Card,
    FormLabel,
    FormSelect
} from '@nrc.no/ui-toolkit';
import { FieldTypePicker, FieldType } from '../fieldtype/fieldtype.component';

type FieldBuilderProps = {
}

enum Tab {
    type = "Field Type",
    config = "Field Config"
}

const buildTab = (tab: Tab, currentTab: Tab, setCurrentTab: (React.Dispatch<React.SetStateAction<Tab>>)) => {
    if (tab == currentTab){
        return <a className="nav-link active" aria-current="page" onClick={(event) => {setCurrentTab(tab)}}>{tab}</a>
    } else {
        return <a className="nav-link" onClick={(event) => {setCurrentTab(tab)}}>{tab}</a>
    }
}

export const FieldBuilder: FunctionComponent<FieldBuilderProps> = (props) => {
    const [currentTab, setCurrentTab] = useState<Tab>(Tab.type)
    const [fieldType, setFieldType] = useState<FieldType | undefined>()    
    return (
        <Fragment>
            <div className={'container'}>
                <div className={'row'}>
                    <div className={'col-12 mb-5'}>
                        <ul className="nav nav-tabs">
                            <li className="nav-item">
                                {
                                    buildTab(Tab.type, currentTab, setCurrentTab)
                                }
                            </li>
                            <li className="nav-item">
                                {
                                    buildTab(Tab.config, currentTab, setCurrentTab)
                                }
                            </li>
                        </ul>
                        <br/>
                        {   
                            currentTab == Tab.type &&
                            <form onChange={
                                (event) => {
                                    setFieldType(event.target["value"])
                                }
                            }>
                                <FieldTypePicker value={fieldType} />
                            </form>
                        }
                        {
                            currentTab == Tab.config &&
                            "Current type " + fieldType 
                        }
                    </div>
                </div>
            </div>
        </Fragment>
    )
}