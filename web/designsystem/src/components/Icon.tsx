import React from 'react';
import {IconName, IconProps} from "../types/icons";
import {Icon as IconNB, NativeBaseProvider} from "native-base";
import * as icons from './icons';
import {Path} from "react-native-svg";

const Icon = ({name, variant}: IconProps) => {
    let iconComponent;
    switch (name) {
        case IconName.ATTACHMENT:
            iconComponent = icons.Attachment;
            break;
        case IconName.BENEFICIARY:
            iconComponent = icons.Beneficiary;
            break;
        case IconName.CALENDAR:
            iconComponent = icons.Calender;
            break;
        case IconName.CALL:
            iconComponent = icons.Call;
            break;
        case IconName.CASE:
            iconComponent = icons.Case;
            break;
        case IconName.DELETE:
            iconComponent = icons.Delete;
            break;
        case IconName.EDIT:
            iconComponent = icons.Edit;
            break;
        case IconName.EMAIL:
            iconComponent = icons.Email;
            break;
        case IconName.FACE_ERROR:
            iconComponent = icons.FaceError;
            break;
        case IconName.FACE_GENERAL:
            iconComponent = icons.FaceGeneral;
            break;
        case IconName.FACE_SUCCESS:
            iconComponent = icons.FaceSuccess;
            break;
        case IconName.FEMALE:
            iconComponent = icons.Female;
            break;
        case IconName.FILTER:
            iconComponent = icons.Filter;
            break;
        case IconName.HOME:
            iconComponent = icons.Home;
            break;
        case IconName.IMAGE:
            iconComponent = icons.Image;
            break;
        case IconName.LOCK:
            iconComponent = icons.Lock;
            break;
        case IconName.MALE:
            iconComponent = icons.Male;
            break;
        case IconName.MENU:
            iconComponent = icons.Menu;
            break;
        case IconName.MORE:
            iconComponent = icons.More;
            break;
        case IconName.NEXT:
            iconComponent = icons.Next;
            break;
        case IconName.NOTIFICATION:
            iconComponent = icons.Notification;
            break;
        case IconName.PLUS:
            iconComponent = icons.Plus;
            break;
        case IconName.PRINT:
            iconComponent = icons.Print;
            break;
        case IconName.PROTECTION:
            iconComponent = icons.Protection;
            break;
        case IconName.REGISTRATION:
            iconComponent = icons.Registration;
            break;
        case IconName.RELOAD:
            iconComponent = icons.Reload;
            break;
        case IconName.SAVE:
            iconComponent = icons.Save;
            break;
        case IconName.SEARCH:
            iconComponent = icons.Search;
            break;
        case IconName.SETTING:
            iconComponent = icons.Setting;
            break;
        case IconName.SHARE:
            iconComponent = icons.Share;
            break;
        case IconName.SYNC:
            iconComponent = icons.Sync;
            break;
        case IconName.THUMBNAIL:
            iconComponent = icons.Thumbnail;
            break;
        case IconName.UNLOCK:
            iconComponent = icons.Unlock;
            break;
        case IconName.UPLOAD:
            iconComponent = icons.Upload;
            break;
        case IconName.USER_ADD:
            iconComponent = icons.UserAdd;
            break;
        case IconName.USER_DELETE:
            iconComponent = icons.UserDelete;
            break;
        case IconName.USER_GROUP:
            iconComponent = icons.UserGroup;
            break;
        case IconName.USER_GROUP_ADD:
            iconComponent = icons.UserGroupAdd;
            break;
        case IconName.WHATSAPP:
            iconComponent = icons.Whatsapp;
            break;
        case IconName.ACTIVITY_CAMP:
            iconComponent = icons.ActivityCamp;
            break;
        case IconName.ACTIVITY_EDUCATION:
            iconComponent = icons.ActivityEducation;
            break;
        case IconName.ACTIVITY_FOOD:
            iconComponent = icons.ActivityFood;
            break;
        case IconName.ACTIVITY_ICLA:
            iconComponent = icons.ActivityIcla;
            break;
        case IconName.ACTIVITY_SHELTER:
            iconComponent = icons.ActivityShelter;
            break;
        case IconName.ACTIVITY_WASH:
            iconComponent = icons.ActivityWash;
            break;
        default:
            iconComponent = () => {
                return <Path/>
            };
    }

    return (
        <NativeBaseProvider>
            <IconNB
                viewBox="0 0 40 40"
                style={{height: '40px', width: '40px', display: 'flex'}}
                children={iconComponent(variant)}
            />
        </NativeBaseProvider>
    );
}

export default Icon
