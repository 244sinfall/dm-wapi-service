import { useState } from 'react';
import ContentTitle from "../../common/content-title";
import ActionButton from "../../common/action-button";
import './styles.css'
import TextInput from '../../common/text-input';

type ProtectorNotConnectedProps = {
    onSubmit: ((rawText: string) => Promise<void>) | ((rawText: string) => void),
    onCancel: (() => Promise<void>) | (() => void)
}
const ProtectorNotConnected = (props: ProtectorNotConnectedProps) => {
    const [key, setKey] = useState("")
    return (
        <ContentTitle className="protector-no-access" title={"Вы не подключены"} collapsable={false}>
            <p>Ваш аккаунт не привязан к Darkmoon. <br>
            </br>Чтобы получить расширенный доступ, вы должны его привязать.<br>
            </br>Зайдите в Профиль, в верхнем меню выберите Интеграции<br>
            </br>Создайте ключ и вставьте его в поле ниже.</p>
            <span className="protector-no-access-button">
                <TextInput onChange={setKey} />
            </span>
            <span>
                <ActionButton title={"Отправить"} onClick={() => props.onSubmit(key)}/>
                <ActionButton title={"Отмена"} onClick={props.onCancel}/>
            </span>
        </ContentTitle>
    );
};

export default ProtectorNotConnected;