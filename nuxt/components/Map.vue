<template>
    <div id="map"></div>
</template>

<script>
import { Loader } from "@googlemaps/js-api-loader";
import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

export default {
    data() {
        return {
            map: null,
            scene: null,
            renderer: null,
            camera: null,
            loader: null,
            apiOptions: {
                apiKey: "AIzaSyBKT-QYrE3RY9mY3h_XXMHkfKQBe2jsAWQ",
                version: "beta",
            },
            location: {
                latitude: "35.66093428",
                longitude: "139.7290334",
                altitude: "0",
                identifier: "null",
                timestamp: "4875",
                floor_label: "null",
                horizontal_accuracy: "2.314",
                vertical_accuracy: "0.612",
                confidence_in_location_accuracy: "0.6827",
                activity: "UNKNOWN",
            },
            mapOptions: {
                tilt: 67.5,
                heading: 0,
                zoom: 18,
                center: {
                    lat: 35.66093428,
                    lng: 139.7290334,
                },
                mapId: "f25a14a71e327fab",
            },
        };
    },
    methods: {
        async initMap() {
            const mapDiv = document.getElementById("map");
            const apiLoader = new Loader(this.apiOptions);
            await apiLoader.load();
            this.map = new google.maps.Map(mapDiv, this.mapOptions);
        },
        initWebGLOverlayView() {
            const webGLOverlayView = new google.maps.WebGLOverlayView();

            webGLOverlayView.onAdd = () => {
                // set up the scene
                this.scene = new THREE.Scene();
                this.camera = new THREE.PerspectiveCamera();
                const ambientLight = new THREE.AmbientLight(0xffffff, 0.75); // soft white light
                this.scene.add(ambientLight);
                const directionalLight = new THREE.DirectionalLight(
                    0xffffff,
                    0.25
                );
                directionalLight.position.set(0.5, -1, 0.5);
                this.scene.add(directionalLight);

                // load the model
                this.loader = new GLTFLoader();
                const source = "../assets/pin.gltf";
                this.loader.load(source, (gltf) => {
                    gltf.scene.scale.set(25, 25, 25);
                    gltf.scene.rotation.x = (180 * Math.PI) / 180; // rotations are in radians
                    this.scene.add(gltf.scene);
                });
            };
            webGLOverlayView.onContextRestored = ({ gl }) => {
                // create the three.js renderer, using the
                // maps's WebGL rendering context.
                this.renderer = new THREE.WebGLRenderer({
                    canvas: gl.canvas,
                    context: gl,
                    ...gl.getContextAttributes(),
                });
                this.renderer.autoClear = false;
            };

            webGLOverlayView.onDraw = ({ gl, transformer }) => {
                // update camera matrix to ensure the model is georeferenced correctly on the map
                const latLngAltitudeLiteral = {
                    lat: this.mapOptions.center.lat,
                    lng: this.mapOptions.center.lng,
                    altitude: +this.location.altitude,
                };

                const matrix = transformer.fromLatLngAltitude({
                    ...latLngAltitudeLiteral,
                    altitude: latLngAltitudeLiteral.altitude + 20,
                });
                this.camera.projectionMatrix = new THREE.Matrix4().fromArray(
                    matrix
                );

                webGLOverlayView.requestRedraw();
                this.renderer.render(this.scene, this.camera);

                // always reset the GL state
                this.renderer.resetState();
            };
            webGLOverlayView.setMap(this.map);
        },
    },
    async mounted() {
        await this.initMap();
        this.initWebGLOverlayView();
    },
};
</script>

<style>
#map {
    height: 800px;
    margin: 0;
    padding: 0;
}
</style>
