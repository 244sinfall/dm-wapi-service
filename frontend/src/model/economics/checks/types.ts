export type CheckStatus = "Ожидает" | "Закрыт" | "Отказан"

export const CheckStatusCompanion = {
    list: (): CheckStatus[] => ["Ожидает", "Закрыт", "Отказан"],
}
export const CheckStatusValue: Record<CheckStatus | "Все", string> = {
    "Все": "", "Закрыт": "Закрыт", "Ожидает": "Ожидает", "Отказан": "Отказан"

}

type CheckUser = {
    id: number,
    nickname: string,
    gameId: number,
}

export type ICheck = {
    id: number,
    date: string,
    senderUser: CheckUser,
    receiver: string,
    subject: string,
    body: string,
    money: number,
    gmUser: CheckUser,
    status: CheckStatus,
    items: string
}



export type CheckResponse = {
    checks: ICheck[],
    count: number,
    filteredCount: number,
    types: string[],
    updatedAt: string
}

export const CheckTableParamsCompanion = {
    default: (): CheckTableParams => ({
        limit: 50,
        skip: 0,
        category: "",
        search: "",
        status: "Все"
    }),
    is: (data: string): data is keyof CheckTableParams => {
       return data === "limit" || data === "skip" || data === "category" || data === "search" || data === "status"
    }
}

export type CheckTableParams = {
    limit: 20 | 50 | 100,//Количество чеков на одной странице
    skip: number,// Количество чеков, которые пропускаются (количество чеков на странице * номер страницы)
    category: string,
    search: string,
    status: CheckStatus | "Все"
}

export type ChecksState = {
    params: CheckTableParams,
    isLoading: boolean,
    result: CheckResponse | null
    selectedCheck: ICheck | null,
    error: string,
}

export const ChecksDefaultState: ChecksState = {
    params: CheckTableParamsCompanion.default(),
    isLoading: false,
    result: null,
    selectedCheck: null,
    error: ""
}