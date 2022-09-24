<template>
  <Multiselect
     :value="spectator"
     @input="handleChange"
      :options="spectators"
      :allow-empty="true"
      :taggable="false"
  ></Multiselect>
</template>

<script>
import {mapGetters} from "vuex";

export default {
  name: "Spectator",
  data() {
    return {
      spectator: null,
    }
  },
  computed: {
    ...mapGetters({
        spectators: "spectator/spectators",
    }),
  },
  methods: {
    handleChange(value) {
      this.spectator = value
      this.$store.commit("spectator/SET_SPECTATOR", value);
    },
  },
  watch: {
      // spectator() {
      //       this.$store.commit("spectator/setSpectator", this.spectator);
      //     // this.showMap = false;
      //     // setTimeout(() => (this.showMap = true), 200);
      // },
  },
}
</script>

<style scoped>

</style>