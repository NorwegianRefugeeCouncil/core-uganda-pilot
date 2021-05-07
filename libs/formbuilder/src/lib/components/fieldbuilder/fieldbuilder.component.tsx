import React, { useState, Fragment, FunctionComponent } from 'react';
import { FieldTypePicker, FieldType } from '../fieldtype/fieldtype.component';
import { FieldInfo } from '../fieldinfo/fieldinfo.component';
import { FieldConfig, GenericFieldConfig } from '../fieldconfig/fieldconfig.component';

type FieldBuilderProps = {
}

enum Tab {
    info = "Information",
    type = "Type",
    config = "Config"
}

const buildTabLink = (tab: Tab, currentTab: Tab, setCurrentTab: (React.Dispatch<React.SetStateAction<Tab>>), disabled=false) => {
    if (tab == currentTab){
        return disabled ? 
            <a className="nav-link disabled" aria-disabled="true" aria-current="page" onClick={(event) => {setCurrentTab(tab)}}>{tab}</a> :
            <a className="nav-link active" aria-current="page" onClick={(event) => {setCurrentTab(tab)}}>{tab}</a>
    } else {
        return disabled ? 
            <a className="nav-link disabled" aria-disabled="true" onClick={(event) => {setCurrentTab(tab)}}>{tab}</a> :
            <a className="nav-link" onClick={(event) => {setCurrentTab(tab)}}>{tab}</a>
    }
}

export const FieldBuilder: FunctionComponent<FieldBuilderProps> = (props) => {
    const [currentTab, setCurrentTab] = useState<Tab>(Tab.info)
    const [fieldType, setFieldType] = useState<FieldType | undefined>()
    const [name, setName] = useState("")
    const [description, setDescription] = useState("")
    const [fieldConfig, setFieldConfig] = useState<GenericFieldConfig>({} as GenericFieldConfig)

    const patchFieldConfig = (key: string, value: any) => {
        const tempFieldConfig = fieldConfig
        tempFieldConfig[key] = value
        setFieldConfig(tempFieldConfig)
    }
    
    return (
        <Fragment>
            <div className={'container'}>
                <div className={'row'}>
                    <div className={'col-12 mb-5'}>
                        <ul className="nav nav-tabs">
                            <li className="nav-item">
                                {
                                    buildTabLink(Tab.info, currentTab, setCurrentTab)
                                }
                            </li>
                            <li className="nav-item">
                                {
                                    buildTabLink(Tab.type, currentTab, setCurrentTab)
                                }
                            </li>
                            {
                                fieldType === undefined ? 
                                <li className="nav-item">
                                    {
                                        buildTabLink(Tab.config, currentTab, setCurrentTab, true)
                                    }
                                </li> :
                                    <li className="nav-item">
                                    {
                                        buildTabLink(Tab.config, currentTab, setCurrentTab)
                                    }
                                </li>
                            }
                            
                        </ul>
                        <br/>
                        {
                            currentTab == Tab.info &&
                            <form onChange={
                                (event) => {
                                    console.log(event.target["name"])
                                    console.log(event.target["value"])
                                    switch (event.target["name"]) {
                                        case "name":
                                            setName(event.target["value"])
                                            break;

                                        case "description":
                                            setDescription(event.target["value"])
                                            break;
                                    
                                        default:
                                            break;
                                    }
                                }
                            }>
                                <FieldInfo name={name} description={description} />
                            </form>
                        }
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
                            <form onChange={
                                (event) => {
                                    patchFieldConfig(event.target["name"], event.target["value"])
                                }
                            }>
                                <FieldConfig fieldType={fieldType} fieldProps={fieldConfig} />
                            </form>
                        }
                    </div>
                </div>
            </div>
        </Fragment>
    )
}