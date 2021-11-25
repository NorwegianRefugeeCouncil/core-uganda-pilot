import React from 'react';
import {IconName, IconProps} from "../types/icons";
import {extendTheme, Icon as IconNB, NativeBaseProvider} from "native-base";
import ActivityCamp from './icons/ActivityCamp';
import ActivityEducation from './icons/ActivityEducation';
import ActivityFood from './icons/ActivityFood';
import ActivityIcla from './icons/ActivityIcla';
import ActivityShelter from './icons/ActivityShelter';
import ActivityWash from './icons/ActivityWash';
import Attachment from './icons/Attachment';
import Beneficiary from "./icons/Beneficiary";
import Calender from "./icons/Calender";
import Call from "./icons/Call";
import Case from "./icons/Case";
import Delete from "./icons/Delete";
import Edit from "./icons/Edit";
import Email from "./icons/Email";
import FaceError from "./icons/FaceError";
import FaceGeneral from "./icons/FaceGeneral";
import FaceSuccess from "./icons/FaceSuccess";
import Female from "./icons/Female";
import Filter from "./icons/Filter";
import Home from "./icons/Home";
import Image from "./icons/Image";
import Lock from "./icons/Lock";
import Male from "./icons/Male";
import Menu from "./icons/Menu";
import More from "./icons/More";
import Next from "./icons/Next";
import Notification from "./icons/Notification";
import Plus from "./icons/Plus";
import Print from "./icons/Print";
import Protection from "./icons/Protection";
import Registration from "./icons/Registration";
import Reload from "./icons/Reload";
import Save from "./icons/Save";
import Search from "./icons/Search";
import Setting from "./icons/Setting";
import Share from "./icons/Share";
import Sync from "./icons/Sync";
import Thumbnail from "./icons/Thumbnail";
import Unlock from "./icons/Unlock";
import Upload from "./icons/Upload";
import UserAdd from "./icons/UserAdd";
import UserDelete from "./icons/UserDelete";
import UserGroup from "./icons/UserGroup";
import UserGroupAdd from "./icons/UserGroupAdd";
import Whatsapp from "./icons/Whatsapp";
import {Path} from "react-native-svg";

const Icon = ({name, variant}: IconProps) => {
    let iconComponent;
    switch (name) {
        case IconName.ATTACHMENT:
            iconComponent = Attachment;
            break;
        case IconName.BENEFICIARY:
            iconComponent = Beneficiary;
            break;
        case IconName.CALENDAR:
            iconComponent = Calender;
            break;
        case IconName.CALL:
            iconComponent = Call;
            break;
        case IconName.CASE:
            iconComponent = Case;
            break;
        case IconName.DELETE:
            iconComponent = Delete;
            break;
        case IconName.EDIT:
            iconComponent = Edit;
            break;
        case IconName.EMAIL:
            iconComponent = Email;
            break;
        case IconName.FACE_ERROR:
            iconComponent = FaceError;
            break;
        case IconName.FACE_GENERAL:
            iconComponent = FaceGeneral;
            break;
        case IconName.FACE_SUCCESS:
            iconComponent = FaceSuccess;
            break;
        case IconName.FEMALE:
            iconComponent = Female;
            break;
        case IconName.FILTER:
            iconComponent = Filter;
            break;
        case IconName.HOME:
            iconComponent = Home;
            break;
        case IconName.IMAGE:
            iconComponent = Image;
            break;
        case IconName.LOCK:
            iconComponent = Lock;
            break;
        case IconName.MALE:
            iconComponent = Male;
            break;
        case IconName.MENU:
            iconComponent = Menu;
            break;
        case IconName.MORE:
            iconComponent = More;
            break;
        case IconName.NEXT:
            iconComponent = Next;
            break;
        case IconName.NOTIFICATION:
            iconComponent = Notification;
            break;
        case IconName.PLUS:
            iconComponent = Plus;
            break;
        case IconName.PRINT:
            iconComponent = Print;
            break;
        case IconName.PROTECTION:
            iconComponent = Protection;
            break;
        case IconName.REGISTRATION:
            iconComponent = Registration;
            break;
        case IconName.RELOAD:
            iconComponent = Reload;
            break;
        case IconName.SAVE:
            iconComponent = Save;
            break;
        case IconName.SEARCH:
            iconComponent = Search;
            break;
        case IconName.SETTING:
            iconComponent = Setting;
            break;
        case IconName.SHARE:
            iconComponent = Share;
            break;
        case IconName.SYNC:
            iconComponent = Sync;
            break;
        case IconName.THUMBNAIL:
            iconComponent = Thumbnail;
            break;
        case IconName.UNLOCK:
            iconComponent = Unlock;
            break;
        case IconName.UPLOAD:
            iconComponent = Upload;
            break;
        case IconName.USER_ADD:
            iconComponent = UserAdd;
            break;
        case IconName.USER_DELETE:
            iconComponent = UserDelete;
            break;
        case IconName.USER_GROUP:
            iconComponent = UserGroup;
            break;
        case IconName.USER_GROUP_ADD:
            iconComponent = UserGroupAdd;
            break;
        case IconName.WHATSAPP:
            iconComponent = Whatsapp;
            break;
        case IconName.ACTIVITY_CAMP:
            iconComponent = ActivityCamp;
            break;
        case IconName.ACTIVITY_EDUCATION:
            iconComponent = ActivityEducation;
            break;
        case IconName.ACTIVITY_FOOD:
            iconComponent = ActivityFood;
            break;
        case IconName.ACTIVITY_ICLA:
            iconComponent = ActivityIcla;
            break;
        case IconName.ACTIVITY_SHELTER:
            iconComponent = ActivityShelter;
            break;
        case IconName.ACTIVITY_WASH:
            iconComponent = ActivityWash;
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
