import React from 'react';
import ContentTitle from "../../common/content-title";
import ActionButton from "../../common/action-button";
import './styles.css'

type WelcomeExistingUserProps = {
    name: string,
    permissionName: string,
    isConnected: boolean
    onLogout: (() => Promise<void>) | (() => void)
    onConnect?: (() => Promise<void>) | (() => void)
}

const WelcomeExistingUser = (props: WelcomeExistingUserProps) => {
    return (
        <ContentTitle className="welcome-message" title="Аккаунт" collapsable={false}>
            <p>Привет, {props.name}</p>
            <p>Уровень доступа: {props.permissionName}</p>
            <span className="welcome-message__controls">
                <ActionButton title="Выйти" onClick={props.onLogout}/>
                {!props.isConnected && <ActionButton title="Привязать DM" onClick={props.onConnect}/>}
            </span>
        </ContentTitle>
    );
};

export default React.memo(WelcomeExistingUser);