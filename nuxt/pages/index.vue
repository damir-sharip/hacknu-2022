<template>
    <div class="map__container d-flex w-100">
        <div class="navbar">
            <h2 class="navbar__title">All Traces</h2>
            <Multiselect
                v-model="trace"
                :options="traces"
                :allow-empty="false"
                :taggable="false"
                class="mb-4"
            ></Multiselect>

            <h2 class="navbar__title">History</h2>
            <vue-slider
                v-if="show"
                v-model="tracePoint"
                :max="duration"
                class="custom-vue-slider w-100"
            />

            <b-spinner
                v-else
                variant="light"
                class="my-4"
                label="Loading..."
            ></b-spinner>
            <b-button class="custom-button" variant="light">Start</b-button>

            <Xlsx class="mt-2" />
        </div>
        <div class="map">
            <Map
                v-if="showMap"
                :trace="trace"
                :duration="duration"
                :tracePoint="tracePoint"
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
            trace: "dev11",
            tracePoint: 0,
            traceLabes: ["A", "B", "C", "D"],
            show: false,
            showMap: true,
            duration: 10000000000000000,
        };
    },
    computed: {
        ...mapGetters({
            routes: "routes/routes",
        }),
        traces() {
            return Object.keys(this.routes || {});
        },
    },
    watch: {
        trace() {
            this.showMap = false;
            // console.log(this.routes[this.trace], "trace")
            setTimeout(() => (this.showMap = true), 200);
        },
    },
    mounted() {
        setTimeout(() => {
            this.show = true;
        }, 400);
    },
};
</script>
