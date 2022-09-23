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

import { Loader } from '@googlemaps/js-api-loader';
import * as THREE from 'three';
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js';

const apiOptions = {
  apiKey: 'AIzaSyBKT-QYrE3RY9mY3h_XXMHkfKQBe2jsAWQ',
  version: "beta"
};

const location = {
  latitude: "35.66093428",
  longitude: "139.7290334",
  altitude: "0",
  identifier: "null",
  timestamp: "4875",
  floor_label: "null",
  horizontal_accuracy: "2.314",
  vertical_accuracy: "0.612",
  confidence_in_location_accuracy: "0.6827",
  activity: "UNKNOWN"
}

const mapOptions = {
  "tilt": 67.5,
  "heading": 0,
  "zoom": 18,
  "center": { lat: +location.latitude, lng: +location.longitude },
  "mapId": "f25a14a71e327fab"
}

async function initMap() {
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
    const directionalLight = new THREE.DirectionalLight(0xffffff, 0.25);
    directionalLight.position.set(0.5, -1, 0.5);
    scene.add(directionalLight);

    // load the model    
    loader = new GLTFLoader();
    const source = "pin.gltf";
    loader.load(
      source,
      gltf => {
        gltf.scene.scale.set(15, 15, 15);
        gltf.scene.rotation.x = 180 * Math.PI / 180; // rotations are in radians
        scene.add(gltf.scene);
      }
    );
  }

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
  }

  webGLOverlayView.onDraw = ({ gl, transformer }) => {
    // update camera matrix to ensure the model is georeferenced correctly on the map
    const latLngAltitudeLiteral = {
      lat: mapOptions.center.lat,
      lng: mapOptions.center.lng,
      altitude: +location.altitude
    }

    const matrix = transformer.fromLatLngAltitude({ ...latLngAltitudeLiteral, altitude: latLngAltitudeLiteral.altitude + 20 });
    camera.projectionMatrix = new THREE.Matrix4().fromArray(matrix);

    webGLOverlayView.requestRedraw();
    renderer.render(scene, camera);

    // always reset the GL state
    renderer.resetState();
  }
  webGLOverlayView.setMap(map);
}

(async () => {
  const map = await initMap();
  initWebGLOverlayView(map);
})();