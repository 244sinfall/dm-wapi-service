import Config from './config'

export default class API {
    protected config = Config

    /**
     *
     * @param endpoint Эндпоинт из конфигурации
     * @param params Query параметры
     * @param payload Тело запроса
     * @throws APIResponseKnownError
     */
    async createRequest<PayloadType extends BodyInit | void = void>(
                            endpoint: keyof typeof Config["endpoints"],
                            params: string = "",
                            payload?: PayloadType,
                            token?: string): Promise<Response> {

        const init: RequestInit = {}
        init.headers = []
        init.method = Config.endpoints[endpoint].method
        if(token) {
            init.headers.push(["Authorization", token])
        }
        if(typeof payload === "object") {
            init.headers.push(["Content-Type", "application/json"])
        }
        init.headers.push(["Accept", Config.endpoints[endpoint].accept])
        if(payload) init.body = payload
        return await fetch(`${Config.address}${Config.endpoints[endpoint].url}${params}`, init)
        // if(!response.ok) throw new APIResponseKnownError(response)
        // return response
    }
}