import React, {useCallback, useState} from 'react';
import {PermissionNameByValue} from "../../model/user/types";
import WelcomeExistingUser from "../../components/auth/welcome-existing";
import {useNavigate} from "react-router-dom";
import WelcomeNewUser from "../../components/auth/welcome-new";
import {UserLoginCredentials, UserRegisterCredentials} from "../../model/auth/types";
import { useAppSelector} from "../../store";
import { createUserWithEmailAndPassword, signInWithEmailAndPassword, signOut, updateProfile } from 'firebase/auth';
import { auth } from '../../auth';


const AccountManager = () => {
    const state = useAppSelector(state => ({
        user: state.user.user,
        isLoading: state.user.isLoading
    }))
    const [errMsg, setErrMsg] = useState("")
    const [isCaptchaDone, setIsCaptchaDone] = useState(process.env.NODE_ENV === "development")
    const nav = useNavigate();
    const validateCaptcha = useCallback(() => {
        if(isCaptchaDone) return true;
        setErrMsg("Вы не прошли проверку!")
        return false;
    }, [isCaptchaDone])
    const callbacks = {
        onRegister: useCallback(async(credentials: UserRegisterCredentials) => {
            if(!validateCaptcha()) return
            if (credentials.password !== credentials.passwordCheck) throw Error("Пароли не совпадают")
            if (!credentials.password || !credentials.email || !credentials.login) throw Error("Не все поля заполнены")
            const {email, password} = credentials
            const result = await createUserWithEmailAndPassword(auth, email, password)
            await updateProfile(result.user, {displayName: credentials.login})
        }, [validateCaptcha]),
        onLogin: useCallback(async(credentials: UserLoginCredentials) => {
            if(!validateCaptcha()) return
            if(!credentials.password || !credentials.email || !credentials.email.includes("@")) throw Error("Не все поля заполнены")
            await signInWithEmailAndPassword(auth, credentials.email, credentials.password)
        }, [validateCaptcha]),
        // onReset: useCallback(async(credentials: UserLoginCredentials) => {
        //     if(!validateCaptcha()) return
        //     return authorizer.reset(credentials)
        //         .catch((e: unknown) => {
        //             if(e instanceof Error)
        //                 setErrMsg(e.message)
        //         })
        // }, [authorizer, validateCaptcha]),
        onCaptcha: useCallback((success: boolean) => {
            setIsCaptchaDone(process.env.NODE_ENV === "development" ? true : success);
        }, [])
    }
    return (
        <>
            {state.user.email ?
                <WelcomeExistingUser name={state.user.name || "Гость"}
                                     isConnected={state.user.apiUser != null}
                                     onConnect={() => nav('/connect')}
                                     permissionName={state.isLoading ? "Загрузка..." : PermissionNameByValue[(state.user.apiUser?.permission ?? 0) as keyof typeof PermissionNameByValue]}
                                     onLogout={async()=> await signOut(auth)}/>
            :
                <WelcomeNewUser onLogin={callbacks.onLogin}
                                error={errMsg}
                                isLoading={state.isLoading}
                                onRegister={callbacks.onRegister}
                                onCaptcha={callbacks.onCaptcha}
                                onReset={async() => console.log('hello')}/>
            }
        </>
    );
};

export default React.memo(AccountManager);