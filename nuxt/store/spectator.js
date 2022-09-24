export const state = () => ({
    spectator: null,
    spectators: [],
})


export const getters = {
    spectator(state) {
        return state.spectator;
    },
    spectators(state) {
        return state.spectators;
    },
}

export const mutations = {
    SET_SPECTATOR(state, spectator) {
        state.spectator = spectator;
    },
    SET_SPECTATORS(state, spectators) {
        state.spectators = spectators;
    }
}