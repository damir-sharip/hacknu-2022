<template>
  <div class="map__container d-flex w-100">
    <div class="navbar align-items-start">
      <div class="custom-logo">
        <img src="~/assets/images/logo.svg" />
      </div>
      <div class="custom-form w-100">
        <h3 class="navbar__title">Upload Excel Document</h3>
        <Xlsx class="mt-2 mb-4 w-100" />
        <h3 class="navbar__title">Path</h3>
        <Multiselect
          v-model="trace"
          :options="traces"
          :allow-empty="false"
          :taggable="false"
          class="mb-4"
        ></Multiselect>

        <h3 class="navbar__title mb-2">Person</h3>
        <Spectator class="mb-4" />
        <h3 class="navbar__title">Timeline</h3>
        <div v-for="spectator in spectators" :key="spectator" class="w-100">
          <div
            v-if="
              show &&
              sliderData[spectator].length >= 2 &&
              sliderData[spectator][0] !== sliderData[spectator][1]
            "
          >
            <div class="timeline__name">{{ spectator }}</div>
            <vue-slider
              :value="tracePoints[spectator]"
              @change="handleChange($event, spectator)"
              :data="sliderData[spectator].map((_, index) => index)"
              :marks="true"
              class="custom-vue-slider w-100"
            />
          </div>
        </div>
        <b-button class="custom-button" @click="togglePlay">{{
          isPlay ? "Stop" : "Play"
        }}</b-button>
      </div>
    </div>
    <div class="map">
      <Map
        v-if="showMap"
        :trace="trace"
        :duration="duration"
        :traceDurations="traceDurations"
        :isPlay="isPlay"
      />
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  name: "IndexPage",
  data() {
    return {
      trace: "example",
      show: false,
      showMap: true,
      duration: 18000,
      isPlay: false,
      tracePoints: {},
      sliderData: {},
      tracePoint: 0,
      traceDurations: {},
    };
  },
  computed: {
    ...mapGetters({
      routes: "routes/routes",
      spectators: "spectator/spectators",
      identifiers: "spectator/identifiers",
    }),
    traces() {
      return Object.keys(this.routes || {});
    },
  },
  watch: {
    trace() {
      this.$store.commit("spectator/SET_SPECTATOR", null);
      this.$root.$emit("resetSpectator");
      this.showMap = false;
      setTimeout(() => (this.showMap = true), 200);
    },
    spectators() {
      this.tracePoints = {};
      this.traceDurations = {};
      this.spectators.forEach((spectator) => {
        this.tracePoints[spectator] = null;
        this.traceDurations[spectator] = null;
      });
    },
    identifiers() {
      this.sliderData = {};
      for (const key in this.identifiers) {
        let timestamps = this.identifiers[key].map((point) => point?.timestamp);
        this.sliderData[key] = timestamps.filter((item) => !!item);
      }
    },
  },
  methods: {
    togglePlay() {
      this.isPlay = !this.isPlay;
    },
    handleChange(value, spectator) {
      this.tracePoints[spectator] = value;
      const max = this.sliderData[spectator].length;
      this.traceDurations[spectator] = (value * this.duration) / max;
    },
  },
  mounted() {
    setTimeout(() => (this.show = true), 100);
  },
};
</script>

<style>
.vue-slider-mark-label {
  display: none !important;
}
</style>
