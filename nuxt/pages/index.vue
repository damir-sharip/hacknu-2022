<template>
    <div class="map__container d-flex w-100">
        <div class="navbar">
            <h2 class="navbar__title mb-4">All Traces</h2>
            <Multiselect
                v-model="trace"
                :options="options.traces"
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
            <div class="d-flex mt-2">
                <h3>Info about current trace</h3>
            </div>
            <div class="d-flex mt-2">
                <Xlsx />
            </div>

        </div>
        <div class="map">
            <Map :trace="trace" />
        </div>
    </div>
</template>

<script>
import options from "~/helpers/options";

export default {
    name: "IndexPage",
    data() {
        return {
            options,
            trace: "dev11",
            tracePoint: "B",
            traceLabes: ["A", "B", "C", "D"],
            show: false,
        };
    },
    mounted() {
        console.log(options, "options");
        setTimeout(() => {
            this.show = true;
        }, 400);
    },
};
</script>
