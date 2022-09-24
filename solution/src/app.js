// Copyright 2021 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { Loader } from "@googlemaps/js-api-loader";
import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";
import data from "./data.json";

const dev1 = data["dev2"][0];

const apiOptions = {
  apiKey: "AIzaSyBKT-QYrE3RY9mY3h_XXMHkfKQBe2jsAWQ",
  version: "beta",
};

const location = {
  latitude: dev1.latitude,
  longitude: dev1.longitude,
  altitude: dev1.altitude,
  identifier: dev1.identifier,
  timestamp: dev1.timestamp,
  floor_label: dev1.floor_label,
  horizontal_accuracy: dev1.horizontal_accuracy,
  vertical_accuracy: dev1.vertical_accuracy,
  confidence_in_location_accuracy: dev1.confidence_in_location_accuracy,
  activity: dev1.activity,
};

const mapOptions = {
  tilt: 67.5,
  heading: 0,
  zoom: 20,
  center: { lat: +location.latitude, lng: +location.longitude },
  mapId: "f25a14a71e327fab",
};

function minMaxAccuracy(obj, type) {
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
}

function generateObject(object) {
  // const mainTargetGeometry = new THREE.SphereGeometry(6, 6, 48);
  // const mainTargetMaterial = new THREE.MeshBasicMaterial({
  //   color: "#0031ff",
  // });
  // const mainTargetSphere = new THREE.Mesh(
  //   mainTargetGeometry,
  //   mainTargetMaterial
  // );

  const mainTargetGeometry = new THREE.ConeGeometry(5, 20, 32);
  const mainTargetMaterial = new THREE.MeshBasicMaterial({
    color: 0xffff00,
    transparent: true,
    opacity: 0.1,
  });
  const mainTargetCone = new THREE.Mesh(mainTargetGeometry, mainTargetMaterial);

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

  console.log(
    object.horizontal_accuracy,
    minMaxAccuracy(dev1, "min"),
    minMaxAccuracy(dev1, "max")
  );

  const horizontalMinAccuracyGeometry = new THREE.CircleGeometry(
    minMaxAccuracy(dev1, "min")
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
    minMaxAccuracy(dev1, "max")
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
}

async function initMap() {
  console.log(location, "dev1");
  const mapDiv = document.getElementById("map");
  const apiLoader = new Loader(apiOptions);
  await apiLoader.load();
  return new google.maps.Map(mapDiv, mapOptions);
}

function initWebGLOverlayView(map) {
  let scene, renderer, camera, loader;
  const webGLOverlayView = new google.maps.WebGLOverlayView();

  webGLOverlayView.onAdd = () => {
    // set up the scene
    scene = new THREE.Scene();
    camera = new THREE.PerspectiveCamera();
    const ambientLight = new THREE.AmbientLight(0xffffff, 0.75); // soft white light
    scene.add(ambientLight);
    const directionalLight = new THREE.DirectionalLight("0xffffff", 0.25);
    directionalLight.position.set(0.5, -1, 0.5);
    scene.add(directionalLight);

    const mainTargetSphere = generateObject(dev1);
    scene.add(mainTargetSphere);
  };

  webGLOverlayView.onContextRestored = ({ gl }) => {
    // create the three.js renderer, using the
    // maps's WebGL rendering context.
    renderer = new THREE.WebGLRenderer({
      canvas: gl.canvas,
      context: gl,
      ...gl.getContextAttributes(),
    });
    renderer.autoClear = false;

    // const left = document.getElementById("left");
    // const right = document.getElementById("right");

    // const leftInterval = setInterval(() => {
    //   mapOptions.heading += 0.2;
    // }, 200)
    // function onLeftClick() {

    //   // while (mapOptions.heading <= 360) {
    //   //   mapOptions.heading += 0.2;
    //   //   console.log("left")
    //   // }
    //   // renderer.setAnimationLoop(null);
    // }
    // function onRightClick() {
    //   while (mapOptions.heading > 0) {
    //     mapOptions.heading -= 0.2;
    //     console.log("right")
    //   }
    //   // renderer.setAnimationLoop(null);
    // }

    // left.addEventListener("mousedown", onLeftClick);
    // left.addEventListener("mouseup", onLeftClick);
    // right.addEventListener("mousedown", onRightClick)

    // wait to move the camera until the 3D model loads
    // loader.manager.onLoad = () => {
    //   renderer.setAnimationLoop(() => {
    //     map.moveCamera({
    //       "tilt": mapOptions.tilt,
    //       "heading": mapOptions.heading,
    //       "zoom": mapOptions.zoom
    //     });

    //     // rotate the map 360 degrees
    //     if (mapOptions.tilt < 67.5) {
    //       mapOptions.tilt += 0.5
    //     } else if (mapOptions.heading <= 360) {
    //       mapOptions.heading += 0.2;
    //     } else {
    //       renderer.setAnimationLoop(null)
    //     }
    //   });
    // }
  };

  webGLOverlayView.onDraw = ({ gl, transformer }) => {
    // update camera matrix to ensure the model is georeferenced correctly on the map
    const latLngAltitudeLiteral = {
      lat: mapOptions.center.lat,
      lng: mapOptions.center.lng,
      altitude: +location.altitude,
    };

    const matrix = transformer.fromLatLngAltitude({
      ...latLngAltitudeLiteral,
      altitude: latLngAltitudeLiteral.altitude + 20,
    });
    camera.projectionMatrix = new THREE.Matrix4().fromArray(matrix);

    webGLOverlayView.requestRedraw();
    renderer.render(scene, camera);

    // always reset the GL state
    renderer.resetState();
  };
  webGLOverlayView.setMap(map);
}

(async () => {
  const map = await initMap();
  initWebGLOverlayView(map);
})();
