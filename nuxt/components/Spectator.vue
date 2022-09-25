<template>
    <Multiselect
        :value="spectator"
        @input="handleChange"
        :options="spectators"
        :allow-empty="true"
        :taggable="false"
        placeholder="select person"
    ></Multiselect>
</template>

<script>
import { mapGetters } from "vuex";

export default {
    name: "Spectator",
    data() {
        return {
            spectator: null,
        };
    },
    computed: {
        ...mapGetters({
            spectators: "spectator/spectators",
        }),
    },
    methods: {
        handleChange(value) {
            this.spectator = value;
            this.$store.commit("spectator/SET_SPECTATOR", value);
        },
    },
    mounted() {
        this.$root.$on("resetSpectator", () => (this.spectator = null));
    },
};
</script>

<style scoped>
</style>