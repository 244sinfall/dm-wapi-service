import {
    createSlice,
    PayloadAction
} from "@reduxjs/toolkit";
import {PERMISSION, PermissionValue} from "./index";
import {signOut, User} from "firebase/auth";
import {auth, db} from "../../services/services/authorizer/firebase";
import {doc, getDoc} from "firebase/firestore";
import { createAppAsyncThunk } from "../reduxTypes";

export type UserStateUserInfo = {
    name: string,
    email: string,
    permission: PermissionValue,
    apiUser: ApiAuthUser | null 
}

export const DefaultUserState = {
    name: "Гость",
    email: "",
    permission: PERMISSION.Player,
    apiUser: null
}

interface UserState {
    user: UserStateUserInfo
    isLoading: boolean
}

const userInitialState: UserState = {
    user: DefaultUserState,
    isLoading: !!localStorage.getItem("hasFirebaseSession")
}

type ApiUser = {
    id: number
}

type ApiScope = {
    securityLevel: number,
    rbac: number[],
    root: boolean
}

type ApiAuthUser = {
    user: ApiUser,
    scope: ApiScope,
    integrationUserId: string,
    permission: number
}

export const connectToDarkmoon = createAppAsyncThunk("user/connectToDarkmoon", async(code: string, thunkAPI) => {
    const response = await thunkAPI.extra.get("API").createRequest("users.connect", "", JSON.stringify({code: code}));
    const responseJson = await response.json();
    if(!response.ok) {
        return null;
    }
    return responseJson as ApiAuthUser
})

export const restoreSession = createAppAsyncThunk("user/restoreSession", async(user: User, thunkAPI) => {
    let apiAuthUser: ApiAuthUser | null = null
    let permission: PermissionValue = 0
    const response = await thunkAPI.extra.get("API").createRequest("users.me")
    if(response.ok){
        const json = await response.json();
        const userData = json as ApiAuthUser
        apiAuthUser = userData
        permission = userData.permission as PermissionValue
    }
    localStorage.setItem("hasFirebaseSession", "true")
    return {name: user.displayName ?? "Пользователь", email: user.email ?? "",
        permission, apiUser: apiAuthUser}
})

export const destroySession = createAppAsyncThunk("user/destroySession", async() => {
    await signOut(auth)
    localStorage.removeItem("hasFirebaseSession")
})

export const userSlice = createSlice({name: "user", initialState: userInitialState,
    reducers: {},
    extraReducers: builder => builder
        .addCase(destroySession.fulfilled, state => {
            state.user = DefaultUserState
        })
        .addCase(restoreSession.pending, state => {state.isLoading = true})
        .addCase(restoreSession.fulfilled, (state, action: PayloadAction<UserStateUserInfo>) => {
            state.user = action.payload
            state.isLoading = false
        })
        .addCase(restoreSession.rejected, (state) => {
            state.isLoading = false
        })
        .addCase(connectToDarkmoon.pending, state => {state.isLoading = true})
        .addCase(connectToDarkmoon.fulfilled, (state, action: PayloadAction<ApiAuthUser | null>) => {
            if(action.payload)
            {
                state.user.apiUser = action.payload
                state.user.permission = action.payload.permission as PermissionValue
            }
            state.isLoading = false
        })
})