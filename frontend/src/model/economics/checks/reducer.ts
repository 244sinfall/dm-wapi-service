import {createSlice, PayloadAction} from "@reduxjs/toolkit";

import {
    CheckResponse,
    ChecksDefaultState,
    CheckStatusValue,
    CheckTableParams,
    CheckTableParamsCompanion,
    ICheck
} from "./types";
import {APIResponseKnownError} from "../../exceptions";
import { createAppAsyncThunk } from "../../../thunk";


export const fetchChecks = createAppAsyncThunk("checks/fetch", async(_, thunkAPI) => {
    let delimiter = "&"
    let query = "?"
    const state = thunkAPI.getState().checks.params
    for(let prop in state) {
        if(!CheckTableParamsCompanion.is(prop)) continue
        if(state[prop]) {
            if(prop === "status" && state[prop] === "Все") continue;
            if(prop === "category" && state[prop] === "Все получатели") continue;
            if(query.length !== 1) query += delimiter
            query += `${prop}=${String(prop === "status" ? CheckStatusValue[state[prop] as keyof typeof CheckStatusValue] : state[prop])}`
        }
    }
    function isApiResponseIsData(data: unknown): data is CheckResponse {
        return typeof data === "object" && data != null && "checks" in data
    }
    const response = await thunkAPI.extra.createRequest("checks.get",
        query)
    const data = await response.json()
    if(!isApiResponseIsData(data))
        throw new APIResponseKnownError(response)
    data.types = ["Все получатели", ...data.types.filter(t => t && t !== "-")]
    data.updatedAt = new Date(data.updatedAt).toLocaleString("ru")
    return data
})

export const checksSlice = createSlice({
    name: 'checks',
    initialState: ChecksDefaultState,
    reducers: {
        setSelectedCheck: (state, action: PayloadAction<ICheck>) => {
            state.selectedCheck = action.payload
        },
        removeSelectedCheck: (state) => {
            state.selectedCheck = null
        },
        setParams: (state, action: PayloadAction<CheckTableParams>) => {
            state.params = action.payload
        }
    },
    extraReducers: (builder) => {
        builder
            .addCase(fetchChecks.pending, (state) => {
                state.isLoading = true
            })
            .addCase(fetchChecks.fulfilled, (state, action: PayloadAction<CheckResponse>) => {
                state.result = action.payload
                state.isLoading = false
                state.error = ""
            })
            .addCase(fetchChecks.rejected, (state, action: PayloadAction<unknown>)=> {
                if(action.payload instanceof Error) {
                    state.error = action.payload.message
                }
                state.isLoading = false
            })
    }
});

export const { setSelectedCheck, removeSelectedCheck, setParams } = checksSlice.actions;
export default checksSlice.reducer