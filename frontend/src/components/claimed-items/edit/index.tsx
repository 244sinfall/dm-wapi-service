import React, {useCallback, useMemo, useState} from 'react';
import {ClaimedItemInterface} from "../../../model/claimed-items/types";
import ModalTitle from "../../common/modal-title";
import TextInput from "../../common/text-input";
import {PERMISSION} from "../../../model/user/types";
import ActionButton from "../../common/action-button";
import Field from "../../common/field";
import './styles.css'
import apiStringDateToString from "../../../utils/api-string-date-to-string";
import {UserStateUserInfo} from "../../../model/user/types";

export type ClaimedItemEditorProps = {
    onEdit: (editedItem: ClaimedItemInterface) => Promise<void>,
    onApprove: (id: string) => Promise<void>,
    onDelete: (id: string) => Promise<void>,
    onClose: () => void
    item: ClaimedItemInterface
    user: UserStateUserInfo
}
const ClaimedItemEditor = (props: ClaimedItemEditorProps) => {
    const [changeable, setChangeable] = useState<ClaimedItemInterface>(props.item)
    const handleChange = useCallback(<K extends keyof ClaimedItemInterface, V extends ClaimedItemInterface[K]>(key: K, value: V) => {
        setChangeable(prev => ({...prev, [key]: value}))
    }, [])
    const containerOptions = useMemo(() => ({collapsedOptions: {widthToCollapse: 480}}), [])
    return (
        <ModalTitle className="claimed-items-table-edit" title="Редактировать предмет" closeCallback={props.onClose}>

            <Field title="Качество" containerOptions={containerOptions}>
                <TextInput disabled={true} value={props.item.quality} />
            </Field>
            <Field title="Название" containerOptions={containerOptions}>
                <TextInput maxLength={256}
                           disabled={props.user.apiUser != null && props.user.apiUser.permission < PERMISSION.Admin}
                           onChange={newName => handleChange("name", newName)}
                           value={changeable.name}/>
            </Field>
            <Field title="Ссылка на предмет" containerOptions={containerOptions}>
                <TextInput maxLength={256}
                           disabled={props.user.apiUser != null && props.user.apiUser.permission < PERMISSION.Admin}
                           onChange={newLink => handleChange("link", newLink)}
                           value={changeable.link}/>
            </Field>
            <Field title="Владелец" containerOptions={containerOptions}>
                <TextInput maxLength={256}
                           value={changeable.owner}
                           onChange={newOwner => handleChange("owner", newOwner)}/>
            </Field>
            <Field title="Профиль владельца" containerOptions={containerOptions}>
                <TextInput maxLength={256}
                           value={changeable.ownerProfile}
                           onChange={newOwnerProfile => handleChange("ownerProfile", newOwnerProfile)}/>
            </Field>
            <Field title="Доказательство отыгрыша" containerOptions={containerOptions}>
                <TextInput maxLength={256}
                           value={changeable.ownerProofLink}
                           onChange={newProofLink => handleChange("ownerProofLink", newProofLink)}/>
            </Field>
            <Field title="Согласовавший рецензент" containerOptions={containerOptions}>
                <TextInput maxLength={256}
                           disabled={props.user.apiUser != null && props.user.apiUser.permission < PERMISSION.Admin}
                           onChange={newReviewer => handleChange("reviewer", newReviewer)}
                           value={changeable.reviewer}/>
            </Field>
            {props.item.accepted && <>
              <Field title="Утвержден" containerOptions={containerOptions}>
                <TextInput maxLength={256}
                           disabled={true}
                           value={props.item.acceptor}/>
              </Field>
              <Field title="Дата утверждения" containerOptions={containerOptions}>
                <TextInput maxLength={256}
                           disabled={true}
                           value={apiStringDateToString(props.item.acceptedAt)}/>
              </Field>
            </>}
            <Field title="Дата добавления" containerOptions={containerOptions}>
                <TextInput maxLength={256}
                           disabled={true}
                           value={apiStringDateToString(props.item.addedAt)}/>
            </Field>
            <Field title="Доп. инфо" containerOptions={containerOptions}>
                <TextInput maxLength={256}
                           onChange={newInfo => handleChange("additionalInfo", newInfo)}
                           value={changeable.additionalInfo}/>
            </Field>
            <div className="claimed-item-editor-controls">
                {props.user.apiUser != null && props.user.apiUser.permission >= PERMISSION.GM &&
                  <ActionButton title="Изменить"
                                onClick={() => props.onEdit(changeable)}/>}
                {props.user.apiUser != null && props.user.apiUser.permission >= PERMISSION.Admin &&
                  <ActionButton title="Утвердить"
                                onClick={() => props.onApprove(props.item.id)}/>}
                {props.user.apiUser != null && props.user.apiUser.permission >= PERMISSION.Admin &&
                  <ActionButton title="Удалить"
                                onClick={() => props.onDelete(props.item.id)} />}
            </div>
        </ModalTitle>
    )
};

export default ClaimedItemEditor;