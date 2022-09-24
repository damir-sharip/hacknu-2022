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
            ANIMATION_DURATION: 200000,
            ROUTES: [],
            ANIMATION_POINTS: [
                { lat: 53.554473, lng: 10.008226 },
                { lat: 53.554913, lng: 10.008124 },
                { lat: 53.554986, lng: 10.007928 },
                { lat: 53.554775, lng: 10.006363 },
                { lat: 53.554674, lng: 10.006383 },
                { lat: 53.554473, lng: 10.006681 },
                { lat: 53.554363, lng: 10.006971 },
                { lat: 53.554453, lng: 10.008091 },
                { lat: 53.554424, lng: 10.008201 },
                { lat: 53.554473, lng: 10.008226 },
            ],
            tmpVec3: new Vector3(),
        };
    },
    methods: {
        async initMap() {
            const mapDiv = document.getElementById("map");
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
            const mainTargetGeometry = new THREE.SphereGeometry(6, 6, 48);
            const mainTargetMaterial = new THREE.MeshBasicMaterial({
                color: "#0031ff",
            });
            const mainTargetCone = new THREE.Mesh(
                mainTargetGeometry,
                mainTargetMaterial
            );

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
                this.minMaxAccuracy(this.ANIMATION_POINTS[0], "min")
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
                this.minMaxAccuracy(this.ANIMATION_POINTS[0], "max")
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
    },
    async mounted() {
        routes.dev11.map((route) => {
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

        const map = await this.initMap();

        const overlay = new ThreeJSOverlayView({
            lat: 51.46988,
            lng: -0.45197,
        });
        const scene = overlay.getScene();

        overlay.setMap(map);

        this.ANIMATION_POINTS = routes.dev11.map((route) => {
            return {
                lat: route.latitude,
                lng: route.longitude,
                altitude: route.altitude,
            };
        });

        // create a Catmull-Rom spline from the points to smooth out the corners
        // for the animation

        const identifiers = { Bob: {}, Jane: {} };

        const points = this.IDENTIFIERS.Bob.map((p) =>
            overlay.latLngAltToVector3(p)
        );
        const curve = new CatmullRomCurve3(points, false, "catmullrom", 0.2);
        curve.updateArcLengths();
        identifiers.Bob.curve = curve; // 

        const points2 = this.IDENTIFIERS.Jane.map((p) =>
            overlay.latLngAltToVector3(p)
        );
        const curve2 = new CatmullRomCurve3(points2, false, "catmullrom", 0.2);
        curve2.updateArcLengths();

        const trackLine = this.createTrackLine(curve);
        identifiers.Bob.trackLine = trackLine; //
        scene.add(trackLine);

        const obj = this.generateObject(this.IDENTIFIERS.Bob[0]);
        identifiers.Bob.obj = obj; // 
        scene.add(obj);

        const trackLine2 = this.createTrackLine(curve2);
        scene.add(trackLine2);

        const obj2 = this.generateObject(this.IDENTIFIERS.Jane[0]);
        scene.add(obj2);

        overlay.requestRedraw();

        // the update-function will animate the object along the spline
        overlay.update = () => {
            identifiers.Bob.trackLine.material.resolution.copy(overlay.getViewportSize());
            trackLine2.material.resolution.copy(overlay.getViewportSize());

            if (!identifiers.Bob.obj) return;

            const animationProgress =
                (performance.now() % this.ANIMATION_DURATION) /
                this.ANIMATION_DURATION;

            identifiers.Bob.curve.getPointAt(animationProgress, identifiers.Bob.obj.position);
            identifiers.Bob.curve.getTangentAt(animationProgress, this.tmpVec3);
            identifiers.Bob.obj.quaternion.setFromUnitVectors(this.CAR_FRONT, this.tmpVec3);

            curve2.getPointAt(animationProgress, obj2.position);
            curve2.getTangentAt(animationProgress, this.tmpVec3);
            obj2.quaternion.setFromUnitVectors(this.CAR_FRONT, this.tmpVec3);

            overlay.requestRedraw();
        };
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
