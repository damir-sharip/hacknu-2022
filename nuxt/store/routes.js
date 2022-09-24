export const state = () => ({
    routes: {},
})


export const getters = {
    routes(state) {
        return state.routes;
    },
}

export const mutations = {
    SET_ROUTES(state, routes) {
        state.routes = routes;
    }
}