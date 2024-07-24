import React from 'react';
import Protector from "../../protector";
import {PERMISSION} from "../../../model/user";

const ConnectPage = () => {
    return (
        <Protector accessLevel={PERMISSION.Player}>
            <div></div>
        </Protector>
    );
};

export default React.memo(ConnectPage);