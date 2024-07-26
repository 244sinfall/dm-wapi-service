import React, {useMemo} from 'react';
import {PermissionValue} from "../../model/user/types";
import AccountManager from "../auth";
import ProtectorFrame from "../../components/protector/frame";
import ProtectorNoAccess from "../../components/protector/no-access";
import {useAppDispatch, useAppSelector} from "../../store";
import ProtectorNotConnected from '../../components/protector/not-connected';
import { useNavigate } from 'react-router-dom';
import { connectToDarkmoon } from '../../model/user/';

const Protector = (props: {children: React.ReactNode[] | React.ReactNode, accessLevel: PermissionValue}) => {
    const currentUser = useAppSelector(state => state.user.user)
    const nav = useNavigate()
    const dispatch = useAppDispatch()
    const protector = useMemo(() => {
        if(!currentUser.token) {
            return <AccountManager />
        }
        if(currentUser.apiUser == null) {
            return <ProtectorNotConnected onSubmit={async (key) => {
                await dispatch(connectToDarkmoon(key))
                nav('/')
            }} onCancel={() => nav('/')}/>
        }
        if(currentUser.apiUser.permission < props.accessLevel) {
            return <ProtectorNoAccess/>
        }
        
    }, [currentUser, props.accessLevel, dispatch, nav])
    
    return (
        <>
            {protector && 
            <ProtectorFrame>
                {protector}
            </ProtectorFrame>}
            {!protector && props.children}
        </>
    );
};

export default Protector;