export const state = () => ({
    data: {},
})


export const getters = {
    data() {
        return state.data;
    }
}

export const mutations = {
    setData(state, data) {
        state.data = data;
    }
}