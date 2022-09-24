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
      CAR_FRONT: new Vector3(0, 1, 0),

      API_OPTIONS: {
        apiKey: "AIzaSyBKT-QYrE3RY9mY3h_XXMHkfKQBe2jsAWQ",
        version: "beta",
      },
      VIEW_PARAMS: {
        center: { lat: 53.554486, lng: 10.007479 },
        zoom: 18,
        heading: 40,
        tilt: 65,
        mapId: "f25a14a71e327fab",
      },

      ANIMATION_DURATION: 12000,
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
      dev1: [
        {
          latitude: 53.554486,
          longitude: 10.007479,
          altitude: 0,
          identifier: null,
          timestamp: 4875,
          floor_label: null,
          horizontal_accuracy: 2.314,
          vertical_accuracy: 0.612,
          confidence_in_location_accuracy: 0.6827,
          activity: "UNKNOWN",
        },
      ],
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

    //   const mainTargetGeometry = new THREE.ConeGeometry(5, 20, 32);
    //   const mainTargetMaterial = new THREE.MeshBasicMaterial({
    //     color: 0xffff00,
    //     transparent: true,
    //     opacity: 0.1,
    //   });
    //   const mainTargetCone = new THREE.Mesh(
    //     mainTargetGeometry,
    //     mainTargetMaterial
    //   );

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
        this.minMaxAccuracy(this.dev1, "min")
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
        this.minMaxAccuracy(this.dev1, "max")
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
      quaternion.setFromAxisAngle(new THREE.Vector3(0, 1, 0), Math.PI / 2);
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
          (obj.horizontal_accuracy * obj.confidence_in_location_accuracy) / 100
        );
      } else {
        return (
          obj.horizontal_accuracy +
          (obj.horizontal_accuracy * obj.confidence_in_location_accuracy) / 100
        );
      }
    },
  },
  async mounted() {
    const map = await this.initMap();

    const overlay = new ThreeJSOverlayView({ lat: 53.554486, lng: 10.007479 });
    const scene = overlay.getScene();

    overlay.setMap(map);

    // create a Catmull-Rom spline from the points to smooth out the corners
    // for the animation
    const points = this.ANIMATION_POINTS.map((p) =>
      overlay.latLngAltToVector3(p)
    );
    const curve = new CatmullRomCurve3(points, true, "catmullrom", 0.2);
    curve.updateArcLengths();

    const trackLine = this.createTrackLine(curve);
    scene.add(trackLine);

    const obj = this.generateObject(this.dev1);
    scene.add(obj);

    overlay.requestRedraw();

    // the update-function will animate the car along the spline
    overlay.update = () => {
      trackLine.material.resolution.copy(overlay.getViewportSize());

      if (!obj) return;

      const animationProgress =
        (performance.now() % this.ANIMATION_DURATION) / this.ANIMATION_DURATION;

      curve.getPointAt(animationProgress, obj.position);
      curve.getTangentAt(animationProgress, this.tmpVec3);
      obj.quaternion.setFromUnitVectors(this.CAR_FRONT, this.tmpVec3);

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
