export const state = () => ({
    spectator: null,
    spectators: [],
    identifiers: {}
})


export const getters = {
    spectator(state) {
        return state.spectator;
    },
    spectators(state) {
        return state.spectators;
    },
    identifiers(state) {
        return state.identifiers;
    }
}

export const mutations = {
    SET_SPECTATOR(state, spectator) {
        state.spectator = spectator;
    },
    SET_SPECTATORS(state, spectators) {
        state.spectators = spectators;
    },
    SET_IDENTIFIERS(state, identifiers) {
        state.identifiers = identifiers;
    }
}