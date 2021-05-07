import React, { Fragment, FunctionComponent } from 'react';
import {
    FormLabel,
    FormInput
} from '@nrc.no/ui-toolkit';
import { FieldTranslation } from '../fieldtranslation/fieldtranslation.component';

type FieldInfoProps = {
    name: string;
    description: string; 
    tooltip: string;
}

const renderInfoField = (label: string, name: string, value: string, placeholder: string, collapsed = true) => {
    return (
        <Fragment>
            <div className="accordion" id={`accordion-${name}`}>
                <div className="accordion-item">
                    <h2 className="accordion-header" id={`heading-${name}`}>
                    <button className={`accordion-button ${collapsed ? "collapsed" : ""}`} type="button" data-bs-toggle="collapse" data-bs-target={`#collapse-${name}`} aria-controls={`collapse-${name}`}>
                        {label}
                    </button>
                    </h2>
                    <div id={`collapse-${name}`} className={`accordion-collapse collapse ${collapsed ? "" : "show"}`} aria-labelledby={`heading-${name}`} data-bs-parent={`accordion-${name}`}>
                        <div className="accordion-body">
                            <FormInput name={name} value={value} placeholder={placeholder} /><br/>
                            <div className="accordion" id={`accordion-inner-${name}`}>
                                <div className="accordion-item">
                                    <h2 className="accordion-header" id={`heading-inner-${name}`}>
                                    <button className="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target={`#collapse-inner-${name}`} aria-controls={`collapse-inner-${name}`}>
                                        Translations
                                    </button>
                                    </h2>
                                    <div id={`collapse-inner-${name}`} className="accordion-collapse collapse" aria-labelledby={`heading-inner-${name}`} data-bs-parent={`accordion-inner-${name}`}>
                                        <div className="accordion-body">
                                            <FieldTranslation name={name} translation={{}}/>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <br />
                        </div>
                    </div>
                </div>
            </div>
            <br/>
        </Fragment>
    )
}

export const FieldInfo: FunctionComponent<FieldInfoProps> = (props) => {
    return (
        <Fragment>
            {
                renderInfoField("Name", "name", props.name, "Default Name", false)
            }
            {
                renderInfoField("Description", "description", props.description, "Default Description")
            }
            {
                renderInfoField("Tooltip", "tooltip", props.tooltip, "Default Tooltip")
            }
        </Fragment>
    )
}