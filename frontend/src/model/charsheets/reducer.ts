import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {Rate} from "./types";
import DefaultCharsheetReviewState from "./types";
import {APIResponseKnownError, ControllerException} from "../exceptions";
import { createAppAsyncThunk } from "../../thunk";


export const generateCharsheetTemplate = createAppAsyncThunk("charsheet/generate", async (_, thunkAPI) => {
    const token = thunkAPI.getState().user.user.token;
    if(!token){
        throw new Error("Not authorized")
    }
    const info = thunkAPI.getState().charsheet.info
    if(!info.charName || !info.reviewerDiscord || !info.reviewerProfile) return thunkAPI.rejectWithValue(new Error("Поля не заполнены"))
    try {
        function isReview(response: unknown): response is { review: string } {
            return typeof response === "object" && response != null && "review" in response;
        }
        function isError(response: unknown): response is { error: string } {
            return typeof response === "object" && response != null && "error" in response;
        }
        const response = await thunkAPI.extra.createRequest(
            "charReview.generate",
            "",
            JSON.stringify(info),
            token)
        const json: unknown = await response.json()
        if(isReview(json)) return json.review;
        if(isError(json)) throw new ControllerException(json.error)
        return ""
    } catch (e: unknown) {
        if (e instanceof ControllerException || e instanceof APIResponseKnownError) {
            return thunkAPI.rejectWithValue(e)
        }
        throw e
    }
})

const charsheetReviewSLice = createSlice({
    name: 'charsheet',
    initialState: DefaultCharsheetReviewState,
    reducers: {
        updateRate: (state, action: PayloadAction<Rate>) => {
            const newRates = [...state.info.rates]
            const index = newRates.findIndex(rate => rate.rateName === action.payload.rateName)
            if(index === -1) return
            newRates[index] = action.payload;
            state.info.rates = newRates;
            state.info.totalRate = Math.floor(state.info.rates.reduce((acc, curr) => acc + curr.rateValue, 0) / state.info.rates.length)
        },
        setCharName: (state, action: PayloadAction<string>) => {
            state.info.charName = action.payload
        },
        setReviewerProfile: (state, action: PayloadAction<string>) => {
            state.info.reviewerProfile = action.payload
            localStorage.setItem("profileLink", action.payload)
        },
        setReviewerDiscord: (state, action: PayloadAction<string>) => {
            state.info.reviewerDiscord = action.payload
            localStorage.setItem("discordProfile", action.payload)
        }
    },
    extraReducers: builder => {
        builder
            .addCase(generateCharsheetTemplate.fulfilled, (state, action: PayloadAction<string>) => {
                state.result = action.payload
                state.error = ''
            })
            .addCase(generateCharsheetTemplate.rejected, (state, action: PayloadAction<unknown>) => {
                if(action.payload instanceof Error) {
                    state.error = action.payload.message
                    setTimeout(() => state.error = "", 1500)
                }
            })
    }
});

export default charsheetReviewSLice.reducer

export const { updateRate, setCharName, setReviewerProfile, setReviewerDiscord } = charsheetReviewSLice.actions;

