import {configureStore} from '@reduxjs/toolkit'
import ClaimedItemsReducer from "./model/claimed-items/reducer";
import UserReducer from "./model/user";
import ThemeReducer from "./model/theme";
import CharsheetReviewGeneratorReducer from './model/charsheets/reducer'
import CheckReducer from "./model/economics/checks/reducer";
import ArbiterItemsReducer from "./model/arbiters/items/reducer";
import ArbiterBusinessRewardReducer from './model/arbiters/business/reducer'
import ArbiterEventRewardDistributionReducer from './model/arbiters/event-rewards/reducer'
import GobSearcherReducer from './model/gob-searcher/reducer'
import {TypedUseSelectorHook, useDispatch, useSelector} from "react-redux";
import API from './api';

const ConfiguredStore = (api: API) => configureStore({
    reducer: {
        charsheet: CharsheetReviewGeneratorReducer,
        user: UserReducer,
        theme: ThemeReducer,
        claimedItems: ClaimedItemsReducer,
        checks:  CheckReducer,
        business: ArbiterBusinessRewardReducer,
        eventReward: ArbiterEventRewardDistributionReducer,
        gobSearcher: GobSearcherReducer,
        arbiterItems: ArbiterItemsReducer,
        // comments: commentsReducer,
        // users: usersReducer,
    },
    middleware: getDefaultMiddleware => getDefaultMiddleware({
        thunk: {
            extraArgument: api
        }
    })
})
type StoreType = ReturnType<typeof ConfiguredStore>
// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<StoreType["getState"]>
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = StoreType["dispatch"]
export const useAppDispatch: () => AppDispatch = useDispatch
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector
const store = ConfiguredStore(new API())

export default store
