import React, {useCallback, useEffect, useRef} from 'react';
import {PERMISSION} from "../../../model/user/types";
import CheckRow from "../../../components/economics/checktable/table/row";
import {CheckTableParams, ICheck} from "../../../model/economics/checks/types";
import CheckTableWrapper from "../../../components/economics/checktable";
import {useAppDispatch, useAppSelector} from "../../../store";
import {fetchChecks, removeSelectedCheck, setParams, setSelectedCheck} from "../../../model/economics/checks/reducer";

const ChecksTable = () => {
    const state = useAppSelector(state => ({
        checks: state.checks,
        user: state.user.user
    }))
    
    const dispatch = useAppDispatch()
    
    const searchDebounce = useRef<NodeJS.Timeout | null>(null)

    useEffect(() => {
        dispatch(fetchChecks())
    },[state.checks.params, dispatch])
    
    const callbacks = {
        onParamsChange: useCallback(<K extends keyof CheckTableParams,V extends CheckTableParams[K]>(key: K, value: V) => {
            if(key === "search") {
                if(searchDebounce.current) clearTimeout(searchDebounce.current)
                searchDebounce.current = setTimeout(() => {
                    dispatch(setParams({...state.checks.params, [key]: value, skip: 0}))
                }, 1000)
                return
            }
            if(key !== "skip") {
                return dispatch(setParams({...state.checks.params, [key]: value, skip: 0}))
            }
            return dispatch(setParams({...state.checks.params, [key]: value}))
        }, [state.checks.params, dispatch]),
        renderCheck: useCallback((check: ICheck) => {
            return <CheckRow key={check.id} check={check} onClick={() => {
                dispatch(setSelectedCheck(check))
            }} />
        }, [dispatch])
    }

    return (
        <CheckTableWrapper renderCheck={callbacks.renderCheck}
                           isLoading={state.checks.isLoading}
                           onParamsChange={callbacks.onParamsChange}
                           params={state.checks.params}
                           error={state.checks.error}
                           response={state.checks.result}
                           modal={state.checks.selectedCheck && state.user.apiUser && state.user.apiUser.permission >= PERMISSION.Arbiter ? {
                               check: state.checks.selectedCheck,
                               onClose: () => dispatch(removeSelectedCheck())
                           }: undefined}/>
    );
};

export default ChecksTable;