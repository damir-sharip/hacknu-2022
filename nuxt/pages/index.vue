<template>
    <div class="map__container d-flex w-100">
        <div class="navbar">
            <h2 class="navbar__title mb-4">All Traces</h2>
            <Multiselect
                v-model="trace"
                :options="traces"
                :allow-empty="false"
                :taggable="false"
            ></Multiselect>
            <vue-slider
                v-if="show"
                class="custom-vue-slider w-100"
                :data="traceLabes"
                v-model="tracePoint"
                :marks="true"
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
            <Map v-if="showMap" :trace="trace" />
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
            tracePoint: "B",
            traceLabes: ["A", "B", "C", "D"],
            show: false,
            showMap: true,
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
