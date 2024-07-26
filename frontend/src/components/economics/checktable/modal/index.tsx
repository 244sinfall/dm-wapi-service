import React from "react";
import ModalTitle from "../../../../components/common/modal-title";
import CheckModalCopyOption from "./modal-field";
import './styles.css'
import {ICheck} from "../../../../model/economics/checks/types";


function CheckTableModal(props: {check: ICheck, closeHandler: () => void}) {
    const rejectCommand = `.check return ${props.check.id}`
    const openCommand = `.check open ${props.check.id}`
    const closeCommand = `.check close ${props.check.id}`
    return (
        <ModalTitle className="check-info-modal" title="Макросы для чека" closeCallback={props.closeHandler}>
            {props.check.status === "Отказан" && <p>Этот чек отказан. Его невозможно изменить. Игроку необходимо отправить новый чек.</p>}
            {props.check.status === "Закрыт" && <CheckModalCopyOption title="Переоткрыть чек" command={openCommand}/>}
            {props.check.status === "Ожидает" &&
              <>
                <CheckModalCopyOption title={"Закрыть чек"} command={closeCommand}/>
                <CheckModalCopyOption title={"Отказать чек"} command={rejectCommand}/>
              </>
            }
        </ModalTitle>
    )
}

export default React.memo(CheckTableModal)