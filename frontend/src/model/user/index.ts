import {
    createSlice,
    PayloadAction,

} from "@reduxjs/toolkit";
import {ApiAuthUser, UserStateUserInfo} from "./types";
import { createAppAsyncThunk } from "../../thunk";
import { User } from "firebase/auth";


const DefaultUserState: UserStateUserInfo = {
    email: "",
    token: "",
    name:"Гость",
    apiUser: null,
}

const userInitialState = {
    user: DefaultUserState,
    isLoading: !!localStorage.getItem("hasFirebaseSession")
}

export const connectToDarkmoon = createAppAsyncThunk("user/connectToDarkmoon", async(code: string, thunkAPI) => {
    const response = await thunkAPI.extra.createRequest("users.connect", "", JSON.stringify({code: code}));
    const responseJson = await response.json();
    if(!response.ok) {
        return null;
    }
    thunkAPI.dispatch(userSlice.actions.setApiUser(responseJson as ApiAuthUser))
})

export const restoreSession = createAppAsyncThunk("user/restoreSession", async(user: User | null, thunkAPI) => {
    if(!user){
        thunkAPI.dispatch(userSlice.actions.resetUser())
        return
    }
    const token = await user.getIdToken()
    let apiAuthUser: ApiAuthUser | null = null
    const response = await thunkAPI.extra.createRequest("users.me", undefined, undefined, token)
    if(response.ok){
        const json = await response.json();
        const userData = json as ApiAuthUser
        apiAuthUser = userData
    }
    thunkAPI.dispatch(userSlice.actions.setUser({ email: user.email ?? "unknown", token, name: user.displayName ?? "Пользователь", apiUser: apiAuthUser }))
})

const userSlice = createSlice({
    name: "user", 
    initialState: userInitialState,
    reducers: {
        setUser(state, action: PayloadAction<UserStateUserInfo>) {
            state.user = action.payload
            state.isLoading = false
        },
        resetUser(state) {
            state.user = DefaultUserState
        },
        setApiUser(state, action: PayloadAction<ApiAuthUser>) {
            state.user.apiUser = action.payload
            state.isLoading = false
        },
    }
})

export default userSlice.reducer