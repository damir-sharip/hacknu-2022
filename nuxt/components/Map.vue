<template>
    <div id="map"></div>
</template>

<script>
import * as THREE from "three";
import { CatmullRomCurve3, Vector3 } from "three";
import { Loader } from "@googlemaps/js-api-loader";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

import { Line2 } from "three/examples/jsm/lines/Line2.js";
import { LineMaterial } from "three/examples/jsm/lines/LineMaterial.js";
import { LineGeometry } from "three/examples/jsm/lines/LineGeometry.js";

import ThreeJSOverlayView from "@ubilabs/threejs-overlay-view";
import routes from "~/static/routes";
import { mapGetters } from "vuex";

export default {
    props: ["trace"],
    data() {
        return {
            map: null,
            scene: null,
            renderer: null,
            loader: null,
            apiOptions: {
                apiKey: "AIzaSyBKT-QYrE3RY9mY3h_XXMHkfKQBe2jsAWQ",
                version: "beta",
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
            CAR_FRONT: new Vector3(0, 1, 0),
            API_OPTIONS: {
                apiKey: "AIzaSyBKT-QYrE3RY9mY3h_XXMHkfKQBe2jsAWQ",
                version: "beta",
            },
            VIEW_PARAMS: {
                center: { lat: 51.46988, lng: -0.45197 },
                zoom: 18,
                heading: 40,
                tilt: 65,
                mapId: "f25a14a71e327fab",
            },
            IDENTIFIERS: {},
            ANIMATION_DURATION: 18000,
            ROUTES: [],
            COLORS: [
                "#fd8ab5",
                "#d30229",
                "#048ce9",
                "#00554c",
                "#25424f",
                "#3e6e83",
                "#ceb9ff",
            ],
            tmpVec3: new Vector3(),
        };
    },
    watch: {
        trace() {
            this.launch();
        },
    },
    computed: {
        ...mapGetters({
            routes: "routes/routes",
            spectator: "spectator/spectator",
        }),
    },
    methods: {
        async initMap() {
            const mapDiv = document.getElementById("map");
            console.log(mapDiv, "mapDiv");
            const apiLoader = new Loader(this.API_OPTIONS);
            await apiLoader.load();
            return new google.maps.Map(mapDiv, {
                ...this.VIEW_PARAMS,
                disableDefaultUI: true,
                backgroundColor: "transparent",
                gestureHandling: "greedy",
            });
        },
        createTrackLine(curve) {
            const numPoints = 10 * curve.points.length;
            const curvePoints = curve.getSpacedPoints(numPoints);
            const positions = new Float32Array(numPoints * 3);

            for (let i = 0; i < numPoints; i++) {
                curvePoints[i].toArray(positions, 3 * i);
            }

            const trackLine = new Line2(
                new LineGeometry(),
                new LineMaterial({
                    color: 0x0f9d58,
                    linewidth: 5,
                })
            );

            trackLine.geometry.setPositions(positions);

            return trackLine;
        },
        generateObject(object) {
            const color =
                this.COLORS[Math.floor(Math.random() * this.COLORS.length)];
            // const mainTargetGeometry = new THREE.SphereGeometry(1, 1, 2);
            // const mainTargetGeometry = new THREE.SphereGeometry(25, 25, 48);

            const mainTargetGeometry = new THREE.CircleGeometry(
                object.horizontal_accuracy || 0,
                32
            );
            const mainTargetMaterial = new THREE.MeshBasicMaterial({
                color: 0xffff00,
            });
            const mainTargetCone = new THREE.Mesh(
                mainTargetGeometry,
                mainTargetMaterial
            );
            // scene.add(circle);

            // const mainTargetMaterial = new THREE.MeshBasicMaterial({ color });
            // const mainTargetCone = new THREE.Mesh(
            //   mainTargetGeometry,
            //   mainTargetMaterial
            // );

            const loader = new GLTFLoader();

            loader.load("./Xbot.glb", (gltf) => {
                let model = gltf.scene;

                gltf.scene.scale.set(25, 25, 25);
                gltf.scene.rotation.x = (180 * Math.PI) / 180;
                gltf.scene.rotation.y = (180 * Math.PI) / 180;
                // gltf.scene.rotation.z = 180 * Math.PI/180;
                const quaternion = new THREE.Quaternion();
                quaternion.setFromAxisAngle(
                    new THREE.Vector3(-1, 0, 0),
                    Math.PI / 2
                );
                gltf.scene.applyQuaternion(quaternion);
                mainTargetCone.add(gltf.scene);
            });

            const horizontalAccuracyGeometry = new THREE.CircleGeometry(
                object.horizontal_accuracy
            );
            const horizontalAccuracyMaterial = new THREE.MeshBasicMaterial({
                color: "#003133",
                transparent: true,
                opacity: 0.5,
            });

            const horizontalAccuracyCircle = new THREE.Mesh(
                horizontalAccuracyGeometry,
                horizontalAccuracyMaterial
            );

            const horizontalMinAccuracyGeometry = new THREE.CircleGeometry(
                this.minMaxAccuracy(object, "min")
            );
            const horizontalMinAccuracyMaterial = new THREE.MeshBasicMaterial({
                color: "#003155",
                transparent: true,
                opacity: 0.5,
            });

            const horizontalMinAccuracyCircle = new THREE.Mesh(
                horizontalMinAccuracyGeometry,
                horizontalMinAccuracyMaterial
            );

            const horizontalMaxAccuracyGeometry = new THREE.CircleGeometry(
                this.minMaxAccuracy(object, "max")
            );
            const horizontalMaxAccuracyMaterial = new THREE.MeshBasicMaterial({
                color: "#003177",
                transparent: true,
                opacity: 0.5,
            });

            const horizontalMaxAccuracyCircle = new THREE.Mesh(
                horizontalMaxAccuracyGeometry,
                horizontalMaxAccuracyMaterial
            );

            const curve = new THREE.EllipseCurve(
                0,
                0, // ax, aY
                5,
                object.vertical_accuracy, // xRadius, yRadius
                0,
                2 * Math.PI, // aStartAngle, aEndAngle
                true, // aClockwise
                0 // aRotation
            );

            const points = curve.getPoints(20);
            const geometry = new THREE.BufferGeometry().setFromPoints(points);

            const material = new THREE.LineBasicMaterial({ color: 0xff0000 });

            // Create the final object to add to the scene
            const ellipse = new THREE.Line(geometry, material);

            const quaternion = new THREE.Quaternion();
            quaternion.setFromAxisAngle(
                new THREE.Vector3(0, 1, 0),
                Math.PI / 2
            );
            ellipse.applyQuaternion(quaternion);

            mainTargetCone.add(ellipse);

            mainTargetCone.add(horizontalAccuracyCircle);
            mainTargetCone.add(horizontalMinAccuracyCircle);
            mainTargetCone.add(horizontalMaxAccuracyCircle);

            return mainTargetCone;
        },
        minMaxAccuracy(obj, type) {
            if (type === "min") {
                return (
                    obj.horizontal_accuracy -
                    (obj.horizontal_accuracy *
                        obj.confidence_in_location_accuracy) /
                        100
                );
            } else {
                return (
                    obj.horizontal_accuracy +
                    (obj.horizontal_accuracy *
                        obj.confidence_in_location_accuracy) /
                        100
                );
            }
        },
        async launch() {
            let items = routes;
            if (Object.keys(this.routes).length) items = this.routes;
            items[this.trace].map((route) => {
                const temp = {
                    lat: route.latitude,
                    lng: route.longitude,
                    altitude: route.altitude,
                };
                if (this.IDENTIFIERS.hasOwnProperty(route.identifier)) {
                    this.IDENTIFIERS[route.identifier].push(temp);
                } else {
                    this.IDENTIFIERS[route.identifier] = [temp];
                }
            });

            for (const key in this.IDENTIFIERS) {
                if (this.IDENTIFIERS[key].length === 1) {
                    this.IDENTIFIERS[key].push(this.IDENTIFIERS[key][0]);
                }
            }
            this.$store.commit("spectator/SET_SPECTATORS", Object.keys(this.IDENTIFIERS));

            const center = Object.values(this.IDENTIFIERS)[0][0];
            this.VIEW_PARAMS.center = {
                lat: center.lat,
                lng: center.lng,
            };

            const map = await this.initMap();

            const overlay = new ThreeJSOverlayView(this.VIEW_PARAMS.center);

            // const overlay = new ThreeJSOverlayView({
            //     lat: 51.46988,
            //     lng: -0.45197,
            // });
            const scene = overlay.getScene();

            overlay.setMap(map);

            // create a Catmull-Rom spline from the points to smooth out the corners
            // for the animation

            const identifiers = {};
            for (const key in this.IDENTIFIERS) {
                identifiers[key] = {};

                const points = this.IDENTIFIERS[key].map((p) =>
                    overlay.latLngAltToVector3(p)
                );
                const curve = new CatmullRomCurve3(
                    points,
                    false,
                    "catmullrom",
                    0.2
                );
                curve.updateArcLengths();
                identifiers[key].curve = curve; //

                const trackLine = this.createTrackLine(curve);
                identifiers[key].trackLine = trackLine; //
                scene.add(trackLine);

                const obj = this.generateObject(this.IDENTIFIERS[key][0]);
                identifiers[key].obj = obj; //
                scene.add(obj);
            }

            overlay.requestRedraw();

            // the update-function will animate the object along the spline
            overlay.update = () => {
                for (const key in this.IDENTIFIERS) {
                    identifiers[key].trackLine.material.resolution.copy(
                        overlay.getViewportSize()
                    );

                    if (!identifiers[key].obj) return;

                    const animationProgress =
                        (performance.now() % this.ANIMATION_DURATION) /
                        this.ANIMATION_DURATION;

                    identifiers[key].curve.getPointAt(
                        animationProgress,
                        identifiers[key].obj.position
                    );
                    identifiers[key].curve.getTangentAt(
                        animationProgress,
                        this.tmpVec3
                    );
                    identifiers[key].obj.quaternion.setFromUnitVectors(
                        this.CAR_FRONT,
                        this.tmpVec3
                    );
                }
                if (this.spectator) {
                    // move map camera to follow the object along the spline
                    // console.log(map)
                    map.setCenter(
                        overlay.vector3ToLatLngAlt(
                            identifiers[this.spectator].obj.position
                        )
                    );
                }
                
                // console.log(overlay/)
                overlay.requestRedraw();
            };
        },
    },
    async mounted() {
        await this.launch();
    },
};
</script>

<style>
#map {
    height: 100%;
    margin: 12px;
    padding: 0;
    border-radius: 20px;
    border: 4px solid white;
}
</style>
